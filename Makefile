GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
GIT_TAG:=$(shell git describe --tags --abbrev=0)
PREFIX:=pst_${GIT_TAG}
OS=$(uname)
EXT=""
ifeq ($(OS),Windows_NT)
    EXT := .exe
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        EXT :=
    endif
    ifeq ($(UNAME_S),Darwin)
        EXT :=
    endif
endif

.PHONY: init
# 初始化
init:
	go mod download

.PHONY: build
# 构建
build:
	rm -rf dist/ && mkdir -p dist/

	rm -rf assets && rm -rf index.html && rm -rf pal-conf.html
	cd web && pnpm i && pnpm build && cd ..
	git submodule update --init --recursive
	cd pal-conf && pnpm i && pnpm build && cd ..
	mv pal-conf/dist/assets/* assets/
	mv pal-conf/dist/index.html ./pal-conf.html

	cd module && pip install -r requirements.txt
	cd module && pyinstaller --onefile sav_cli.py
	mv module/dist/sav_cli${EXT} ./dist/

	cp example/config.yaml dist/config.yaml
	go build -o ./dist/pst${EXT} main.go

.PHONY: build-pub
# 为所有平台构建，确保 module/dist 中有所有平台的 sav_cli
build-pub:
	rm -rf dist/ && mkdir -p dist/

	rm -rf assets && rm -rf index.html && rm -rf pal-conf.html
	cd web && pnpm i && pnpm build && cd ..
	git submodule update --init --recursive
	cd pal-conf && pnpm i && pnpm build && cd ..
	mv pal-conf/dist/assets/* assets/
	mv pal-conf/dist/index.html ./pal-conf.html

	mkdir -p dist/windows_x86_64 && mkdir -p dist/linux_x86_64 && mkdir -p dist/linux_aarch64 && mkdir -p dist/darwin_arm64
	GOOS=windows GOARCH=386 go build -ldflags="-s -w -X 'main.version=${GIT_TAG}'" -o ./dist/windows_x86_64/pst.exe main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X 'main.version=${GIT_TAG}'" -o ./dist/linux_x86_64/pst main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X 'main.version=${GIT_TAG}'" -o ./dist/linux_aarch64/pst main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X 'main.version=${GIT_TAG}'" -o ./dist/darwin_arm64/pst main.go

	GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o ./dist/pst-agent_${GIT_TAG}_windows_x86_64.exe ./cmd/pst-agent/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./dist/pst-agent_${GIT_TAG}_linux_x86_64 ./cmd/pst-agent/main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ./dist/pst-agent_${GIT_TAG}_linux_aarch64 ./cmd/pst-agent/main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ./dist/pst-agent_${GIT_TAG}_darwin_arm64 ./cmd/pst-agent/main.go

	cp module/dist/sav_cli_windows_x86_64.exe dist/windows_x86_64/sav_cli.exe
	cp module/dist/sav_cli_linux_x86_64 dist/linux_x86_64/sav_cli
	cp module/dist/sav_cli_linux_aarch64 dist/linux_aarch64/sav_cli
	cp module/dist/sav_cli_darwin_arm64 dist/darwin_arm64/sav_cli

	cp example/config.yaml dist/windows_x86_64/config.yaml
	cp example/config.yaml dist/linux_x86_64/config.yaml
	cp example/config.yaml dist/linux_aarch64/config.yaml
	cp example/config.yaml dist/darwin_arm64/config.yaml

	cp script/start.bat dist/windows_x86_64/start.bat

	cd dist && zip -p -r ${PREFIX}_windows_x86_64.zip windows_x86_64/* && tar -czf ${PREFIX}_linux_x86_64.tar.gz linux_x86_64/* && tar -czf ${PREFIX}_linux_aarch64.tar.gz linux_aarch64/* && tar -czf ${PREFIX}_darwin_arm64.tar.gz darwin_arm64/* && cd ..
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
