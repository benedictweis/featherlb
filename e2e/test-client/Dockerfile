FROM alpine:3.21

RUN apk add --no-cache wrk

ENV LOAD_BALANCER_HOST=featherlb
ENV LOAD_BALANCER_PORT=8080

CMD [ "sh", "-c", "sleep 10 && wrk -t12 -c400 -d10s http://$LOAD_BALANCER_HOST:$LOAD_BALANCER_PORT" ]
