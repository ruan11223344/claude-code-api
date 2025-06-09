#!/bin/bash

# Claude Code API runner script with log management

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Create logs directory if it doesn't exist
mkdir -p logs

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
    echo -e "${GREEN}✓ Loaded environment variables from .env${NC}"
else
    echo -e "${YELLOW}⚠ No .env file found, using defaults${NC}"
fi

# Check if go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go is not installed${NC}"
    exit 1
fi

# Build the application
echo -e "${GREEN}Building application...${NC}"
go build -o claude-code-api main.go
if [ $? -ne 0 ]; then
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

# Get current date for log file
LOG_DATE=$(date +%Y-%m-%d)
LOG_FILE="logs/server_${LOG_DATE}.log"

echo -e "${GREEN}✓ Build successful${NC}"
echo -e "${GREEN}Starting Claude Code API Server...${NC}"
echo -e "${YELLOW}Log file: ${LOG_FILE}${NC}"
echo "====================================="

# Run the application
if [ "$LOG_TO_FILE" = "true" ]; then
    ./claude-code-api "$@"
else
    # If LOG_TO_FILE is not set, we can still capture logs
    ./claude-code-api "$@" 2>&1 | tee -a "$LOG_FILE"
fi