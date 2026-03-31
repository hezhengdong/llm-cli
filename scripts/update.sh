#!/bin/bash

# ========================================================
# 项目名称: llm-cli 更新脚本
# 逻辑说明: 拉取仓库中最新的安装脚本并执行，实现静默升级
# ========================================================

GITHUB_USER="hezhengdong"

INSTALL_SCRIPT_URL="https://raw.githubusercontent.com/$GITHUB_USER/llm-cli/main/scripts/install.sh"

echo ">>> 正在检查 llm-cli 更新..."

# 直接调用安装脚本完成覆盖安装
curl -sL "$INSTALL_SCRIPT_URL" | bash

if [ $? -eq 0 ]; then
    echo ">>> llm-cli 升级程序执行完毕。"
else
    echo "[错误] 更新过程中出现问题，请检查网络。"
    exit 1
fi