.PHONY: run-jet-stream
run-jet-stream:
	docker run -p 4222:4222 nats -js

.PHONY: run
run:
	go run cmd/api/*.go