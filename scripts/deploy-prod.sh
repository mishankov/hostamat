source ./scripts/params.sh

# Run deploy script
ssh -t ${PROD} 'cd ~/hostamat/infra/prod; ./scripts/deploy-app.sh'
