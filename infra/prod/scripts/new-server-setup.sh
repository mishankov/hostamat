set -e  # Exit on any error

apt update
apt upgrade -y


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
# Instruction from here - https://docs.docker.com/engine/install/ubuntu/

echo "Installing Docker..."

# Set up Docker's apt repository
apt-get update
apt-get install -y ca-certificates curl
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null

apt-get update

# Install Docker packages
apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

echo "Docker installation completed successfully!"

# Verify Docker installation
echo "Verifying Docker installation..."
docker run hello-world

# Install GitHub CLI
# Instruction from here - https://github.com/cli/cli/blob/trunk/docs/install_linux.md#debian

echo "Installing GitHub CLI..."

# Install GitHub CLI using official Debian/Ubuntu method
(type -p wget >/dev/null || (apt update && apt install wget -y)) \
	&& mkdir -p -m 755 /etc/apt/keyrings \
	&& out=$(mktemp) && wget -nv -O$out https://cli.github.com/packages/githubcli-archive-keyring.gpg \
	&& cat $out | tee /etc/apt/keyrings/githubcli-archive-keyring.gpg > /dev/null \
	&& chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg \
	&& mkdir -p -m 755 /etc/apt/sources.list.d \
	&& echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
	&& apt update \
	&& apt install gh -y

echo "GitHub CLI installation completed successfully!"
echo "You can now authenticate with: gh auth login"
