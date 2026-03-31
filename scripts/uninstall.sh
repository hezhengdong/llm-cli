#!/bin/bash

# 判断系统环境
INSTALL_DIR="/usr/local/bin"
TARGET_PATH="$INSTALL_DIR/llm"

if echo "$(uname -s)" | grep -qE "(MINGW|CYGWIN|MSYS)"; then
    TARGET_PATH="${TARGET_PATH}.exe"
fi

# 提权删除二进制文件
SUDO=""
if [ ! -w "$INSTALL_DIR" ] && [ "$(uname -s)" != "MINGW"* ]; then
    SUDO="sudo"
fi

echo "🗑️ 正在删除 llm 可执行文件..."
$SUDO rm -f "$TARGET_PATH"

# 询问用户是否删除配置文件
read -p "❓ 是否同时删除配置文件目录 (~/.config/llm-cli)? [y/N] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    rm -rf "$HOME/.config/llm-cli"
    echo "✅ 配置文件已删除。"
fi

echo "👋 卸载完成。"