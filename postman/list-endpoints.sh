#!/bin/bash

# List all endpoints in the Postman collection

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COLLECTION_FILE="$SCRIPT_DIR/algoholic-api.postman_collection.json"

# Colors
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Algoholic API - Endpoint List${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Extract collection info
COLLECTION_NAME=$(jq -r '.info.name' "$COLLECTION_FILE")
COLLECTION_DESC=$(jq -r '.info.description' "$COLLECTION_FILE")
FOLDER_COUNT=$(jq '.item | length' "$COLLECTION_FILE")

echo -e "${GREEN}Collection:${NC} $COLLECTION_NAME"
echo -e "${GREEN}Description:${NC} $COLLECTION_DESC"
echo -e "${GREEN}Folders:${NC} $FOLDER_COUNT"
echo ""

# Count total requests
TOTAL_REQUESTS=$(jq '[.item[].item | length] | add' "$COLLECTION_FILE")
echo -e "${YELLOW}Total Endpoints: $TOTAL_REQUESTS${NC}"
echo ""

# List all folders and their requests
jq -r '.item[] |
  "→ \(.name) (\(.item | length) endpoints)\n" +
  (.item[] | "  • [\(.request.method)] \(.name)")' "$COLLECTION_FILE"

echo ""
echo -e "${BLUE}========================================${NC}"
