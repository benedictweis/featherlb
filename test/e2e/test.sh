#!/usr/bin/env bash

command -v docker >/dev/null 2>&1 || {
    echo >&2 "Docker is required but not installed. Aborting."
    exit 1
}

cd "$(dirname "$0")" || exit

rm -rf ./config/featherlb.yaml

mkdir -p ./runs
logfile=./runs/$(date '+%Y%m%d%H%M%S').log
touch "$logfile"

for config in roundrobin random iphash; do
    echo "Testing $config..."
    echo "Config: $config" >>"$logfile"

    cp "./config/$config.yaml" ./config/featherlb.yaml

    docker compose up -d --build

    echo "Waiting for the test-client to to exit"
    docker compose wait test-client

    echo "Results for test $config:"

    logs=$(docker compose logs test-client)
    echo "$logs"
    echo "$logs" >>"$logfile"

    for backend in backend1 backend2; do
        echo "$backend requests:"
        request_count=$(docker compose exec "$backend" sh -c "cat /var/log/nginx/access.log | wc -l")
        echo "$request_count"
        echo "$backend: $request_count" >>"$logfile"
    done

    docker compose down
done

rm -rf ./config/featherlb.yaml
