# AGENTS.md

## 沟通约定

- 始终使用简体中文回复。

## Release 发布流程

除非维护者明确要求采用其他方式，否则发布正式版本时遵循以下流程。

### 1. 发布前检查

1. 所有待发布改动先提交到 `dev`，不要直接在 `main` 上开发。
2. 确认工作树干净，并确保 Git LFS 对象已经上传：

   ```bash
   git switch dev
   git pull --ff-only origin dev
   git status
   git lfs status
   git push origin dev
   ```

3. 验证地图瓦片完整，不能包含未展开的 LFS 指针：

   ```bash
   python3 script/verify_map.py
   ```

4. 创建 `dev -> main` 的正式 PR。PR 描述应覆盖相较于 `main` 的全部变更、影响范围和验证结果。
5. Review 没有阻塞问题后合并 PR。合并到 `main` 不应触发发布构建；正式构建只由版本 tag 触发。

### 2. 创建正式版本 tag

只有已经进入 `main` 的提交才能打正式 tag。tag 格式必须为 `v<major>.<minor>.<patch>`，例如 `v0.12.0`。

仓库历史版本使用轻量 tag，继续采用相同格式：

```bash
git fetch origin --tags
git switch main
git pull --ff-only origin main

VERSION=v0.12.0
git tag "$VERSION"
git push origin "$VERSION"
```

不要在 PR 合并前打 tag。工作流会验证 tag 指向的 commit 是否属于 `origin/main`；不满足时发布会失败。

### 3. 自动执行的工作流

推送正式 tag 后，GitHub Actions 会自动运行：

- `Native Build and Release`
  - 构建 Linux amd64、Linux arm64、Windows amd64 和 macOS arm64 原生文件。
  - 自动创建或更新对应的 GitHub Release。
  - 自动生成 `SHA256SUMS.txt`。
- `PST Docker Image CI`
  - 发布 `jokerwho/palworld-server-tool:<version>` 和 `latest`。
  - Docker 版本标签会去掉开头的 `v`，例如 `v0.12.0` 对应 `0.12.0`。
- `PST-Agent Docker Image CI`
  - 发布 `jokerwho/palworld-server-tool-agent:<version>` 和 `latest`。

仅合并到 `main` 不会运行以上三条发布工作流。

`Native Build and Release` 会先将各平台文件保存为临时 Actions artifacts；最后的 `release` job 会自动下载这些文件、生成校验和，并将全部文件上传为 GitHub Release assets。不要在 tag 推送后手工执行 `gh release create`，也不需要在网页中选择构建文件，以免产生草稿、重复 Release 或资产缺失。

### 4. Release 资产检查

正式 Release 应包含以下 9 个资产：

- `pst-agent_<version>_darwin_aarch64`
- `pst-agent_<version>_linux_aarch64`
- `pst-agent_<version>_linux_x86_64`
- `pst-agent_<version>_windows_x86_64.exe`
- `pst_<version>_darwin_aarch64.tar.gz`
- `pst_<version>_linux_aarch64.tar.gz`
- `pst_<version>_linux_x86_64.tar.gz`
- `pst_<version>_windows_x86_64.zip`
- `SHA256SUMS.txt`

检查 Release：

```bash
gh release view "$VERSION" --repo zaigie/palworld-server-tool
gh release view "$VERSION" \
  --repo zaigie/palworld-server-tool \
  --json assets \
  --jq '.assets[].name'
```

### 5. Docker 多架构检查

两个正式镜像都必须包含 `linux/amd64` 和 `linux/arm64`：

```bash
DOCKER_VERSION="${VERSION#v}"

docker buildx imagetools inspect \
  "jokerwho/palworld-server-tool:${DOCKER_VERSION}"

docker buildx imagetools inspect \
  "jokerwho/palworld-server-tool-agent:${DOCKER_VERSION}"
```

同时检查对应版本标签与 `latest` 是否指向同一个顶层 manifest digest。

### 6. Release 说明

工作流默认生成 GitHub 自动发布说明。需要保持中英双语说明时，等待 Release 自动创建完成，再准备 `RELEASE_NOTES.md` 并更新：

```bash
gh release edit "$VERSION" \
  --repo zaigie/palworld-server-tool \
  --title "$VERSION" \
  --notes-file RELEASE_NOTES.md
```

说明至少应包含中文变更、English changes 和相较于上一版本的 Full Changelog 链接。

### 7. 失败恢复

- 网络波动或单个 job 失败时，在 GitHub Actions 页面选择 `Re-run failed jobs`。
- 不要为了重试而重复创建 Release，也不要随意删除或移动已经公开的 tag。
- 原生 Release 成功、Docker 失败时，只重跑失败的 Docker workflow。
- Docker 成功、原生构建失败时，只重跑 `Native Build and Release` 的失败任务。
- 本机磁盘空间有限，不要默认在本地执行完整双架构镜像构建；优先使用 GitHub Actions。
- Docker Hub 登录依赖仓库 secrets：`DOCKER_HUB_USERNAME` 和 `DOCKER_HUB_PWD`。

## 仅测试 PST Docker 镜像

需要验证 Dockerfile、但不能创建 Release 或覆盖正式镜像标签时，手动运行：

```bash
gh workflow run docker-image.yml \
  --repo zaigie/palworld-server-tool \
  --ref dev
```

手动运行只发布 `manual-<short-sha>` 标签，不更新 `edge`、语义版本标签或 `latest`，也不会触发原生 Release 和 PST-Agent 构建。
