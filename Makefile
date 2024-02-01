GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
GIT_TAG:=$(shell git describe --tags)
PREFIX:=pst_${GIT_TAG}

.PHONY: init
# 初始化
init:
	go mod download

.PHONY: build
# 构建
build:
	rm -rf dist/ && mkdir -p dist/
	go build -o ./dist/pst main.go

.PHONY: build-all
# 为所有平台构建，确保 module/dist 中有所有平台的 sav_cli
build-all:
	rm -rf dist/ && mkdir -p dist/

	cd web && pnpm i && pnpm build && cd ..

	mkdir -p dist/windows_x86 && mkdir -p dist/linux_amd64 && mkdir -p dist/linux_arm64 && mkdir -p dist/darwin_arm64
	GOOS=windows GOARCH=386 go build -ldflags="-s -w -X 'main.version=${GIT_TAG}'" -o ./dist/windows_x86/pst.exe main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X 'main.version=${GIT_TAG}'" -o ./dist/linux_amd64/pst main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X 'main.version=${GIT_TAG}'" -o ./dist/linux_arm64/pst main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X 'main.version=${GIT_TAG}'" -o ./dist/darwin_arm64/pst main.go

	GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o ./dist/pst-agent_${GIT_TAG}_windows_x86.exe ./cmd/pst-agent/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./dist/pst-agent_${GIT_TAG}_linux_amd64 ./cmd/pst-agent/main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ./dist/pst-agent_${GIT_TAG}_linux_arm64 ./cmd/pst-agent/main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ./dist/pst-agent_${GIT_TAG}_darwin_arm64 ./cmd/pst-agent/main.go

	cp module/dist/sav_cli_windows_x86.exe dist/windows_x86/sav_cli.exe
	cp module/dist/sav_cli_linux_amd64 dist/linux_amd64/sav_cli
	cp module/dist/sav_cli_linux_arm64 dist/linux_arm64/sav_cli
	cp module/dist/sav_cli_darwin_arm64 dist/darwin_arm64/sav_cli

	cp test/config.yaml dist/windows_x86/config.yaml
	cp test/config.yaml dist/linux_amd64/config.yaml
	cp test/config.yaml dist/linux_arm64/config.yaml
	cp test/config.yaml dist/darwin_arm64/config.yaml

	cp script/start.bat dist/windows_x86/start.bat

	cd dist && zip -p -r ${PREFIX}_windows_x86.zip windows_x86/* && tar -czf ${PREFIX}_linux_amd64.tar.gz linux_amd64/* && tar -czf ${PREFIX}_linux_arm64.tar.gz linux_arm64/* && tar -czf ${PREFIX}_darwin_arm64.tar.gz darwin_arm64/* && cd ..
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
