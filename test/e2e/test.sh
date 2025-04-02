#!/usr/bin/env bash

cd "$(dirname "$0")" || exit

rm -rf ./config/featherlb.yaml

for config in roundrobin random iphash; do
    echo "Testing $config..."

    cp "./config/$config.yaml" ./config/featherlb.yaml

    docker compose up -d --build

    echo "Waiting for the test-client to to exit"
    docker compose wait test-client

    echo "Results for test $config:"

    docker compose logs test-client

    for backend in backend1 backend2; do
        echo "$backend requests:"
        docker compose exec "$backend" sh -c "cat /var/log/nginx/access.log | wc -l"
    done

    docker compose down
done

rm -rf ./config/featherlb.yaml
