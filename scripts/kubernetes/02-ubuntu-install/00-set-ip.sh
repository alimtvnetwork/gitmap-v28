#!/bin/bash

# Check if at least one argument is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <ip> [route ip]"
  exit 1
fi

IP=$1
ROUTE_IP="${2:-192.168.0.1}"

# Remove the existing netplan configuration file if it exists
if [ -f /etc/netplan/00-installer-config.yaml ]; then
  sudo rm /etc/netplan/00-installer-config.yaml
fi

# Create the netplan configuration directory if it doesn't exist
sudo mkdir -p /etc/netplan

echo "Route IP : $ROUTE_IP"

# Create a new netplan configuration file
sudo bash -c "cat > /etc/netplan/00-installer-config.yaml <<EOL
network:
  renderer: networkd
  ethernets:
    ens33:
      dhcp4: false
      addresses:
        - $IP/24
      routes:
        - to: default
          via: $ROUTE_IP
EOL"

# Add the nameservers
sudo bash -c "cat >> /etc/netplan/00-installer-config.yaml <<EOL
      nameservers:
          addresses: [8.8.8.8, 1.1.1.1]
  version: 2
EOL"

cat /etc/netplan/00-installer-config.yaml
# Validate and apply the netplan configuration
sudo netplan try && sudo netplan apply

sleep 2
ping google.com