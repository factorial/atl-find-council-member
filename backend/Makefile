PID=.hot-compile.pid
PWD=$(shell pwd)
BIN_DIR=$(PWD)/bin
SRC_DIR=$(PWD)/src
CLI_BIN_NAME=cli
SERVER_BIN_NAME=server
DATA_DIR?=$(PWD)/data

all: build

build: build-server build-cli

initialize-dependencies:
	-@go mod init
	-@go mod vendor

update-dependencies:
	cd $(SRC_DIR) && go get -u ./...

build-server:
	cd $(SRC_DIR) && go build -v -o $(BIN_DIR)/$(SERVER_BIN_NAME) backend/cmd/server

build-cli:
	cd $(SRC_DIR) && go build -v -o $(BIN_DIR)/$(CLI_BIN_NAME) backend/cmd/cli


container:
	docker image build -t goeieware.ca/mapatlapi:latest .
	docker push goeieware.ca/mapatlapi:latest

test:
	cd $(SRC_DIR) && go test -coverprofile=coverage.out ./...

test-coverage:
	-cd $(SRC_DIR) && @go tool cover -html=coverage.out

watch-test: test
	@fswatch -l 0.5 -o $(SRC_DIR) -e ".*" -i "\\.go$$" | xargs -n1 -I{} make test

start: kill build
	@$(BIN_DIR)/$(BIN_NAME) -d $(DATA_DIR) & echo $$! > $(PID)

restart: notify start

notify:
	@echo Recompiling and reloading...

kill:
	-@kill `cat $(PID) 2>/dev/null` 2>/dev/null || true
	-@rm $(PID) || true

watch: restart
	@fswatch -l 0.5 -o $(SRC_DIR) -e ".*" -i "\\.go$$" | xargs -n1 -I{}  make restart

.PHONY: clean
clean:
	-rm $(BIN_DIR)/$(BIN_NAME)
