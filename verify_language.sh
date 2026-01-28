#!/bin/bash

# Define the endpoint
API_URL="http://localhost:8080/api/plan/generate"

# Define the payload with the new language parameter
PAYLOAD='{
  "userId": "test-user-id",
  "gender": "male",
  "age": "25",
  "height": "180cm",
  "weight": "75kg",
  "mainGoal": "muscle gain",
  "workoutDuration": "60 mins",
  "experienceLevel": "intermediate",
  "language": "Chinese"
}'

# Send the request
echo "Sending request to $API_URL with payload:"
echo $PAYLOAD

curl -X POST "$API_URL" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_TEST_TOKEN" \
     -d "$PAYLOAD"

echo "\n\nRequest sent. Check server logs or response for output."
