#!/bin/bash
set -e

# 为了代码复用并避免误碰配置，我们可以直接调用网络上的 install.sh
# 这里采用独立的更新逻辑，只覆盖二进制文件
echo "🔄 正在更新 llm-cli..."
curl -sSL https://raw.githubusercontent.com/<YOUR_GITHUB_USERNAME>/llm-cli/main/scripts/install.sh | bash
echo "✅ 更新成功！"