# palworld-server-tool

通过 可视化界面及 REST 接口或命令行工具管理幻兽帕鲁 / PalWorld 专用服务器，基于 RCON 实现

![PC](./doc/img/pc.png)

基于官方提供的 RCON 命令（仅服务器可用的）实现功能如下：

- [x] 获取服务器信息
- [x] 玩家列表（历史玩家数据仅 pst-server）
- [x] 踢出/封禁玩家
- [x] 游戏内广播
- [x] 平滑关闭服务器并广播消息

请在以下地址下载最新版可执行文件

- [Github Releases](https://github.com/zaigie/palworld-server-tool/releases)
- [(国内) Gitee Releases](https://gitee.com/jokerwho/palworld-server-tool/releases)

## 如何开启私服 RCON

需要开启服务器的 RCON 功能，如果你的私服教程有写更好，没有的话，修改 `PalWorldSettings.ini` 文件

**也就是修改游戏内各种倍数、概率的那个文件**，里面最后的位置有如下：

```txt
RCONEnabled=true,RCONPort=25575
```

![RCON](./doc/img/rcon.png)

请先关闭服务器，然后将 `RCONEnabled` 和 `RCONPort` 填写如上，再重启服务器即可

## pst-server: 可视化工具

![Mobile](./doc/img/mobile.png)

服务使用 bbolt 单文件数据库，用来存历史玩家数据，并且每五分钟会定时查询一次在线玩家列表，更新最后在线时间。

> [!CAUTION]
> 由于数据库文件系统变更，**v0.3.0 以下版本**请删除原来的 `players.db`

### Linux

1. 下载文件并重命名

```bash
# 下载 pst-server_{version}_{platform}_{arch} 文件并重命名
mv pst-server_v0.3.2_linux_amd64 pst-server
```

2. 运行

如果和服务器部署在同一机器上

```bash
./pst-server -p {你的 AdminPassword}
```

如果和服务器不在同一机器上

```bash
./pst-server -a {服务器IP:RCON端口} -p {你的 AdminPassword}
```

3. 后台运行

```bash
# 后台运行并将日志保存在 server.log
nohup ./pst-server -a 127.0.0.1:25575 -p {你的 AdminPassword} > server.log 2>&1 &
# 查看日志
tail -f server.log
```

5. 关闭后台程序

```bash
# 关闭程序
kill $(ps aux | grep 'pst-server' | awk '{print $2}') | head -n 1
```

请通过浏览器访问 http://127.0.0.1:8080 或 http://{局域网 IP}:8080

云服务器也可以访问 http://{服务器 IP}:8080

> [!WARNING]
> 如果你想变更工具服务运行的端口（默认 8080），则可以在命令上加上 --port 8000 来自定义

### Windows

1. 下载文件并重命名

将 pst-server_v0.3.2_windows_x86.exe 重命名为 pst-server.exe

2. 按下 `Win + R`，输入 `powershell` 打开 Powershell，通过 `cd` 命令到下载的可执行文件目录

3. 持续运行，不要关闭窗口

如果和服务器部署在同一机器上

```powershell
.\pst-server.exe -p {你的 AdminPassword}
```

如果和服务器不在同一机器上

```powershell
.\pst-server.exe -a {服务器IP:RCON端口} -p {你的 AdminPassword}
```

请通过浏览器访问 http://127.0.0.1:8080 或 http://{局域网 IP}:8080
云服务器也可以访问 http://{服务器 IP}:8080

> [!WARNING]
> 如果你想变更工具服务运行的端口（默认 8080），则可以在命令上加上 --port 8000 来自定义

若要自己开发前端界面或用作它用请移步 [接口文档](./API.md)

---

## pst-cli: 命令行工具

```bash
# 下载 pst-cli_{version}_{platform}_{arch} 文件并重命名
mv pst-cli_{version}_{platform}_{arch} pst-cli
```

首次运行会自动生成 `config.yaml` ，修改文件

```yaml
host: 127.0.0.1:25575
password: "你的 AdminPassword"
timeout: 10
```

### 玩家

#### 在线玩家列表

```bash
./pst-cli player list
```

```
+-------------------------------------------+
| Pal World 在线玩家列表                    |
+----------+------------+-------------------+
| 昵称     | PLAYERUID  | STEAMID           |
+----------+------------+-------------------+
| 香菇包子 | 2398722357 | xxxxx |
| 梵音丶   | 2144044083 | xxxxx |
| 狐狸     | 1333009711 | xxxxx |
| Baoz     | <null/err> | <null/err>        |
+----------+------------+-------------------+
|          | 在线人数   | 4                 |
+----------+------------+-------------------+
```

> <null/err> 是帕鲁服务器的错误，待官方修复

#### 踢出/封禁玩家

```bash
./pst-cli kick -s <SteamID>
./pst-cli ban -s <SteamID>
```

### 广播

```bash
./pst-cli broadcast -m "<message>"
```

> [!WARNING]
> message 中不能包含中文

### 服务器

#### 关闭服务器

```bash
./pst-cli server shutdown -s <seconds> -m "Server Will Shutdown"
```
