#!/bin/bash

# Notification/Review/Shipping Integration Test Runner
# This script runs all integration tests for the notification flow

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Notification/Review/Shipping Flow Tests${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Check if services are running
echo -e "${YELLOW}Checking service availability...${NC}"

# Check Notification Service
if ! nc -z localhost 22017 2>/dev/null; then
    echo -e "${RED}❌ Notification database (port 22017) is not running${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Notification database is running${NC}"

# Check Review Service
if ! nc -z localhost 22018 2>/dev/null; then
    echo -e "${RED}❌ Review database (port 22018) is not running${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Review database is running${NC}"

# Check Shipping Service
if ! nc -z localhost 22016 2>/dev/null; then
    echo -e "${RED}❌ Shipping database (port 22016) is not running${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Shipping database is running${NC}"

echo ""

# Install dependencies
echo -e "${YELLOW}Installing dependencies...${NC}"
go mod download
echo -e "${GREEN}✅ Dependencies installed${NC}"
echo ""

# Run tests
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Running Tests${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Run notification tests
echo -e "${YELLOW}Running Notification Tests...${NC}"
go test -v -run TestOrderConfirmationNotification
go test -v -run TestPaymentSuccessNotification
go test -v -run TestShippingUpdateNotification
go test -v -run TestNotificationFailureHandling
go test -v -run TestMultipleNotificationsForRecipient
echo ""

# Run review tests
echo -e "${YELLOW}Running Review Tests...${NC}"
go test -v -run TestCreateReview
go test -v -run TestReviewRatingValidation
go test -v -run TestGetReviewsByProduct
go test -v -run TestAverageRating
go test -v -run TestUpdateReview
go test -v -run TestReviewEditableUntil
go test -v -run TestDuplicateReviewPrevention
echo ""

# Run shipping tests
echo -e "${YELLOW}Running Shipping Tests...${NC}"
go test -v -run TestCreateShipment
go test -v -run TestShipmentStatusFlow
go test -v -run TestTrackingEvents
go test -v -run TestShippingWithNotification
go test -v -run TestMultipleTrackingEventsChronology
go test -v -run TestShipmentUniqueOrderID
echo ""

# Run all tests with coverage
echo -e "${YELLOW}Running All Tests with Coverage...${NC}"
go test -v -coverprofile=coverage.out ./...

# Generate coverage report
if [ -f coverage.out ]; then
    echo ""
    echo -e "${YELLOW}Coverage Summary:${NC}"
    go tool cover -func=coverage.out | tail -n 1
    echo ""
    echo -e "${GREEN}Coverage report generated: coverage.out${NC}"
    echo -e "${GREEN}View HTML report with: go tool cover -html=coverage.out${NC}"
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✅ All Tests Completed${NC}"
echo -e "${BLUE}========================================${NC}"
