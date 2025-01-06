#!/bin/sh

function exportVariables(){
  export NAME_OF_ORGANIZATION=$NAME_OF_ORGANIZATION
  export DOMAIN_OF_ORGANIZATION=$DOMAIN_OF_ORGANIZATION
  export HOST_COMPUTER_IP_ADDRESS=$HOST_COMPUTER_IP_ADDRESS
  export ORGANIZATION_NAME_LOWERCASE=`echo "$NAME_OF_ORGANIZATION" | tr '[:upper:]' '[:lower:]'`
  
  # Security defaults
  export COUCH_DB_USERNAME=admin
  export COUCH_DB_PASSWORD=adminpw
  export CA_ADMIN_USER=admin
  export CA_ADMIN_PASSWORD=adminpw
  export ORDERER_PASSWORD=adminpw
  export PEER_PASSWORD=peerpw
}

# Define organization-specific variables based on input parameter
function setOrgVariables() {
    case $1 in
        "ISTAD")
            NAME_OF_ORGANIZATION=ISTAD
            HOST_COMPUTER_IP_ADDRESS=172.17.0.1
            DOMAIN_OF_ORGANIZATION=$NAME_OF_ORGANIZATION.idonate.istad.co
            export CA_ADDRESS_PORT=ca.$DOMAIN_OF_ORGANIZATION:7054
            COMPOSE_FILE=docker-compose-istad.yaml
            ;;
        "SampleOrg1")
            NAME_OF_ORGANIZATION=SampleOrg1
            HOST_COMPUTER_IP_ADDRESS=172.17.0.1
            DOMAIN_OF_ORGANIZATION=$NAME_OF_ORGANIZATION.idonate.istad.co
            export CA_ADDRESS_PORT=ca.$DOMAIN_OF_ORGANIZATION:8054
            COMPOSE_FILE=docker-compose-sample1.yaml
            ;;
        "SampleOrg2")
            NAME_OF_ORGANIZATION=SampleOrg2
            HOST_COMPUTER_IP_ADDRESS=172.17.0.1
            DOMAIN_OF_ORGANIZATION=$NAME_OF_ORGANIZATION.idonate.istad.co
            export CA_ADDRESS_PORT=ca.$DOMAIN_OF_ORGANIZATION:9054
            COMPOSE_FILE=docker-compose-sample2.yaml
            ;;
        *)
            echo "Invalid organization name. Use ISTAD, Sample1, or Sample2"
            exit 1
            ;;
    esac
}

# Check if organization name is provided
if [ -z "$1" ]; then
    echo "Please provide organization name (ISTAD, Sample1, or Sample2)"
    exit 1
fi

setOrgVariables $1
exportVariables

# Clean up existing artifacts
./clean-all.sh

# Generate configtx.yaml
sed -e "s/organization_name/$NAME_OF_ORGANIZATION/g" \
    -e "s/organization_domain/$DOMAIN_OF_ORGANIZATION/g" \
    -e "s/ip_address/$HOST_COMPUTER_IP_ADDRESS/g" \
    configtx_template.yaml > configtx.yaml

# Start the CA
docker-compose -p fabric-network -f $COMPOSE_FILE up -d ca
sleep 3

# Generate orderer identities
for ORDERER_NUMBER in 1 2 3; do
    docker exec ca.$DOMAIN_OF_ORGANIZATION /bin/bash -c "cd /etc/hyperledger/artifacts/ && \
    ./orderer-identity.sh $CA_ADDRESS_PORT $DOMAIN_OF_ORGANIZATION $HOST_COMPUTER_IP_ADDRESS \
    $CA_ADMIN_USER $CA_ADMIN_PASSWORD $ORDERER_NUMBER $ORDERER_PASSWORD"
done

# Generate peer identities
for PEER_NUMBER in 1 2; do
    docker exec ca.$DOMAIN_OF_ORGANIZATION /bin/bash -c "cd /etc/hyperledger/artifacts/ && \
    ./peer-identity.sh $CA_ADDRESS_PORT $DOMAIN_OF_ORGANIZATION $HOST_COMPUTER_IP_ADDRESS \
    $PEER_PASSWORD $PEER_NUMBER"
done

# Move and setup crypto materials
sudo mv ./${ORGANIZATION_NAME_LOWERCASE}Ca/client/crypto-config ./
sudo chmod -R 777 ./crypto-config

