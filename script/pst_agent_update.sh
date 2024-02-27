#!/bin/bash

# GitHub仓库用户名和项目名
USER="zaigie"
REPO="palworld-server-tool"

#检测包管理工具
if command -v yum >/dev/null 2>&1; then
    yum install -y jq
elif command -v dnf >/dev/null 2>&1; then
    dnf install -y jq
else
    # 检查apt是否存在并可用
    if command -v apt >/dev/null 2>&1; then
        apt install -y jq
    else
        # 如果没有找到任何包管理工具，可以执行其他操作或退出脚本
	exit
    fi
fi

# 获取最新Release的信息
TAGNAME=$(curl -s "https://api.github.com/repos/$USER/$REPO/releases" | jq '.[0].tag_name' --raw-output)
#获取系统架构
MACHINE_TYPE=$(uname -m)
DOWNLOAD_URL="https://github.com/$USER/$REPO/releases/download/$TAGNAME/pst-agent_${TAGNAME}_linux_$MACHINE_TYPE"
echo "$DOWNLOAD_URL"
# 下载文件
curl -o "pst-agent" "$DOWNLOAD_URL"

# 检查下载是否成功
if [ $? -eq 0 ]; then
	    echo "最新版本 ($TAGNAME) 已成功下载"
    else
	    echo "下载最新版本失败，请检查网络连接和GitHub API访问情况。"
fi
