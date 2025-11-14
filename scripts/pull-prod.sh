source ./scripts/params.sh

# Pull from git
ssh -t ${PROD} 'cd ~/hostamat; git pull'
