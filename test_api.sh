#!/bin/bash

# ================================
# Country API Endpoint Tester
# Works with Go (Gin + GORM) version
# ================================

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="${1:-http://localhost:3000}"
TIMEOUT=30
SERVER_PID=""

# Test results
PASSED=0
FAILED=0

log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}PASS${NC} $1"
    ((PASSED++))
}

fail() {
    echo -e "${RED}FAIL${NC} $1"
    ((FAILED++))
}

check_status() {
    local status=$1
    local expected=$2
    local name=$3
    if [ "$status" -eq "$expected" ]; then
        success "$name"
    else
        fail "$name (got $status, expected $expected)"
    fi
}

check_json_field() {
    local json=$1
    local field=$2
    local name=$3
    if echo "$json" | grep -q "\"$field\":"; then
        success "$name contains '$field'"
    else
        fail "$name missing '$field'"
    fi
}

wait_for_server() {
    log "Waiting for server at $BASE_URL (max ${TIMEOUT}s)..."
    local elapsed=0
    while ! curl -s -o /dev/null "$BASE_URL/status"; do
        sleep 1
        ((elapsed++))
        if [ $elapsed -ge $TIMEOUT ]; then
            fail "Server did not start in $TIMEOUT seconds"
            exit 1
        fi
    done
    log "Server is up!"
}

start_server() {
    if curl -s -o /dev/null "$BASE_URL/status"; then
        log "Server already running at $BASE_URL"
        return
    fi

    log "Starting Go server in background..."
    go run main.go > server.log 2>&1 &
    SERVER_PID=$!
    wait_for_server
}

cleanup() {
    if [ -n "$SERVER_PID" ] && kill -0 "$SERVER_PID" 2>/dev/null; then
        log "Stopping server (PID: $SERVER_PID)..."
        kill "$SERVER_PID" || true
        wait "$SERVER_PID" 2>/dev/null || true
    fi
    rm -f summary.png
}

trap cleanup EXIT

# ================================
# Run Tests
# ================================

echo -e "${YELLOW}Starting API Endpoint Tests...${NC}"
echo "Base URL: $BASE_URL"
echo "========================================"

start_server

# 1. POST /countries/refresh
log "1. Testing POST /countries/refresh"
REFRESH_RESP=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/countries/refresh" -o refresh.json)
REFRESH_CODE=${REFRESH_RESP: -3}
check_status "$REFRESH_CODE" "200" "POST /countries/refresh"
if [ "$REFRESH_CODE" -eq 200 ]; then
    check_json_field "$(cat refresh.json)" "message" "Refresh response"
fi

sleep 3  # Allow image generation

# 2. GET /status
log "2. Testing GET /status"
STATUS_RESP=$(curl -s -w "%{http_code}" "$BASE_URL/status" -o status.json)
STATUS_CODE=${STATUS_RESP: -3}
check_status "$STATUS_CODE" "200" "GET /status"
if [ "$STATUS_CODE" -eq 200 ]; then
    check_json_field "$(cat status.json)" "total_countries" "Status has total_countries"
    check_json_field "$(cat status.json)" "last_refreshed_at" "Status has last_refreshed_at"
fi

# 3. GET /countries?region=Africa
log "3. Testing GET /countries?region=Africa"
AFRICA_RESP=$(curl -s -w "%{http_code}" "$BASE_URL/countries?region=Africa" -o africa.json)
AFRICA_CODE=${AFRICA_RESP: -3}
check_status "$AFRICA_CODE" "200" "GET /countries?region=Africa"
if [ "$AFRICA_CODE" -eq 200 ]; then
    AFRICA_COUNT=$(jq '. | length' africa.json 2>/dev/null || echo "0")
    if [ "$AFRICA_COUNT" -gt 0 ]; then
        success "Africa filter returned $AFRICA_COUNT countries"
    else
        fail "Africa filter returned 0 countries"
    fi
fi

# 4. GET /countries?sort=gdp_desc
log "4. Testing GET /countries?sort=gdp_desc"
SORT_RESP=$(curl -s -w "%{http_code}" "$BASE_URL/countries?sort=gdp_desc" -o sort.json)
SORT_CODE=${SORT_RESP: -3}
check_status "$SORT_CODE" "200" "GET /countries?sort=gdp_desc"
if [ "$SORT_CODE" -eq 200 ]; then
    GDP1=$(jq '.[0].estimated_gdp // 0' sort.json 2>/dev/null || echo "0")
    GDP2=$(jq '.[1].estimated_gdp // 0' sort.json 2>/dev/null || echo "0")
    if (( $(echo "$GDP1 >= $GDP2" | bc -l 2>/dev/null || echo 0) )); then
        success "GDP descending order correct"
    else
        fail "GDP not in descending order"
    fi
fi

# 5. GET /countries/Nigeria
log "5. Testing GET /countries/Nigeria"
NIGERIA_RESP=$(curl -s -w "%{http_code}" "$BASE_URL/countries/Nigeria" -o nigeria.json)
NIGERIA_CODE=${NIGERIA_RESP: -3}
check_status "$NIGERIA_CODE" "200" "GET /countries/Nigeria"
if [ "$NIGERIA_CODE" -eq 200 ]; then
    check_json_field "$(cat nigeria.json)" "currency_code" "Nigeria has currency_code"
    check_json_field "$(cat nigeria.json)" "estimated_gdp" "Nigeria has estimated_gdp"
fi

# 6. GET /countries/image
log "6. Testing GET /countries/image"
curl -s -f -o summary.png "$BASE_URL/countries/image"
if [ $? -eq 0 ] && [ -f summary.png ] && [ -s summary.png ]; then
    success "Summary image downloaded ($(stat -f%z summary.png 2>/dev/null || stat -c%s summary.png) bytes)"
else
    fail "Failed to download summary image"
fi

# 7. DELETE /countries/TestCountry (should 404)
log "7. Testing DELETE /countries/TestCountry (expect 404)"
DELETE_RESP=$(curl -s -w "%{http_code}" -X DELETE "$BASE_URL/countries/TestCountry" -o /dev/null)
DELETE_CODE=${DELETE_RESP: -3}
check_status "$DELETE_CODE" "404" "DELETE non-existent country"

# 8. Try refresh again (should work)
log "8. Testing second refresh"
REFRESH2_RESP=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/countries/refresh" -o /dev/null)
REFRESH2_CODE=${REFRESH2_RESP: -3}
check_status "$REFRESH2_CODE" "200" "Second refresh"

# ================================
# Final Results
# ================================

echo "========================================"
echo -e "${YELLOW}Test Summary:${NC}"
echo -e "  PASSED: $PASSED"
if [ $FAILED -eq 0 ]; then
    echo -e "  ${GREEN}FAILED: $FAILED${NC}"
    echo -e "\n${GREEN}All tests passed! API is working perfectly.${NC}"
else
    echo -e "  ${RED}FAILED: $FAILED${NC}"
    echo -e "\n${RED}Some tests failed. Check output above.${NC}"
fi

exit $FAILED
