#!/bin/bash

# ========================================================
# 项目名称: llm-cli 卸载脚本
# ========================================================

BIN_NAME="llm"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.config/llm-cli"

echo ">>> 正在启动 llm-cli 卸载程序..."

# 1. 查找并移除二进制文件
BIN_PATH=$(which $BIN_NAME 2>/dev/null)

if [ -n "$BIN_PATH" ]; then
    echo ">>> 发现程序路径: ${BIN_PATH}"
    echo "提示: 正在请求管理员权限以移除程序文件..."
    sudo rm -f "$BIN_PATH"
    echo ">>> 二进制文件已成功删除。"
else
    echo ">>> 系统路径中未发现 ${BIN_NAME} 命令，尝试清理默认安装目录..."
    sudo rm -f "${INSTALL_DIR}/${BIN_NAME}"
    sudo rm -f "${INSTALL_DIR}/${BIN_NAME}.exe"
fi

# 2. 处理配置文件
if [ -d "$CONFIG_DIR" ]; then
    echo "------------------------------------------------"
    echo "检测到配置文件目录: ${CONFIG_DIR}"
    echo "如果您希望彻底删除所有数据（包括 API 配置），请输入 'y'。"
    echo "如果希望保留配置以便日后再次使用，请直接按回车。"
    read -p "是否删除配置文件？(y/N): " CONFIRM
    
    if [[ "$CONFIRM" =~ ^[Yy]$ ]]; then
        rm -rf "$CONFIG_DIR"
        echo ">>> 配置文件目录已清除。"
    else
        echo ">>> 已保留配置文件。"
    fi
fi

echo "------------------------------------------------"
echo "llm-cli 已成功从您的系统中卸载。"
echo "------------------------------------------------"