# SPDX-License-Identifier: Apache-2.0
export PATH=$GOPATH/src/github.com/hyperledger/fabric/build/bin:${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}
CHANNEL_NAME=$1
ORGANIZATION_NAME=$2

# Create config and crypto-config if not exists
mkdir -p config/

# Remove previous crypto material and config transactions
rm -fr config/*

# Generate channel configuration transaction
configtxgen -profile Channel -outputBlock ./config/channel.tx -channelID $CHANNEL_NAME
if [ "$?" -ne 0 ]; then
  echo "Failed to generate channel configuration transaction..."
  exit 1
fi

# Generate anchor peer transaction
configtxgen -profile Channel -outputAnchorPeersUpdate ./config/${ORGANIZATION_NAME}MSPanchors.tx -channelID $CHANNEL_NAME -asOrg ${ORGANIZATION_NAME}MSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for MSP..."
  exit 1
fi

# # Create a directory to store connection profiles for each organization
mkdir -p connection-profiles

# Update variables in template for each organization and store the connection profiles
sed -e 's/NAME_OF_ORGANIZATION/'$NAME_OF_ORGANIZATION'/g' \
    -e 's/DOMAIN_OF_ORGANIZATION/'$DOMAIN_OF_ORGANIZATION'/g' \
    -e 's/HOST_COMPUTER_IP_ADDRESS/'$HOST_COMPUTER_IP_ADDRESS'/g' \
    -e 's/ORGANIZATION_NAME_LOWERCASE/'$ORGANIZATION_NAME_LOWERCASE'/g' \
    ./template/connection-org.json > ./connection-profiles/connection-${ORGANIZATION_NAME}.json

# Update for connections.yml template
sed -e 's/NAME_OF_ORGANIZATION/'$NAME_OF_ORGANIZATION'/g' \
    -e 's/DOMAIN_OF_ORGANIZATION/'$DOMAIN_OF_ORGANIZATION'/g' \
    -e 's/HOST_COMPUTER_IP_ADDRESS/'$HOST_COMPUTER_IP_ADDRESS'/g' \
    -e 's/ORGANIZATION_NAME_LOWERCASE/'$ORGANIZATION_NAME_LOWERCASE'/g' \
    ./template/connections.yml > ./connection-profiles/connections.yml

# Optionally, if you have multiple organizations, repeat the above steps for each organization
# For example, you can loop through an array of organization names and generate profiles for each
# Example:
# ORGANIZATION_NAMES=("Org1" "Org2" "Org3")
# for org in "${ORGANIZATION_NAMES[@]}"; do
#   # Repeat the sed and connection file generation logic for each organization
# done

echo "Connection profiles for ${ORGANIZATION_NAME} have been stored in ./connection-profiles"
