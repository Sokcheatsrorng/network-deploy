
# Check if the script is run as root (required to modify /etc/hosts)
if [ "$(id -u)" -ne 0 ]; then
  echo "You must run this script as root or with sudo."
  exit 1
fi

# Check for the correct number of arguments
if [ $# -ne 2 ]; then
  echo "Usage: $0 <IP_ADDRESS> <ORG_NAME>"
  exit 1
fi

# Variables from the arguments
IP_ADDRESS=$1
ORG_NAME=$2
HOST_NAME="ca.${ORG_NAME}.idonate.istad.co"

# Backup the current /etc/hosts file
cp /etc/hosts /etc/hosts.bak

# Add the new host entry to /etc/hosts
echo "${IP_ADDRESS} ${HOST_NAME}" >> /etc/hosts

# Confirm the change
echo "Added ${HOST_NAME} -> ${IP_ADDRESS} to /etc/hosts"

# Optionally, you can also display the new line in /etc/hosts
grep "${HOST_NAME}" /etc/hosts