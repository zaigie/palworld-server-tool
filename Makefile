GIT_TAG:=$(shell git describe --tags --abbrev=0)
IMAGE?=palworld-server-tool

.PHONY: init
# 初始化子模块
init:
	git submodule update --init --recursive

.PHONY: build
# 使用当前 Dockerfile 构建本机架构镜像
build: init
	docker build --build-arg version=$(GIT_TAG) -t $(IMAGE):$(GIT_TAG) .

.PHONY: build-pub
# 使用当前 Dockerfile 构建 ARM64/AMD64 OCI 镜像包
build-pub: init
	rm -rf dist/ && mkdir -p dist/
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg version=$(GIT_TAG) \
		--output type=oci,dest=dist/$(IMAGE)_$(GIT_TAG).tar \
		.

.PHONY: frontend
# 仅构建前端开发产物
frontend: init
	rm -rf assets index.html pal-conf.html
	cd web && pnpm install --frozen-lockfile && pnpm build
	cd pal-conf && pnpm install --frozen-lockfile && pnpm build
	mv pal-conf/dist/assets/* assets/
	mv pal-conf/dist/index.html ./pal-conf.html

.PHONY: test-sav
# 构建双架构镜像并验证真实存档
test-sav:
	python3 script/test_sav_cli.py --no-cache

# 显示帮助
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