# Setup orderer TLS certificates
for ORDERER_NUMBER in 1 2 3; do
    ORDERER_DIRECTORY=./crypto-config/ordererOrganizations/orderers
    sudo mv $ORDERER_DIRECTORY/orderer$ORDERER_NUMBER.$DOMAIN_OF_ORGANIZATION/tls/signcerts/cert.pem \
        $ORDERER_DIRECTORY/orderer$ORDERER_NUMBER.$DOMAIN_OF_ORGANIZATION/tls/server.crt
    sudo mv $ORDERER_DIRECTORY/orderer$ORDERER_NUMBER.$DOMAIN_OF_ORGANIZATION/tls/keystore/*_sk \
        $ORDERER_DIRECTORY/orderer$ORDERER_NUMBER.$DOMAIN_OF_ORGANIZATION/tls/server.key
    sudo mv $ORDERER_DIRECTORY/orderer$ORDERER_NUMBER.$DOMAIN_OF_ORGANIZATION/tls/tlscacerts/*.pem \
        $ORDERER_DIRECTORY/orderer$ORDERER_NUMBER.$DOMAIN_OF_ORGANIZATION/tls/ca.crt
    sudo rm -rf $ORDERER_DIRECTORY/orderer$ORDERER_NUMBER.$DOMAIN_OF_ORGANIZATION/tls/{cacerts,keystore,signcerts,tlscacerts,user}
done

# Setup peer TLS certificates
for PEER_NUMBER in 1 2; do
    PEER_DIRECTORY=./crypto-config/peerOrganizations/peers/peer$PEER_NUMBER.$DOMAIN_OF_ORGANIZATION
    sudo mv $PEER_DIRECTORY/tls/signcerts/cert.pem $PEER_DIRECTORY/tls/server.crt
    sudo mv $PEER_DIRECTORY/tls/keystore/*_sk $PEER_DIRECTORY/tls/server.key
    sudo mv $PEER_DIRECTORY/tls/tlscacerts/*.pem $PEER_DIRECTORY/tls/ca.crt
    sudo rm -rf $PEER_DIRECTORY/tls/{cacerts,keystore,signcerts,tlscacerts,user}
done

# Generate channel configuration
./generate.sh idonatechannel $NAME_OF_ORGANIZATION
sleep 2

# Start network components
docker-compose -f $COMPOSE_FILE up -d peer peer2 couchdb cli
sleep 2
docker-compose -f $COMPOSE_FILE up -d orderer2 orderer3

# Join orderers to channel
for ORDERER_NUMBER in 1 2 3; do
    PORT=$((7053 + (ORDERER_NUMBER-1)*1000))
    docker exec cli osnadmin channel join \
        -o orderer${ORDERER_NUMBER}.$DOMAIN_OF_ORGANIZATION:${PORT} \
        --channelID idonatechannel \
        --config-block /etc/hyperledger/artifacts/channel.tx \
        --ca-file /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer${ORDERER_NUMBER}.$DOMAIN_OF_ORGANIZATION/tls/ca.crt \
        --client-cert /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer${ORDERER_NUMBER}.$DOMAIN_OF_ORGANIZATION/tls/server.crt \
        --client-key /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer${ORDERER_NUMBER}.$DOMAIN_OF_ORGANIZATION/tls/server.key
done

sleep 3

# Join peers to channel
docker exec cli peer channel fetch 0 channel.block -c idonatechannel \
    -o orderer1.${DOMAIN_OF_ORGANIZATION}:7050 --tls \
    --cafile /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.${DOMAIN_OF_ORGANIZATION}/tls/ca.crt

docker exec cli peer channel join -b channel.block

# Join peer2 to channel
docker exec -e CORE_PEER_LOCALMSPID="${NAME_OF_ORGANIZATION}MSP" \
    -e CORE_PEER_ADDRESS="peer2.$DOMAIN_OF_ORGANIZATION:7051" \
    -e CORE_PEER_MSPCONFIGPATH="/etc/hyperledger/crypto-config/peerOrganizations/users/Admin@peer2.$DOMAIN_OF_ORGANIZATION/msp" \
    -e CORE_PEER_TLS_ROOTCERT_FILE="/etc/hyperledger/crypto-config/peerOrganizations/peers/peer2.$DOMAIN_OF_ORGANIZATION/tls/ca.crt" \
    cli peer channel fetch 0 channel.block -c idonatechannel \
    -o orderer1.${DOMAIN_OF_ORGANIZATION}:7050 --tls \
    --cafile /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.${DOMAIN_OF_ORGANIZATION}/tls/ca.crt

docker exec -e CORE_PEER_LOCALMSPID="${NAME_OF_ORGANIZATION}MSP" \
    -e CORE_PEER_ADDRESS="peer2.$DOMAIN_OF_ORGANIZATION:7051" \
    -e CORE_PEER_MSPCONFIGPATH="/etc/hyperledger/crypto-config/peerOrganizations/users/Admin@peer2.$DOMAIN_OF_ORGANIZATION/msp" \
    -e CORE_PEER_TLS_ROOTCERT_FILE="/etc/hyperledger/crypto-config/peerOrganizations/peers/peer2.$DOMAIN_OF_ORGANIZATION/tls/ca.crt" \
    cli peer channel join -b channel.block

# Install and approve chaincode
docker exec cli peer lifecycle chaincode package donation.tar.gz --path /etc/hyperledger/chaincode --lang golang --label donation_v1
docker exec cli peer lifecycle chaincode install donation.tar.gz
docker exec cli peer lifecycle chaincode queryinstalled >&log.txt
export PACKAGE_ID=`sed -n '/Package/{s/^Package ID: //; s/, Label:.*$//; p;}' log.txt`

docker exec cli peer lifecycle chaincode approveformyorg \
    -o orderer1.$DOMAIN_OF_ORGANIZATION:7050 \
    --ordererTLSHostnameOverride orderer1.$DOMAIN_OF_ORGANIZATION \
    --channelID idonatechannel --name chaincode --version 1.0 --sequence 1 \
    --tls --cafile /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt \
    --package-id ${PACKAGE_ID}

# Check commit readiness and commit chaincode
docker exec cli peer lifecycle chaincode checkcommitreadiness \
    --channelID idonatechannel --name chaincode --version 1.0 --sequence 1 \
    --tls true --cafile /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt \
    --output json

docker exec cli peer lifecycle chaincode commit \
    -o orderer1.$DOMAIN_OF_ORGANIZATION:7050 \
    --channelID idonatechannel --name chaincode --version 1.0 --sequence 1 \
    --tls true --cafile /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt \
    --peerAddresses peer1.$DOMAIN_OF_ORGANIZATION:7051 \
    --tlsRootCertFiles /etc/hyperledger/crypto-config/peerOrganizations/peers/peer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt


docker exec cli peer chaincode invoke -o orderer1.$DOMAIN_OF_ORGANIZATION:7050 --tls --cafile "/etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt" -C idonatechannel -n chaincode --peerAddresses peer1.$DOMAIN_OF_ORGANIZATION:7051 --tlsRootCertFiles "/etc/hyperledger/crypto-config/peerOrganizations/peers/peer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt" -c '{"function":"CreateDonationEvent","Args":[
    "kevinEventOK9999911223344556677889",
        "Education form Org334 ",
        "Kay Kang",
        "The Pupose of creating this event is to motivite young generate to work harder!",
        "ISTAD Organization"
]}'


docker exec cli peer chaincode invoke -o orderer1.$DOMAIN_OF_ORGANIZATION:7050 --tls --cafile "/etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt" -C idonatechannel -n chaincode --peerAddresses peer1.$DOMAIN_OF_ORGANIZATION:7051 --tlsRootCertFiles "/etc/hyperledger/crypto-config/peerOrganizations/peers/peer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'

sleep 2 

docker exec cli peer chaincode query -o orderer1.$DOMAIN_OF_ORGANIZATION:7050 -C idonatechannel -n chaincode -c '{"function":"ReadDonationEvent","Args":["kevinEventOK9999911223344556677889"]}' --tls --cafile /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt | jq .

# Query Installed chaincode on peer
docker exec cli peer lifecycle chaincode queryinstalled --peerAddresses peer1.$DOMAIN_OF_ORGANIZATION:7051 --tlsRootCertFiles /etc/hyperledger/crypto-config/peerOrganizations/peers/peer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt

# Query commited chaincode on the channel
docker exec cli peer lifecycle chaincode querycommitted -o orderer.$DOMAIN_OF_ORGANIZATION:7050 --channelID idonatechannel --tls --cafile /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt --peerAddresses peer1.$DOMAIN_OF_ORGANIZATION:7051 --tlsRootCertFiles /etc/hyperledger/crypto-config/peerOrganizations/peers/peer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt

echo "Deployment completed for $NAME_OF_ORGANIZATION"