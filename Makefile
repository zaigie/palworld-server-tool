GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
GIT_TAG:=$(shell git describe --tags)

.PHONY: init
# 初始化
init:
	go mod download

.PHONY: build
# 构建
build:
	rm -rf dist/ && mkdir -p dist/
	go build -o ./dist/manager main.go

.PHONY: build-all
# 为所有平台构建
build-all:
	rm -rf dist/ && mkdir -p dist/
	GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o ./dist/pst-${GIT_TAG}_windows_x86.exe main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./dist/pst-${GIT_TAG}_linux_amd64 main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ./dist/pst-${GIT_TAG}_linux_arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ./dist/pst-${GIT_TAG}_darwin_amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ./dist/pst-${GIT_TAG}_darwin_arm64 main.go
# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
