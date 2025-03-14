FROM alpine/curl

ENV LOAD_BALANCER_HOST=featherlb
ENV LOAD_BALANCER_PORT=8080

CMD [ "sh", "-c", "sleep 10 && for i in $(seq 1 10); do curl -s $LOAD_BALANCER_HOST:$LOAD_BALANCER_PORT; echo ''; sleep 0.5; done && sleep 100000" ]
