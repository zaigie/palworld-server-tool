<h1 align='center'>pst-agent deployment</h1>

<p align="center">
   <a href="/docs/README.agent.md">简体中文</a> | <strong>English</strong> | <a href="/docs/README.agent.ja.md">日本語</a>
</p>

### Linux

This section covers a Linux game server with the PST main application deployed elsewhere. Follow the regular [installation guide](./README.en.md#installation), then select pst-agent as the save source in the Web settings.

#### Download

Download the pst-agent tool, rename it, and ensure it's executable

```bash
# Download and rename
mv pst-agent_v0.10.0_linux_x86_64 pst-agent
chmod +x pst-agent
```

#### Run

```bash
# ./pst-agent --port 8081 -d {Absolute path of the Level.sav save file}
# For example:
./pst-agent --port 8081 -d /home/lighthouse/game/Saved/
```

After confirming it's running normally, run it in the background (it will continue to run after closing the ssh window)

```bash
# Run in the background and save logs in agent.log
nohup ./pst-agent --port 8081 -d ...{manually omitted}.../Saved > agent.log 2>&1 &
# View the log
tail -f agent.log
```

#### Open Firewall/Security Group

If pst-agent and pst main body are not in the same network group, you need to open the corresponding public network port of the game server (such as 8081, or other custom ports)

#### Configuration

Enter administrator mode in the **PST main application (not pst-agent)** and open “PST Settings”. Select “pst-agent” as the save source and enter `http://game-server-address:port/sync`. The new source is used for subsequent syncs immediately after saving.

#### Close Background Operation

```bash
kill $(ps aux | grep 'pst-agent' | awk '{print $2}') | head -n 1
```

### Windows

This section covers a Windows game server with the PST main application deployed elsewhere. Follow the regular [installation guide](./README.en.md#installation), then select pst-agent as the save source in the Web settings.

#### Download

Download the pst-agent tool and rename it, e.g., rename `pst-agent_v0.10.0_windows_x86_64.exe` to `pst-agent.exe`

#### Run

Press `Win + R`, type `powershell` to open Powershell, use the `cd` command to navigate to the directory of the downloaded executable

```powershell
# .\pst-agent.exe --port Access Port -d Save file Level.sav location
.\pst-agent.exe --port 8081 -d C:\Users\ZaiGie\...\Pal\Saved
```

After successful operation, please keep the window open

#### Configuration

Enter administrator mode in the **PST main application (not pst-agent)** and open “PST Settings”. Select “pst-agent” as the save source and enter `http://game-server-address:port/sync`. The new source is used for subsequent syncs immediately after saving.
