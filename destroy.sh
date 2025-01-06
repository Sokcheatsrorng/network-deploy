# SPDX-License-Identifier: Apache-2.0
# Stop the containers instead of removing them
# CONTAINER_IDS=$(docker ps -aq)
# if [ -z "$CONTAINER_IDS" ] || [ "$CONTAINER_IDS" == " " ]; then
#   echo "---- No containers available to stop ----"
# else
#   echo "Stopping existing containers..."
#   docker stop $CONTAINER_IDS
# fi

# Clean the container environment for deployment
CONTAINER_IDS=$(docker ps -aq)
if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" == " " ]; then
  echo "---- No containers available for deletion ----"
else
  docker rm -f $CONTAINER_IDS
fi

DOCKER_IMAGE_IDS=$(docker images | awk '$1 ~ /dev-peer*/ {print $3}')
  if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
    echo "---- No images available for deletion ----"
  else
    docker rmi -f $DOCKER_IMAGE_IDS
fi

# Delete old certificates and channel artifacts
sudo rm -rf ${ORGANIZATION_NAME_LOWERCASE}Ca/ crypto-config/ artifacts/ config/ chaincode/node_modules configtx.yaml

if [ ! -d artifacts ]; then
   mkdir artifacts
fi

if [ ! -d config ];then
   mkdir config
fi
