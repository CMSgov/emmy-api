#!/bin/bash

# Emmy API End-to-End Test (Newman-based script)
#
# This script performs the "Obtain Token" and "Get Enrollment Data" flow
# by using Newman to run the bundled Postman collection (emmy.yaml).

set -e

# Configuration (Defaults)
CLIENT_ID="${CLIENT_ID:-6lqtiagq4o3n5bsa3s5mad5h7l}"
AUTH_BASE="${AUTH_BASE:-https://emmy-uat.auth.us-east-1.amazoncognito.com/oauth2/token}"
API_BASE="${API_BASE:-https://api.uat.emmy.cms.gov}"

# Paths to collection
COLLECTION_YAML="examples/postman/emmy.yaml"
COLLECTION_JSON="examples/postman/emmy.json"
echo "Starting E2E Test Flow (Using Newman and emmy.yaml)..."
echo "----------------------------------------------------"

# Ensure Newman is installed
if ! command -v newman &> /dev/null; then
  echo "Error: Newman is not installed or not in PATH. Please install it with 'npm install -g newman'."
  exit 1
fi

# Sync YAML to JSON (Newman requires JSON for many features/versions)
if [ -f "$COLLECTION_YAML" ]; then
  node -e "const yaml = require('js-yaml'); const fs = require('fs'); const doc = yaml.load(fs.readFileSync('$COLLECTION_YAML', 'utf8')); console.log(JSON.stringify(doc, null, 2));" > "$COLLECTION_JSON"
fi

# Execute Newman run
# We pass the variables as collection variables using --env-var
# --insecure is used to bypass self-signed certificate issues in non-prod environments
newman run "$COLLECTION_JSON" \
  --env-var "client_id=$CLIENT_ID" \
  --env-var "client_secret=$SECRET_ID" \
  --env-var "auth_base=$AUTH_BASE" \
  --env-var "api_base=$API_BASE" \
  --reporters cli \
  --insecure

echo -e "\n----------------------------------------------------"
echo "E2E Test completed successfully using Newman."
