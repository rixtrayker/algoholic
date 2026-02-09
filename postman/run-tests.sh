#!/bin/bash

# Algoholic API - Newman Test Runner
# This script runs the Postman collection using Newman

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COLLECTION_FILE="$SCRIPT_DIR/algoholic-api.postman_collection.json"
ENVIRONMENT_FILE="$SCRIPT_DIR/algoholic-local.postman_environment.json"
REPORT_DIR="$SCRIPT_DIR/reports"

# Create reports directory if it doesn't exist
mkdir -p "$REPORT_DIR"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Algoholic API Test Suite${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Check if backend is running
echo -e "${BLUE}Checking if backend is running...${NC}"
if curl -s http://localhost:4000/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Backend is running${NC}"
else
    echo -e "${RED}✗ Backend is not running!${NC}"
    echo -e "${RED}Please start the backend with: cd backend && go run main.go${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}Running tests...${NC}"
echo ""

# Run Newman with detailed output
newman run "$COLLECTION_FILE" \
    --environment "$ENVIRONMENT_FILE" \
    --reporters cli,json,junit \
    --reporter-json-export "$REPORT_DIR/newman-report.json" \
    --reporter-junit-export "$REPORT_DIR/newman-report.xml" \
    --bail \
    --color on \
    --timeout-request 10000

EXIT_CODE=$?

echo ""
if [ $EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  All tests passed! ✓${NC}"
    echo -e "${GREEN}========================================${NC}"
else
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}  Some tests failed! ✗${NC}"
    echo -e "${RED}========================================${NC}"
fi

echo ""
echo -e "${BLUE}Reports saved to: $REPORT_DIR${NC}"
echo ""

exit $EXIT_CODE
