#!/bin/bash
set -e

# === 1. 探测操作系统和架构 ===
OS="$(uname -s)"
ARCH="$(uname -m)"

if [ "$OS" = "Linux" ]; then
    TARGET_OS="linux"
elif [ "$OS" = "Darwin" ]; then
    TARGET_OS="darwin"
elif echo "$OS" | grep -qE "(MINGW|CYGWIN|MSYS)"; then
    TARGET_OS="windows"
else
    echo "❌ 不支持的操作系统: $OS"
    exit 1
fi

if [ "$ARCH" = "x86_64" ] || [ "$ARCH" = "amd64" ]; then
    TARGET_ARCH="amd64"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then
    TARGET_ARCH="arm64"
else
    echo "❌ 不支持的架构: $ARCH"
    exit 1
fi

FILENAME="llm-cli-${TARGET_OS}-${TARGET_ARCH}"
if [ "$TARGET_OS" = "windows" ]; then
    FILENAME="${FILENAME}.exe"
fi

# === 2. 下载最新的二进制文件 ===
REPO="<YOUR_GITHUB_USERNAME>/llm-cli" # TODO: 替换为你的 GitHub 仓库
LATEST_URL="https://github.com/$REPO/releases/latest/download/$FILENAME"

echo "⬇️ 正在从 $LATEST_URL 下载最新版本..."
TMP_FILE=$(mktemp)
curl -sL "$LATEST_URL" -o "$TMP_FILE"

# === 3. 安装到系统路径 ===
INSTALL_DIR="/usr/local/bin"
# 如果没有写入权限，尝试使用 sudo
SUDO=""
if [ ! -w "$INSTALL_DIR" ] && [ "$TARGET_OS" != "windows" ]; then
    SUDO="sudo"
    echo "🔑 需要管理员权限将 llm 安装到 $INSTALL_DIR"
fi

$SUDO mkdir -p "$INSTALL_DIR"
TARGET_PATH="$INSTALL_DIR/llm"
if [ "$TARGET_OS" = "windows" ]; then
    TARGET_PATH="${TARGET_PATH}.exe"
fi

$SUDO mv "$TMP_FILE" "$TARGET_PATH"
$SUDO chmod +x "$TARGET_PATH"

# === 4. 处理 Mac 安全机制 (解除 quarantine 拦截) ===
if [ "$TARGET_OS" = "darwin" ]; then
    $SUDO xattr -d com.apple.quarantine "$TARGET_PATH" 2>/dev/null || true
fi

# === 5. 初始化配置文件 ===
CONFIG_DIR="$HOME/.config/llm-cli"
CONFIG_FILE="$CONFIG_DIR/config.yaml"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "⚙️ 正在创建默认配置文件..."
    mkdir -p "$CONFIG_DIR"
    cat << 'EOF' > "$CONFIG_FILE"
default_provider: moonshot
default_prompt: default

providers:
  moonshot:
    api_key: ${API-KEY}
    base_url: https://api.moonshot.cn/v1
    model: kimi-k2-turbo-preview
  openrouter:
    api_key: ${API-KEY}
    base_url: https://openrouter.ai/api/v1
    model: openrouter/free

tools:
  tavily:
    api_key: ${API-KEY}

prompts:
  default: "你是AI助手, 回复精炼, 一针见血, 200字以内"
EOF
    echo "✅ 配置文件已创建在: $CONFIG_FILE"
else
    echo "✅ 检测到已存在配置文件，跳过初始化。"
fi

echo "🎉 安装完成！你可以直接在终端中运行 'llm' 命令了。"