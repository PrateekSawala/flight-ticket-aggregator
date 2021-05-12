# Build, test, package and more ... A developer everyday tool
BUILD_DIR := bin
TOOLS_DIR := tools

default: all

all: clean lint test build image run

.PHONY: $(BUILD_DIR)/server
bin/server: cmd/*.go
	CGO_ENABLED=0 go build -o ./bin/server ./cmd/

.PHONY: build
build: clean bin/server

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

.PHONY: run
run: build
	bin/server

tools/golangci-lint/golangci-lint:
	mkdir -p tools/
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./tools/golangci-lint latest

.PHONY: lint
lint: $(TOOLS_DIR)/golangci-lint/golangci-lint
	./$(TOOLS_DIR)/golangci-lint/golangci-lint run ./...
	
.PHONY: test
test:
	-go test -v -cover -coverprofile=cover.out -covermode=atomic ./service/...
	go tool cover -html=cover.out -o cover.html	
 
.PHONY: image
image: build
	docker build -t test:test -f Dockerfile ./bin/.