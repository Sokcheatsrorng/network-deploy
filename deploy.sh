# SPDX-License-Identifier: Apache-2.0
#set -ev
#!/bin/sh

function exportVariables(){

  # Organization information that you wish to build and deploy
  export NAME_OF_ORGANIZATION=$NAME_OF_ORGANIZATION
  export DOMAIN_OF_ORGANIZATION=$DOMAIN_OF_ORGANIZATION
  export HOST_COMPUTER_IP_ADDRESS=$HOST_COMPUTER_IP_ADDRESS
  export ORGANIZATION_NAME_LOWERCASE=`echo "$NAME_OF_ORGANIZATION" | tr '[:upper:]' '[:lower:]'`
 

  # Security defaults
  # Couch DB credentials
  export COUCH_DB_USERNAME=admin
  export COUCH_DB_PASSWORD=adminpw

  # Certificate authority credentials
  export CA_ADMIN_USER=admin
  export CA_ADMIN_PASSWORD=adminpw

  # Orderer credentials
  ORDERER_PASSWORD=adminpw

  # Peer credentials
  PEER_PASSWORD=peerpw

}

read -p "Organization Name: "  NAME_OF_ORGANIZATION
# read -p "Organization Domain: " DOMAIN_OF_ORGANIZATION
# read -p "Computer IP Address: " HOST_COMPUTER_IP_ADDRESS
HOST_COMPUTER_IP_ADDRESS=172.17.0.1
DOMAIN_OF_ORGANIZATION=$NAME_OF_ORGANIZATION.idonate.istad.co
export CA_ADDRESS_PORT=ca.$DOMAIN_OF_ORGANIZATION:7054
exportVariables

./clean-all.sh

# Substitutes organizations information in the configtx template to match organizations name, domain and ip address
sed -e 's/organization_name/'$NAME_OF_ORGANIZATION'/g' -e 's/organization_domain/'$DOMAIN_OF_ORGANIZATION'/g' -e 's/ip_address/'$HOST_COMPUTER_IP_ADDRESS'/g'  configtx_template.yaml > configtx.yaml

# Start the certficate authority
docker-compose -p fabric-network -f docker-compose.yml up -d ca
sleep 3

# Generate identity and cryptographic materials for the 3 orderers 
for ORDERER_NUMBER in 1 2 3
do
  docker exec ca.$DOMAIN_OF_ORGANIZATION /bin/bash -c "cd /etc/hyperledger/artifacts/  && ./orderer-identity.sh $CA_ADDRESS_PORT $DOMAIN_OF_ORGANIZATION $HOST_COMPUTER_IP_ADDRESS $CA_ADMIN_USER $CA_ADMIN_PASSWORD $ORDERER_NUMBER $ORDERER_PASSWORD"
done


# Generate identity and cryptographic materials for the peer 
for PEER_NUMBER in 1 2 
do
  docker exec ca.$DOMAIN_OF_ORGANIZATION /bin/bash -c "cd /etc/hyperledger/artifacts/  && ./peer-identity.sh $CA_ADDRESS_PORT $DOMAIN_OF_ORGANIZATION $HOST_COMPUTER_IP_ADDRESS $PEER_PASSWORD $PEER_NUMBER"
done



# Move the crypto-config folder to manipulate it more easily away from the dockers users' restrictions
sudo mv ./${ORGANIZATION_NAME_LOWERCASE}Ca/client/crypto-config ./
sudo chmod -R 777 ./crypto-config

