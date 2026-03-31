#!/bin/bash

# ========================================================
# 项目名称: llm-cli 安装脚本
# 适用平台: Linux (amd64), macOS (amd64/arm64), Windows (Git Bash)
# ========================================================

# --- 配置信息 ---
GITHUB_USER="hezhengdong" 
REPO_NAME="llm-cli"
REPO_PATH="${GITHUB_USER}/${REPO_NAME}"
BIN_NAME="llm"
CONFIG_DIR="$HOME/.config/llm-cli"
CONFIG_FILE="$CONFIG_DIR/config.yaml"
INSTALL_DIR="/usr/local/bin"

echo ">>> 正在启动 llm-cli 安装程序..."

# 1. 环境检测 (操作系统与架构)
OS_TYPE="$(uname -s)"
ARCH_TYPE="$(uname -m)"

case "${OS_TYPE}" in
    Linux*)     PLATFORM=linux;;
    Darwin*)    PLATFORM=darwin;;
    CYGWIN*|MINGW*|MSYS*) PLATFORM=windows;;
    *)          echo "[错误] 不支持的操作系统: ${OS_TYPE}"; exit 1;;
esac

case "${ARCH_TYPE}" in
    x86_64*|amd64*) ARCH=amd64;;
    aarch64*|arm64*) ARCH=arm64;;
    *)              echo "[错误] 不支持的系统架构: ${ARCH_TYPE}"; exit 1;;
esac

# Windows 系统特殊处理：后缀名及安装目录
EXE_EXT=""
if [ "$PLATFORM" = "windows" ]; then
    EXE_EXT=".exe"
    INSTALL_DIR="$HOME/bin"
    mkdir -p "$INSTALL_DIR"
fi

# 2. 获取最新版本号
echo ">>> 正在从 GitHub 检索最新版本..."
LATEST_TAG=$(curl -sL "https://api.github.com/repos/${REPO_PATH}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo "[错误] 无法获取版本信息，请确认仓库地址是否正确且已发布 Release。"
    exit 1
fi
echo ">>> 检测到最新版本: ${LATEST_TAG}"

# 3. 下载二进制文件
# 下载文件名对应 CI/CD 生成的格式：llm-cli-系统-架构
DOWNLOAD_URL="https://github.com/${REPO_PATH}/releases/download/${LATEST_TAG}/llm-cli-${PLATFORM}-${ARCH}${EXE_EXT}"
TMP_FILE="/tmp/${BIN_NAME}_download${EXE_EXT}"

echo ">>> 正在下载: ${DOWNLOAD_URL}"
curl -L -o "$TMP_FILE" "$DOWNLOAD_URL"
if [ $? -ne 0 ]; then
    echo "[错误] 下载失败，请检查网络连接。"
    exit 1
fi

# 4. 授予执行权限并移动到目标目录
chmod +x "$TMP_FILE"
echo ">>> 正在部署二进制文件到: ${INSTALL_DIR}/${BIN_NAME}${EXE_EXT}"

if [ "$PLATFORM" = "windows" ]; then
    mv "$TMP_FILE" "${INSTALL_DIR}/${BIN_NAME}${EXE_EXT}"
else
    # 非 Windows 系统通常需要 sudo 权限写入 /usr/local/bin
    echo "提示: 正在请求管理员权限以完成安装..."
    sudo mv "$TMP_FILE" "${INSTALL_DIR}/${BIN_NAME}"
fi

# 5. macOS 安全机制处理 (Gatekeeper)
if [ "$PLATFORM" = "darwin" ]; then
    echo ">>> 检测到 macOS，正在解除系统安全拦截 (xattr)..."
    sudo xattr -d com.apple.quarantine "${INSTALL_DIR}/${BIN_NAME}" 2>/dev/null || true
fi

# 6. 配置文件初始化
echo ">>> 正在检查配置文件状态..."
if [ ! -f "$CONFIG_FILE" ]; then
    echo ">>> 未发现配置文件，正在从 GitHub 获取默认配置: ${CONFIG_FILE}"
    mkdir -p "$CONFIG_DIR"
    
    # 从 GitHub 仓库的 main 分支获取 config.yaml
    CONFIG_URL="https://raw.githubusercontent.com/hezhengdong/llm-cli/main/config.yaml"
    curl -sL -o "$CONFIG_FILE" "$CONFIG_URL"
    
    if [ $? -eq 0 ]; then
        echo ">>> 默认配置已生成，请稍后手动修改 API-KEY。"
    else
        echo "[警告] 无法从 GitHub 下载默认配置文件，请检查网络连接或手动创建配置。"
    fi
else
    echo ">>> 配置文件已存在，跳过初始化，保护用户数据。"
fi

echo "------------------------------------------------"
echo "安装完成！"
echo "您现在可以在终端任何地方输入 'llm' 来运行程序。"
echo "------------------------------------------------"