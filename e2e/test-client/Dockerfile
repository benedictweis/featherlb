FROM alpine/curl

ENV LOAD_BALANCER_HOST=featherlb
ENV LOAD_BALANCER_PORT=8080

CMD [ "sh", "-c", "curl --retry 10 --retry-delay 1 $LOAD_BALANCER_HOST:$LOAD_BALANCER_PORT" ]