# Move TLS certificates for the 3 orderers 
for ORDERER_NUMBER in 1 2 3
do
  ORDERER_DIRECTORY=./crypto-config/ordererOrganizations/orderers
  sudo mv $ORDERER_DIRECTORY/orderer$ORDERER_NUMBER.$DOMAIN_OF_ORGANIZATION/tls/signcerts/cert.pem $ORDERER_DIRECTORY/orderer$ORDERER_NUMBER.$DOMAIN_OF_ORGANIZATION/tls/server.crt
  sudo mv $ORDERER_DIRECTORY/orderer$ORDERER_NUMBER.$DOMAIN_OF_ORGANIZATION/tls/keystore/*_sk $ORDERER_DIRECTORY/orderer$ORDERER_NUMBER.$DOMAIN_OF_ORGANIZATION/tls/server.key
  sudo mv $ORDERER_DIRECTORY/orderer$ORDERER_NUMBER.$DOMAIN_OF_ORGANIZATION/tls/tlscacerts/*.pem $ORDERER_DIRECTORY/orderer$ORDERER_NUMBER.$DOMAIN_OF_ORGANIZATION/tls/ca.crt

  # Delete empty directories
  sudo rm -rf $ORDERER_DIRECTORY/orderer$ORDERER_NUMBER.$DOMAIN_OF_ORGANIZATION/tls/{cacerts,keystore,signcerts,tlscacerts,user}
done

# Move TLS certificates for the 2 peers 
for PEER_NUMBER in 1 2
do
  # Peer 1 crypto-config directory
  PEER_DIRECTORY=./crypto-config/peerOrganizations/peers/peer$PEER_NUMBER.$DOMAIN_OF_ORGANIZATION

  # Move the Peer TLS files to match cryptogen hierarchy
  sudo mv $PEER_DIRECTORY/tls/signcerts/cert.pem $PEER_DIRECTORY/tls/server.crt
  sudo mv $PEER_DIRECTORY/tls/keystore/*_sk $PEER_DIRECTORY/tls/server.key
  sudo mv $PEER_DIRECTORY/tls/tlscacerts/*.pem $PEER_DIRECTORY/tls/ca.crt

  # Delete the peers empty directory
  sudo rm -rf $PEER_DIRECTORY/tls/{cacerts,keystore,signcerts,tlscacerts,user}
done

# Generate the channel configuration 
./generate.sh idonatechannel $NAME_OF_ORGANIZATION
sleep 2

# Start the network with docker-compose
docker-compose -f docker-compose.yml up -d peer peer2 couchdb cli 
sleep 2

docker-compose -f docker-compose.yml up -d orderer2 
docker-compose -f docker-compose.yml up -d orderer3

# Join the orders to the channel
docker exec cli osnadmin channel join -o orderer1.$DOMAIN_OF_ORGANIZATION:7053 --channelID idonatechannel --config-block /etc/hyperledger/artifacts/channel.tx --ca-file /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt --client-cert /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.$DOMAIN_OF_ORGANIZATION/tls/server.crt --client-key /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.$DOMAIN_OF_ORGANIZATION/tls/server.key 

docker exec cli osnadmin channel join -o orderer2.$DOMAIN_OF_ORGANIZATION:8053 --channelID idonatechannel --config-block /etc/hyperledger/artifacts/channel.tx --ca-file /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer2.$DOMAIN_OF_ORGANIZATION/tls/ca.crt --client-cert /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer2.$DOMAIN_OF_ORGANIZATION/tls/server.crt --client-key /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer2.$DOMAIN_OF_ORGANIZATION/tls/server.key 

docker exec cli osnadmin channel join -o orderer3.$DOMAIN_OF_ORGANIZATION:9053 --channelID idonatechannel --config-block /etc/hyperledger/artifacts/channel.tx --ca-file /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer3.$DOMAIN_OF_ORGANIZATION/tls/ca.crt --client-cert /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer3.$DOMAIN_OF_ORGANIZATION/tls/server.crt --client-key /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer3.$DOMAIN_OF_ORGANIZATION/tls/server.key 

sleep 3

docker exec cli peer channel fetch 0 channel.block -c idonatechannel -o orderer1.${DOMAIN_OF_ORGANIZATION}:7050 --tls --cafile /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.${DOMAIN_OF_ORGANIZATION}/tls/ca.crt

sleep 1

docker exec cli peer channel join -b channel.block

docker exec -e CORE_PEER_LOCALMSPID="${NAME_OF_ORGANIZATION}MSP" -e CORE_PEER_ADDRESS="peer2.$DOMAIN_OF_ORGANIZATION:7051" -e CORE_PEER_MSPCONFIGPATH="/etc/hyperledger/crypto-config/peerOrganizations/users/Admin@peer2.$DOMAIN_OF_ORGANIZATION/msp" -e CORE_PEER_TLS_ROOTCERT_FILE="/etc/hyperledger/crypto-config/peerOrganizations/peers/peer2.$DOMAIN_OF_ORGANIZATION/tls/ca.crt" cli peer channel fetch 0 channel.block -c idonatechannel -o orderer1.${DOMAIN_OF_ORGANIZATION}:7050 --tls --cafile /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.${DOMAIN_OF_ORGANIZATION}/tls/ca.crt

docker exec -e CORE_PEER_LOCALMSPID="${NAME_OF_ORGANIZATION}MSP" -e CORE_PEER_ADDRESS="peer2.$DOMAIN_OF_ORGANIZATION:7051" -e CORE_PEER_MSPCONFIGPATH="/etc/hyperledger/crypto-config/peerOrganizations/users/Admin@peer2.$DOMAIN_OF_ORGANIZATION/msp" -e CORE_PEER_TLS_ROOTCERT_FILE="/etc/hyperledger/crypto-config/peerOrganizations/peers/peer2.$DOMAIN_OF_ORGANIZATION/tls/ca.crt" cli peer channel join -b channel.block

# # Package chaincode
docker exec cli peer lifecycle chaincode package donation.tar.gz --path /etc/hyperledger/chaincode --lang golang --label donation_v1

# Install the chaincode
docker exec cli peer lifecycle chaincode install donation.tar.gz

docker exec cli peer lifecycle chaincode queryinstalled >&log.txt
export PACKAGE_ID=`sed -n '/Package/{s/^Package ID: //; s/, Label:.*$//; p;}' log.txt`

echo $PACKAGE_ID

# Approve chaincode for org
docker exec cli peer lifecycle chaincode approveformyorg -o orderer1.$DOMAIN_OF_ORGANIZATION:7050 --ordererTLSHostnameOverride orderer1.$DOMAIN_OF_ORGANIZATION --channelID idonatechannel --name chaincode --version 1.0 --sequence 1 --tls --cafile /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt --package-id ${PACKAGE_ID}

# Check commit readiness
docker exec cli peer lifecycle chaincode checkcommitreadiness --channelID idonatechannel --name chaincode --version 1.0 --sequence 1 --tls true --cafile /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt --output json

# Commit the chaincode
docker exec cli peer lifecycle chaincode commit -o orderer1.$DOMAIN_OF_ORGANIZATION:7050 --channelID idonatechannel --name chaincode --version 1.0 --sequence 1 --tls true --cafile /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt --peerAddresses peer1.$DOMAIN_OF_ORGANIZATION:7051 --tlsRootCertFiles /etc/hyperledger/crypto-config/peerOrganizations/peers/peer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt 

# docker exec cli peer chaincode invoke -o orderer1.$DOMAIN_OF_ORGANIZATION:7050 -C idonatechannel -n donation -c '{"Args":["invokeTransaction","1","{anythingHereAsJsonPayload}"]}' --tls --cafile /etc/hyperledger/crypto-config/ordererOrganizations/orderers/orderer1.$DOMAIN_OF_ORGANIZATION/tls/ca.crt


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

echo NETWORK DEPLOYMENT COMPLETED SUCCESSFULLY