FROM cimg/go:1.24.1

WORKDIR /featherlb

COPY . .

RUN go mod download

ENV FEATHERLB_CONFIG_PATH=./config/featherlb.yaml

CMD [ "sh", "-c", "go run ./... --config $FEATHERLB_CONFIG_PATH" ]
