# 自助编译及更新指南（linux）

## 当前实现说明

当前 `sav_cli` 使用 `palsav-flex` 和 `palooz` 解析 Palworld 1.0 存档。
它通过 Docker 镜像交付，避免宿主机 Python、编译器和原生依赖版本不同造成问题。
依赖版本和基础镜像均已固定，更新后必须重新运行真实存档验证。

## 编译方法

## 1.下载源码文件

`git clone https://github.com/zaigie/palworld-server-tool.git`

## 2.构建当前应用镜像

仓库只保留根目录的一个 `Dockerfile`。它会构建前端、`pal-conf`、后端、地图
资源和 Palworld 1.0 存档解析器：

`docker build -t palworld-server-tool:latest .`

该镜像仍使用原来的 `/app/sav_cli` 入口，因此已有
`SAVE__DECODE_PATH=/app/sav_cli` 配置不需要迁移。在 Windows Docker Desktop
中运行 Linux 容器时也无需修改这个容器内路径。

## 3.验证当前版本

仓库维护者可使用 `savs/` 下的两份测试存档，从头构建并验证 ARM64、AMD64：

`python3 script/test_sav_cli.py --no-cache`

这个验证会使用临时容器检查镜像架构、依赖版本和真实存档解析结果。
