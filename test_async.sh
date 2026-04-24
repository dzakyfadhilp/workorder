#!/bin/bash

# Test script untuk async processing dengan RabbitMQ

echo "=== Testing Async Processing with RabbitMQ ==="
echo ""

echo "Sending request..."
START=$(date +%s%3N)

RESPONSE=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8080/api/execute \
  -H "Content-Type: application/json" \
  -d '{
    "function": "ff_updateWorkorder",
    "payload": {
      "req": {
        "wonum": "WO'$(date +%s)'",
        "siteid": "REG-6",
        "status": "STARTWA",
        "task": "updateTaskStatus",
        "labor_scmt": "23000036"
      },
      "res": {
        "data": "WO001",
        "status": true,
        "message": "Success"
      },
      "date": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
    }
  }')

END=$(date +%s%3N)
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

echo ""
echo "Response Code: $HTTP_CODE"
echo "Response Body:"
echo "$BODY" | jq .
echo ""
echo "Response Time: $((END - START))ms"
echo ""

if [ "$HTTP_CODE" = "202" ]; then
    echo "✅ Request accepted (202 Accepted)"
    echo "✅ Async processing enabled"
    echo ""
    echo "Check server logs to see background processing..."
else
    echo "❌ Unexpected response code: $HTTP_CODE"
fi
