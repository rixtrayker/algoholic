#!/bin/bash

# Demo script to show Newman usage
# This demonstrates running the Postman collection with various options

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COLLECTION_FILE="$SCRIPT_DIR/algoholic-api.postman_collection.json"
ENVIRONMENT_FILE="$SCRIPT_DIR/algoholic-local.postman_environment.json"

# Colors
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Newman Demo - Algoholic API${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

echo -e "${YELLOW}Note: This is a demo of Newman commands.${NC}"
echo -e "${YELLOW}To actually run tests, ensure the backend is running first.${NC}"
echo ""

# Show available commands
echo -e "${GREEN}1. Basic Run:${NC}"
echo "   newman run $COLLECTION_FILE \\"
echo "     --environment $ENVIRONMENT_FILE"
echo ""

echo -e "${GREEN}2. Run with Verbose Output:${NC}"
echo "   newman run $COLLECTION_FILE \\"
echo "     --environment $ENVIRONMENT_FILE \\"
echo "     --verbose"
echo ""

echo -e "${GREEN}3. Run Specific Folder (e.g., Authentication):${NC}"
echo "   newman run $COLLECTION_FILE \\"
echo "     --environment $ENVIRONMENT_FILE \\"
echo "     --folder \"Authentication\""
echo ""

echo -e "${GREEN}4. Generate JSON Report:${NC}"
echo "   newman run $COLLECTION_FILE \\"
echo "     --environment $ENVIRONMENT_FILE \\"
echo "     --reporters cli,json \\"
echo "     --reporter-json-export reports/newman-report.json"
echo ""

echo -e "${GREEN}5. Run with Delay Between Requests:${NC}"
echo "   newman run $COLLECTION_FILE \\"
echo "     --environment $ENVIRONMENT_FILE \\"
echo "     --delay-request 500"
echo ""

echo -e "${GREEN}6. Run Multiple Iterations:${NC}"
echo "   newman run $COLLECTION_FILE \\"
echo "     --environment $ENVIRONMENT_FILE \\"
echo "     --iteration-count 3"
echo ""

echo -e "${GREEN}7. Run with Custom Timeout:${NC}"
echo "   newman run $COLLECTION_FILE \\"
echo "     --environment $ENVIRONMENT_FILE \\"
echo "     --timeout-request 15000"
echo ""

echo -e "${GREEN}8. Full Production Run with All Reports:${NC}"
echo "   newman run $COLLECTION_FILE \\"
echo "     --environment $ENVIRONMENT_FILE \\"
echo "     --reporters cli,json,junit \\"
echo "     --reporter-json-export reports/newman-report.json \\"
echo "     --reporter-junit-export reports/newman-report.xml \\"
echo "     --bail \\"
echo "     --color on"
echo ""

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}To run tests, use: ./run-tests.sh${NC}"
echo -e "${BLUE}========================================${NC}"
