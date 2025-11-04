set -e  # Exit on any error


# Setup private network
# Instruction from here - https://timeweb.cloud/docs/vpc/managing-bgp-networks/configuring-private-ip#ubuntu

# Configuration variables
INTERFACE_NAME="eth1"
PRIVATE_IP="192.168.0.12/24"
NETPLAN_FILE="/etc/netplan/99-ipv4.yaml"

# Step 1: Configure private IP through Netplan
echo "Configuring Netplan for interface $INTERFACE_NAME..."

cat > "$NETPLAN_FILE" << EOF
network:
  version: 2
  renderer: networkd
  ethernets:
    $INTERFACE_NAME:
      dhcp4: false
      addresses:
        - "$PRIVATE_IP"
EOF

chmod 600 "$NETPLAN_FILE"
netplan apply

# Step 2: Configure non-local binding
echo "Configuring non-local binding..."

echo "net.ipv4.ip_nonlocal_bind = 1" >> /etc/sysctl.conf
echo "net.ipv6.ip_nonlocal_bind = 1" >> /etc/sysctl.conf
sysctl -p /etc/sysctl.conf

# Step 3: Restart network service
echo "Restarting network service..."
systemctl restart systemd-networkd

# Step 4: Verify configuration
echo "Network configuration completed for $INTERFACE_NAME with IP $PRIVATE_IP"
ip a show "$INTERFACE_NAME"


# Install Docker

# Install GitHub CLI
