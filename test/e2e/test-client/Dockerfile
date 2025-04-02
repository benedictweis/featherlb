FROM alpine:3.21

RUN apk add --no-cache wrk curl

ENV LOAD_BALANCER_HOST=featherlb
ENV LOAD_BALANCER_PORT=8080

ENV TARGET_URL=http://$LOAD_BALANCER_HOST:$LOAD_BALANCER_PORT

CMD [ "sh", "-c", "\
    for i in $(seq 1 20); do \
        if curl -s $TARGET_URL > /dev/null; then \
            wrk -t12 -c400 -d10s $TARGET_URL; \
            exit 0; \
        fi; \
        sleep 1; \
    done; \
    echo 'Target not reachable'; \
    exit 1" ]
