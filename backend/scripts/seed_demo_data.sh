#!/bin/bash
set -e

BASE_URL="http://localhost:8080"
EMAIL="demo@flagr.io"
PASSWORD="DemoPassword123456"
ORG_NAME="Demo Organization"

echo "Flagr Demo Seed"
echo "=================="

echo ""
echo "1. Registering user..."
REGISTER=$(curl -s -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\",\"org_name\":\"$ORG_NAME\"}")

TOKEN=$(echo "$REGISTER" | jq -r '.access_token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo "   User exists, logging in..."
  LOGIN=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")
  TOKEN=$(echo "$LOGIN" | jq -r '.access_token')
fi

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo "Auth failed"
  exit 1
fi
echo "   Authenticated"

PROJECT_ID=$(curl -s "$BASE_URL/api/v1/projects" \
  -H "Authorization: Bearer $TOKEN" | jq -r '.[0].id')

ENV_ID=$(curl -s "$BASE_URL/api/v1/environments?project_id=$PROJECT_ID" \
  -H "Authorization: Bearer $TOKEN" | jq -r '.[0].id')

echo "   Project: $PROJECT_ID"
echo "   Environment: $ENV_ID"

echo ""
echo "2. Creating feature flags..."
for FLAG in "checkout-v2:Checkout V2:New checkout flow" \
            "dark-mode:Dark Mode:Enable dark theme" \
            "new-dashboard:New Dashboard:Redesigned dashboard"; do
  KEY=$(echo $FLAG | cut -d: -f1)
  NAME=$(echo $FLAG | cut -d: -f2)
  DESC=$(echo $FLAG | cut -d: -f3)

  RESULT=$(curl -s -X POST "$BASE_URL/api/v1/flags" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"project_id\":\"$PROJECT_ID\",\"key\":\"$KEY\",\"name\":\"$NAME\",\"description\":\"$DESC\",\"type\":\"boolean\"}")

  echo "   Created: $NAME"
done

docker exec flagr_postgres psql -U flagr -d flagr -c "
UPDATE flag_environments fe
SET enabled = true
FROM flags f
WHERE fe.flag_id = f.id
  AND f.key = 'checkout-v2'
  AND fe.environment_id = '$ENV_ID';
" > /dev/null 2>&1
echo "   checkout-v2 enabled"

echo ""
echo "3. Creating SDK key..."
SDK_RESULT=$(curl -s -X POST "$BASE_URL/api/v1/sdk-keys" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"project_id\":\"$PROJECT_ID\",\"environment_id\":\"$ENV_ID\",\"name\":\"Demo SDK Key\"}")

SDK_KEY=$(echo "$SDK_RESULT" | jq -r '.key')
echo "   SDK Key: $SDK_KEY"

echo ""
echo "=================="
echo "Seed complete!"
echo ""
echo "Credentials:"
echo "  Email:    $EMAIL"
echo "  Password: $PASSWORD"
echo ""
echo "SDK Key: $SDK_KEY"
echo ""
echo "Test evaluate:"
echo "  curl -s 'http://localhost:8080/api/v1/flags/evaluate?key=checkout-v2&project_id=$PROJECT_ID&environment_id=$ENV_ID' \\"
echo "    -H 'X-SDK-Key: $SDK_KEY' | jq ."