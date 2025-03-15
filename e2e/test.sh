#!/usr/bin/env bash

cd "$(dirname "$0")" || exit

docker compose up -d --build

echo "Waiting for the test-client to to exit"
docker compose wait test-client

docker compose logs test-client

for backend in backend1 backend2; do
    echo "$backend requests:"
    docker compose exec "$backend" sh -c "cat /var/log/nginx/access.log | wc -l"
done

docker compose down
