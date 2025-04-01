
.PHONY: e2e
e2e:
	./test/e2e/test.sh

.PHONY: run
run:
	go run ./cmd/featherlb --config featherlb.yaml