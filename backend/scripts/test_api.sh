#!/bin/bash
set -e

BASE_URL="http://localhost:8080"
PASS=0
FAIL=0

check() {
  local desc=$1
  local expected=$2
  local actual=$3

  if [ "$actual" -eq "$expected" ] 2>/dev/null; then
    echo "  $desc"
    PASS=$((PASS + 1))
  else
    echo "  $desc (expected $expected, got $actual)"
    FAIL=$((FAIL + 1))
  fi
}

echo "Flagr API Smoke Tests"
echo "========================"

echo ""
echo "Health checks:"
check "GET /health/live" 200 $(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health/live")
check "GET /health/ready" 200 $(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health/ready")

echo ""
echo "Auth:"
UNIQUE=$(date +%s)
REG=$(curl -s -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"test_$UNIQUE@test.com\",\"password\":\"TestPassword123\",\"org_name\":\"Test Org $UNIQUE\"}")

TOKEN=$(echo "$REG" | jq -r '.access_token')
check "POST /auth/register" 1 $([ "$TOKEN" != "null" ] && echo 1 || echo 0)

LOGIN=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"test_$UNIQUE@test.com\",\"password\":\"TestPassword123\"}")
TOKEN=$(echo "$LOGIN" | jq -r '.access_token')
check "POST /auth/login" 1 $([ "$TOKEN" != "null" ] && echo 1 || echo 0)

echo ""
echo "Projects:"
check "GET /projects" 200 $(curl -s -o /dev/null -w "%{http_code}" \
  -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/projects")

PROJECT_ID=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/projects" | jq -r '.[0].id')
ENV_ID=$(curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/environments?project_id=$PROJECT_ID" | jq -r '.[0].id')

echo ""
echo "Flags:"
check "GET /flags" 200 $(curl -s -o /dev/null -w "%{http_code}" \
  -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/v1/flags?project_id=$PROJECT_ID")

check "POST /flags" 201 $(curl -s -o /dev/null -w "%{http_code}" \
  -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d "{\"project_id\":\"$PROJECT_ID\",\"key\":\"smoke-test-$UNIQUE\",\"name\":\"Smoke Test\",\"type\":\"boolean\"}" \
  "$BASE_URL/api/v1/flags")

echo ""
echo "========================"
echo "Results: $PASS passed, $FAIL failed"

[ $FAIL -eq 0 ] && exit 0 || exit 1