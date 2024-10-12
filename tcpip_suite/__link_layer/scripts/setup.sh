#!/bin/bash
# scripts/setup.sh

# Exit immediately if a command exits with a non-zero status.
set -e

# Update and install necessary packages.
echo "Updating package list..."
sudo apt-get update

echo "Installing necessary packages..."
sudo apt-get install -y build-essential libpcap-dev iproute2

# Build the application.
echo "Building the Link Layer application..."
cd /path/to/link_layer
go build -o link_layer cmd/main.go

# Configure network interface (requires root privileges).
echo "Configuring network interface..."
sudo ifconfig eth0 up
sudo ifconfig eth0 mtu 1500

# Run the application.
echo "Running the Link Layer application..."
sudo ./link_layer -config configs/config.yaml