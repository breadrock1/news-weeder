BIN_FILE_PATH := "./bin/news-weeder"

build:
	go build -v -o $(BIN_FILE_PATH) ./cmd/news-weeder

run:
	$(BIN_FILE_PATH) -c ./configs/production.toml

test:
	go test -race ./...

.PHONY: build run test
