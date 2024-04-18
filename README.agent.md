<h1 align='center'>pst-agent 部署</h1>

<p align="center">
   <strong>简体中文</strong> | <a href="/README.agent.en.md">English</a> | <a href="/README.agent.ja.md">日本語</a>
</p>

### Linux

这里指的是，你的游戏服务器是 Linux 版本的，而你想在其它任何地方部署的 pst 本体依然参考前面的 [安装部署](./README.md#安装部署)，pst-agent 对其的影响只是配置文件的修改。

#### 下载

下载 pst-agent 工具并重命名、确保其可执行

```bash
# 下载并重命名
mv pst-agent_v0.7.1_linux_x86_64 pst-agent
chmod +x pst-agent
```

#### 运行

```bash
# ./pst-agent --port 8081 -d {Level.sav 存档所在绝对路径}
# 例如：
./pst-agent --port 8081 -d /home/lighthouse/game/Saved/
```

检查正常运行后，让其后台运行（关闭 ssh 窗口后仍运行）

```bash
# 后台运行并将日志保存在 agent.log
nohup ./pst-agent --port 8081 -d ...{手动省略}.../Saved > agent.log 2>&1 &
# 查看日志
tail -f agent.log
```

#### 开放防火墙/安全组

如果 pst-agent 和 pst 本体完全没在同一组网内，需要放开游戏服务器的相应公网端口（如 8081，也可以是自定义的其它端口）

#### 配置

找到 **pst 本体（注意，不是 pst-agent！！！）** 的 `config.yaml` 文件并修改。

**只需要修改其中的 `save: path`**

```yaml
...
save: # 存档文件解析相关配置
  path: "http://游戏服务器公网IP:端口/sync" # pst-agent 服务接口
  ...
...
```

重启 pst 本体即可

#### 关闭后台运行

```bash
kill $(ps aux | grep 'pst-agent' | awk '{print $2}') | head -n 1
```

### Windows

这里指的是，你的游戏服务器是 Windows 版本的，而你想在其它任何地方部署的 pst 本体依然参考前面的 [安装部署](./README.md#安装部署)，pst-agent 对其的影响只是配置文件的修改。

#### 下载

下载 pst-agent 工具并重命名，如将 `pst-agent_v0.7.1_windows_x86_64.exe` 重命名为 `pst-agent.exe`

#### 运行

按下 `Win + R`，输入 `powershell` 打开 Powershell，通过 `cd` 命令到下载的可执行文件目录

```powershell
# .\pst-agent.exe --port 访问端口 -d 存档文件 Level.sav 所在位置
.\pst-agent.exe --port 8081 -d C:\Users\ZaiGie\...\Pal\Saved
```

成功运行后请保持窗口打开

#### 配置

找到 **pst 本体（注意，不是 pst-agent！！！）** 的 `config.yaml` 文件并修改。

**只需要修改其中的 `save: path`**

```yaml
...
save: # 存档文件解析相关配置
  path: "http://游戏服务器公网IP:端口/sync" # pst-agent 服务接口
  ...
...
```

重启 pst 本体即可
