# 自助编译及更新指南（linux）

## 更新失效问题说明

1.pst对于存档解析的关键工具是palworld-save-tools，负责解析sav存档文件，更新后存档不能解析就是由于这个项目没有适配新存档，进而影响了基于这个项目生成的sav_cli文件。所以只要更新了新版本的palworld-save-tools，存档解析功能就能正常使用。另外，sav_cli无法解析只影响pst右侧的帕鲁/物品展示，不影响其他功能。

2.记录物品和帕鲁名称的文件在web/src/assets文件夹中，分别为items.json和pal.json，对应的图片文件夹是items和pals。新帕鲁和新物品显示异常，参考文件内其他例子，更新这两个文件即可。

下面是具体编译和更新的方法，使用系统版本：ubuntu22.04

注意：执行编译的系统版本尽量选择与运行pst的系统一直的版本，不然编译出来的程序复制过去运行可能会报错。

## 编译方法

## 1.下载源码文件

`git clone https://github.com/zaigie/palworld-server-tool.git`

## 2.编译sav_cli

进入module文件夹，将requirements.txt文件中的palworld-save-tools==0.23.0改到最新版本（如0.23.1），最新版本可以在[pal-save-tools](https://github.com/cheahjs/palworld-save-tools)找到。之后执行如下命令：
`cd module
sudo chmod 777 build.sh
./build.sh`

## 3.编译项目

回到palworld-save-tools文件夹
`cd ..
make init
make build-pub`

## 可能遇到的错误和解决方法

### 报错make: go: 没有那个文件或目录

解决办法：
安装go环境
参考这里https://golang.google.cn/doc/install
如果 go mod download下不动，运行这一句
`go env -w GOPROXY=https://goproxy.io,direct`

### 报错/bin/sh: 1: pnpm: not found

解决办法：
运行下面代码
`curl -fsSL https://get.pnpm.io/install.sh | sh - 
source /home/ubuntu/.bashrc`

### 报错pnpm: not found: node

解决办法：
`curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
source /home/ubuntu/.bashrc
nvm install 18.20.4
nvm alias default 18.20.4`

### 报错FATAL ERROR: Ineffective mark-compacts near heap limit Allocation failed - JavaScript heap out of memory

解决办法：
运行以下命令，不一定8192，反正搞大点
`export NODE_OPTIONS="--max-old-space-size=8192"`

### 报错./build.sh: 行 3: pyinstaller: 未找到命令

解决办法：
`sudo apt install python3-pip
pip3 install --upgrade pip
pip3 install pyinstaller`
然后命令行输入`pip3 show pyinstaller`
记下这一行Location: /home/ubuntu/.local/lib/python3.10/site-packages
然后在.profile文件最后加上这一行PATH=$PATH:/home/ubuntu/.local/lib/python3.10/site-packages
最后运行这一句
`source /home/ubuntu/.profile`