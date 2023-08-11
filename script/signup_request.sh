#!/bin/bash


echo -n "Enter the number of iterations: "
read num_iterations

for ((i=1; i<=$num_iterations; i++))
do
    API_URL="http://localhost:8080/signup"



    month_size_limit=$((RANDOM % 100))
    minute_rate_limit=$((RANDOM % 10))

    PAYLOAD='{
    "month_size_limit": '$month_size_limit',
    "minute_rate_limit": '$minute_rate_limit'
    }'


    curl -X POST "$API_URL" \
    -H "Content-Type: application/json" \
    -d "$PAYLOAD"

done

sleep 20


