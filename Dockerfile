FROM cimg/go:1.23

WORKDIR /featherlb

COPY . .

CMD [ "go", "run", "./..." ]
