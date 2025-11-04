set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${PURPLE}[STEP]${NC} $1"
}

print_header() {
    echo -e "${CYAN}=== $1 ===${NC}"
}

print_header "Starting Ubuntu Server Setup"

print_step "Updating system packages..."
apt update
apt upgrade -y

print_header "Private Network Setup"
print_info "Instruction from: https://timeweb.cloud/docs/vpc/managing-bgp-networks/configuring-private-ip#ubuntu"

# Configuration variables
INTERFACE_NAME="eth1"
PRIVATE_IP="192.168.0.12/24"
NETPLAN_FILE="/etc/netplan/99-ipv4.yaml"

# Step 1: Configure private IP through Netplan
print_step "Configuring Netplan for interface $INTERFACE_NAME..."

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
print_step "Configuring non-local binding..."

echo "net.ipv4.ip_nonlocal_bind = 1" >> /etc/sysctl.conf
echo "net.ipv6.ip_nonlocal_bind = 1" >> /etc/sysctl.conf
sysctl -p /etc/sysctl.conf

# Step 3: Restart network service
print_step "Restarting network service..."
systemctl restart systemd-networkd

# Step 4: Verify configuration
print_success "Network configuration completed for $INTERFACE_NAME with IP $PRIVATE_IP"
print_info "Current network interface status:"
ip a show "$INTERFACE_NAME"


print_header "Docker Installation"
print_info "Instruction from: https://docs.docker.com/engine/install/ubuntu/"

print_step "Installing Docker..."

# Set up Docker's apt repository
print_step "Setting up Docker's apt repository..."
apt-get update
apt-get install -y ca-certificates curl
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources
print_step "Adding Docker repository to apt sources..."
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null

apt-get update

# Install Docker packages
print_step "Installing Docker packages..."
apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

print_success "Docker installation completed successfully!"

# Verify Docker installation
print_step "Verifying Docker installation..."
docker run hello-world

print_header "GitHub CLI Installation"
print_info "Instruction from: https://github.com/cli/cli/blob/trunk/docs/install_linux.md#debian"

print_step "Installing GitHub CLI..."

# Install GitHub CLI using official Debian/Ubuntu method
print_step "Setting up GitHub CLI repository and installing..."
(type -p wget >/dev/null || (apt update && apt install wget -y)) \
	&& mkdir -p -m 755 /etc/apt/keyrings \
	&& out=$(mktemp) && wget -nv -O$out https://cli.github.com/packages/githubcli-archive-keyring.gpg \
	&& cat $out | tee /etc/apt/keyrings/githubcli-archive-keyring.gpg > /dev/null \
	&& chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg \
	&& mkdir -p -m 755 /etc/apt/sources.list.d \
	&& echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
	&& apt update \
	&& apt install gh -y

print_success "GitHub CLI installation completed successfully!"
print_info "You can now authenticate with: ${YELLOW}gh auth login${NC}"

print_header "Setup Complete!"
print_success "All components have been successfully installed:"
print_info "✓ Private network configured on $INTERFACE_NAME"
print_info "✓ Docker Engine installed and verified"
print_info "✓ GitHub CLI (gh) installed"
print_warning "Note: You may need to log out and back in for Docker group changes to take effect"
