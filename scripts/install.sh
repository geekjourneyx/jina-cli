#!/bin/bash
# jina CLI 安装脚本
# 自动检测系统架构并下载对应的二进制文件

set -e

VERSION="1.0.0"
REPO="geekjourneyx/jina-cli"
GITHUB_BASE="https://github.com/${REPO}/releases/download/v${VERSION}"
RAW_BASE="https://raw.githubusercontent.com/${REPO}/v${VERSION}/scripts"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# 检测系统架构
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)

    case "$arch" in
        x86_64|amd64)
            arch="amd64"
            ;;
        aarch64|arm64)
            arch="arm64"
            ;;
        i386|i686)
            arch="386"
            ;;
        *)
            error "不支持的架构: $arch"
            ;;
    esac

    case "$os" in
        darwin)
            os="darwin"
            ;;
        linux)
            os="linux"
            ;;
        msys*|mingw*|cygwin*)
            os="windows"
            ;;
        *)
            error "不支持的操作系统: $os"
            ;;
    esac

    echo "${os}-${arch}"
}

# 检测安装目录
detect_install_dir() {
    # 检查 HOME 目录
    if [ -z "$HOME" ]; then
        error "无法确定 HOME 目录"
    fi

    # 优先使用 $HOME/.local/bin
    if [ -d "$HOME/.local/bin" ] || echo ":$PATH:" | grep -q ":$HOME/.local/bin:"; then
        echo "$HOME/.local/bin"
        return
    fi

    # 其次使用 $HOME/bin
    if [ -d "$HOME/bin" ] || echo ":$PATH:" | grep -q ":$HOME/bin:"; then
        echo "$HOME/bin"
        return
    fi

    # 默认创建 $HOME/.local/bin
    echo "$HOME/.local/bin"
}

# 检查命令是否存在
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# 下载文件
download_file() {
    local url="$1"
    local output="$2"

    if command_exists curl; then
        curl -fsSL --connect-timeout 15 --max-time 300 -o "$output" "$url"
    elif command_exists wget; then
        wget -q --timeout=300 -O "$output" "$url"
    else
        error "需要 curl 或 wget 来下载文件"
    fi
}

# 检查 PATH
check_path() {
    local install_dir="$1"
    local current_shell="${SHELL##*/}"

    if echo ":$PATH:" | grep -q ":$install_dir:"; then
        return 0
    fi

    warn "$install_dir 不在 PATH 中"

    case "$current_shell" in
        bash)
            echo "请将以下行添加到 ~/.bashrc:"
            echo "  export PATH=\"\$PATH:$install_dir\""
            echo "然后运行: source ~/.bashrc"
            ;;
        zsh)
            echo "请将以下行添加到 ~/.zshrc:"
            echo "  export PATH=\"\$PATH:$install_dir\""
            echo "然后运行: source ~/.zshrc"
            ;;
        *)
            echo "请将 $install_dir 添加到您的 PATH 环境变量中"
            ;;
    esac
}

# 主安装流程
main() {
    echo -e "${BLUE}jina CLI 安装脚本 v${VERSION}${NC}"
    echo "======================================"
    echo ""

    # 检测平台
    PLATFORM=$(detect_platform)
    info "检测到平台: $PLATFORM"

    # 确定二进制名称（Release 中的文件名包含平台后缀）
    RELEASE_BINARY="jina-${PLATFORM}"
    BINARY_NAME="jina"
    if [ "$PLATFORM" = "windows-amd64" ]; then
        RELEASE_BINARY="jina-windows-amd64.exe"
        BINARY_NAME="jina.exe"
    fi

    # 确定安装目录
    INSTALL_DIR=$(detect_install_dir)
    info "安装目录: $INSTALL_DIR"

    # 创建安装目录
    if [ ! -d "$INSTALL_DIR" ]; then
        info "创建安装目录: $INSTALL_DIR"
        mkdir -p "$INSTALL_DIR"
    fi

    # 下载 URL（文件名包含平台后缀）
    DOWNLOAD_URL="${GITHUB_BASE}/${RELEASE_BINARY}"
    info "下载 URL: $DOWNLOAD_URL"

    # 临时文件
    TEMP_FILE="/tmp/jina-installer-$$"

    # 下载二进制
    info "下载二进制文件..."
    if ! download_file "$DOWNLOAD_URL" "$TEMP_FILE"; then
        error "下载失败，请检查网络连接或手动下载"
    fi

    # 检查文件大小
    if [ ! -s "$TEMP_FILE" ]; then
        error "下载的文件为空"
    fi

    FILE_SIZE=$(stat -c%s "$TEMP_FILE" 2>/dev/null || stat -f%z "$TEMP_FILE" 2>/dev/null)
    if [ "$FILE_SIZE" -lt 102400 ]; then
        error "下载的文件太小 ($FILE_SIZE bytes)，可能不完整"
    fi

    # 设置可执行权限
    chmod +x "$TEMP_FILE"

    # 安装
    info "安装到 $INSTALL_DIR/$BINARY_NAME"
    mv "$TEMP_FILE" "$INSTALL_DIR/$BINARY_NAME"

    # 清理
    rm -f "$TEMP_FILE"

    echo ""
    info "✓ 安装成功！"
    echo ""
    echo "二进制位置: $INSTALL_DIR/$BINARY_NAME"

    # 检查 PATH
    check_path "$INSTALL_DIR"

    echo ""
    info "快速开始:"
    echo "  jina read --url \"https://example.com\""
    echo "  jina search --query \"golang latest news\""
    echo "  jina config list"
}

# 运行安装
main "$@"
