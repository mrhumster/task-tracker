.DEFAULT_GOAL := build

BIN_NAME := task-cli 


fmt:
	go fmt ./...
.PHONY: fmt

lint: fmt
	golint ./...
.PHONY: lint

vet: fmt
	go vet ./...
.PHONY: vet

build: vet
	go build -o ${BIN_NAME} .
.PHONY: build

install: build
	sudo cp ${BIN_NAME} /usr/local/bin/
.PHONY: install

clean:
		rm -f ${BIN_NAME}
.PHONY: clean
