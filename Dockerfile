FROM cimg/go:1.24.1

WORKDIR /featherlb

COPY . .

ENV FEATHERLB_CONFIG_PATH=./config/featherlb.yaml

CMD [ "sh", "-c", "go run ./... --config $FEATHERLB_CONFIG_PATH" ]
