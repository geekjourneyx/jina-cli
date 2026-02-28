<div align="center">

# jina

### 为 AI Agent 打造的网页阅读 CLI 工具

[English](#english) | [中文](#中文)

</div>

---

## 中文

### 简介

`jina` 是一个轻量级的命令行工具，封装了 [Jina AI Reader API](https://github.com/jina-ai/reader)，可以将任意 URL 转换为 LLM 友好的输入格式。特别适合阅读 X (Twitter)、博客、新闻网站等复杂网页。

**核心功能：**

- **read** - 读取 URL 内容，输出 Markdown/Text/HTML 格式
- **search** - 网络搜索，自动获取并处理搜索结果

---

## 安装指南

jina 有两种独立的安装方式，根据你的使用场景选择：

### 方式一：Claude Code Skill（推荐用于 AI 辅助开发）

**适用场景**：在 Claude Code 中使用 AI 协助你处理网页内容

**安装步骤**：

```bash
# 1. 确保已安装 Node.js 和 Claude Code
node --version
# v18.0.0 或更高版本

# 2. 安装 jina skill
npx skills add https://github.com/geekjourneyx/jina-cli --skill jina-cli
```

**安装后验证**：

```bash
# 在 Claude Code 中可以直接使用
# 无需额外操作，skill 会自动加载
```

**你将获得**：
- ✅ 在 Claude Code 中直接调用 `jina read` 和 `jina search` 命令
- ✅ AI 会自动理解 jina 的功能和使用方式
- ✅ 无需手动安装 CLI 二进制文件

---

### 方式二：CLI 二进制文件（推荐用于终端使用）

**适用场景**：在终端/脚本中使用，或与其他工具集成

#### 一键安装（Linux/macOS）

```bash
curl -fsSL https://raw.githubusercontent.com/geekjourneyx/jina-cli/main/scripts/install.sh | bash
```

安装脚本会自动：
1. 检测你的系统架构（Linux/macOS，amd64/arm64）
2. 下载对应平台的二进制文件
3. 安装到 `~/.local/bin/jina` 或 `~/bin/jina`
4. 提示如何添加到 PATH（如需要）

#### 验证安装

```bash
# 检查是否安装成功
jina --version
# 预期输出: jina version 1.0.0 (构建时间: ..., 提交: ...)

# 测试基本功能
jina read --url "https://example.com"
```

#### 手动安装

如果自动安装失败，可以手动下载：

```bash
# 1. 下载对应平台的二进制
# Linux amd64:
wget https://github.com/geekjourneyx/jina-cli/releases/latest/download/jina-linux-amd64 -O jina
chmod +x jina
sudo mv jina /usr/local/bin/

# macOS ARM64 (Apple Silicon):
wget https://github.com/geekjourneyx/jina-cli/releases/latest/download/jina-darwin-arm64 -O jina
chmod +x jina
sudo mv jina /usr/local/bin/

# 2. 验证安装
jina --version
```

#### 从源码构建

```bash
git clone https://github.com/geekjourneyx/jina-cli.git
cd jina-cli
go build -o jina ./cli
sudo mv jina /usr/local/bin/
```

---

### 两种方式的区别

| 特性 | Claude Code Skill | CLI 二进制 |
|------|-------------------|------------|
| **安装位置** | `~/.claude/skills/` | `~/.local/bin/jina` |
| **使用环境** | 仅 Claude Code | 任何终端/脚本 |
| **AI 集成** | AI 自动理解功能 | 需要手动调用 |
| **更新方式** | `npx skills update` | 重新运行安装脚本 |
| **适用场景** | AI 辅助开发 | 自动化脚本、日常使用 |

**注意**：两种方式完全独立，可以同时安装，互不干扰。

---

### 快速开始

#### 读取网页内容

```bash
# 读取单个 URL
jina read --url "https://example.com"

# 读取 X (Twitter) 帖子
jina read -u "https://x.com/user/status/123456789" --with-alt

# 输出 Markdown 格式
jina read -u "https://example.com" --output markdown

# 保存到文件
jina read -u "https://example.com" --output-file result.md
```

#### 批量处理

```bash
# 从文件读取 URL 列表
cat > urls.txt << EOF
https://example.com/page1
https://example.com/page2
https://x.com/user/status/123
EOF

jina read --file urls.txt
```

#### 网络搜索

```bash
# 搜索关键词
jina search --query "golang latest news"

# 限定搜索站点
jina search -q "AI developments" --site techcrunch.com --site theverge.com

# 限制结果数量
jina search -q "climate change" --limit 10
```

### 配置管理

配置文件位于 `~/.jina-reader/config.yaml`：

```bash
# 查看所有配置
jina config list

# 设置配置项
jina config set timeout 60
jina config set with-generated-alt true

# 获取单个配置
jina config get timeout

# 查看配置文件路径
jina config path
```

### 配置项说明

| 配置项 | 环境变量 | 默认值 | 说明 |
|--------|----------|--------|------|
| `api_base_url` | `JINA_API_BASE_URL` | `https://r.jina.ai/` | Read API 地址 |
| `search_api_url` | `JINA_SEARCH_API_URL` | `https://s.jina.ai/` | Search API 地址 |
| `default_response_format` | `JINA_RESPONSE_FORMAT` | `markdown` | 响应格式 |
| `default_output_format` | `JINA_OUTPUT_FORMAT` | `json` | 输出格式 |
| `timeout` | `JINA_TIMEOUT` | `30` | 请求超时（秒） |
| `with_generated_alt` | `JINA_WITH_GENERATED_ALT` | `false` | 启用图片描述 |
| `proxy_url` | `JINA_PROXY_URL` | `""` | 代理服务器 |

**优先级：** 命令行参数 > 环境变量 > 配置文件 > 默认值

### 输出格式

#### JSON 格式（默认）

```json
{
  "success": true,
  "data": {
    "url": "https://example.com",
    "title": "Example Domain",
    "content": "# Example Domain\n\nThis domain is..."
  }
}
```

#### Markdown 格式

```bash
jina read -u "https://example.com" --output markdown
```

输出：
```markdown
# Example Domain

**Source**: https://example.com

---

# Example Domain

This domain is for use in illustrative examples...
```

### 高级用法

#### 禁用缓存

```bash
jina read -u "https://example.com" --no-cache
```

#### 使用代理

```bash
jina read -u "https://example.com" --proxy "http://proxy.com:8080"
```

#### CSS 选择器提取

```bash
# 只提取指定元素
jina read -u "https://example.com" --target-selector "article.main"

# 等待元素加载
jina read -u "https://example.com" --wait-for-selector "#content"
```

#### 处理 SPA 应用

```bash
# 对于带 hash 路由的 SPA，使用 POST 方法
jina read -u "https://example.com/#/route" --post
```

#### 设置 Cookie

```bash
jina read -u "https://example.com" --cookie "session=abc123"
```

### 与 AI Agent 集成

```bash
# 在 Claude Code 或其他 AI 工具中使用
jina read -u "https://x.com/elonmusk/status/123456" --output json

# 搜索最新信息
jina search -q "golang 1.24 release notes" --output json
```

### 命令参考

```
jina - CLI tool for Jina AI Reader and Search APIs

Usage:
  jina [command]

Available Commands:
  read        Extract and convert content from URLs
  search      Search the web with AI-powered results
  config      Manage configuration
  completion  Generate shell completion
  help        Help about any command

Flags:
  -a, --api-base string   API base URL (overrides config)
  -k, --api-key string    API key (overrides config)
  -o, --output string     Output format: json, markdown (default "json")
  -v, --verbose           Verbose output
  -h, --help              help for jina
      --version           version for jina
```

### 开发

```bash
# 运行测试
go test ./...

# 运行测试并计算覆盖率
go test -cover ./...

# 构建
go build -o jina ./cli

# 交叉编译
GOOS=linux GOARCH=amd64 go build -o jina-linux-amd64 ./cli
GOOS=darwin GOARCH=arm64 go build -o jina-darwin-arm64 ./cli
GOOS=windows GOARCH=amd64 go build -o jina-windows-amd64.exe ./cli
```

### 项目结构

```
jina-cli/
├── cli/
│   ├── main.go          # 入口
│   ├── read.go          # read 命令
│   ├── search.go        # search 命令
│   ├── config.go        # config 命令
│   └── pkg/
│       ├── api/         # HTTP 客户端
│       ├── config/      # 配置管理
│       └── output/      # 输出格式化
└── scripts/
    └── install.sh       # 安装脚本
```

### 许可证

[MIT License](LICENSE)

### 致谢

- [Jina AI Reader API](https://github.com/jina-ai/reader) - 核心 API 服务
- [md2wechat-lite](https://github.com/geekjourneyx/md2wechat-lite) - 项目架构参考

---

## English

### Overview

`jina` is a lightweight CLI tool that wraps the [Jina AI Reader API](https://github.com/jina-ai/reader) to convert any URL into LLM-friendly input. Perfect for reading X (Twitter), blogs, news sites, and other complex web pages.

**Features:**

- **read** - Extract content from URLs in Markdown/Text/HTML format
- **search** - Search the web with AI-powered result processing

---

## Installation Guide

jina offers two independent installation methods. Choose based on your use case:

### Method 1: Claude Code Skill (Recommended for AI-Assisted Development)

**Use case**: Using jina within Claude Code with AI assistance

**Installation steps**:

```bash
# 1. Ensure Node.js and Claude Code are installed
node --version
# v18.0.0 or higher

# 2. Install jina skill
npx skills add https://github.com/geekjourneyx/jina-cli --skill jina-cli
```

**Verify installation**:

```bash
# You can now use jina commands directly in Claude Code
# No additional steps needed, skill loads automatically
```

**You get**:
- ✅ Direct access to `jina read` and `jina search` commands in Claude Code
- ✅ AI automatically understands jina's functionality
- ✅ No manual CLI binary installation required

---

### Method 2: CLI Binary (Recommended for Terminal/Scripting)

**Use case**: Using in terminal/scripts, or integrating with other tools

#### One-line Installation (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/geekjourneyx/jina-cli/main/scripts/install.sh | bash
```

The installation script automatically:
1. Detects your system architecture (Linux/macOS, amd64/arm64)
2. Downloads the appropriate binary for your platform
3. Installs to `~/.local/bin/jina` or `~/bin/jina`
4. Prompts to add to PATH if needed

#### Verify Installation

```bash
# Check if installation succeeded
jina --version
# Expected output: jina version 1.0.0 (build: ..., commit: ...)

# Test basic functionality
jina read --url "https://example.com"
```

#### Manual Installation

If auto-install fails, download manually:

```bash
# 1. Download binary for your platform
# Linux amd64:
wget https://github.com/geekjourneyx/jina-cli/releases/latest/download/jina-linux-amd64 -O jina
chmod +x jina
sudo mv jina /usr/local/bin/

# macOS ARM64 (Apple Silicon):
wget https://github.com/geekjourneyx/jina-cli/releases/latest/download/jina-darwin-arm64 -O jina
chmod +x jina
sudo mv jina /usr/local/bin/

# 2. Verify installation
jina --version
```

#### Build from Source

```bash
git clone https://github.com/geekjourneyx/jina-cli.git
cd jina-cli
go build -o jina ./cli
sudo mv jina /usr/local/bin/
```

---

### Comparison

| Feature | Claude Code Skill | CLI Binary |
|---------|-------------------|------------|
| **Install location** | `~/.claude/skills/` | `~/.local/bin/jina` |
| **Environment** | Claude Code only | Any terminal/script |
| **AI integration** | AI understands functionality | Manual invocation |
| **Updates** | `npx skills update` | Re-run install script |
| **Best for** | AI-assisted development | Automation, daily use |

**Note**: The two methods are completely independent. You can install both without conflicts.

---

### Quick Start

#### Read Web Content

```bash
# Read a single URL
jina read --url "https://example.com"

# Read X (Twitter) post
jina read -u "https://x.com/user/status/123456789" --with-alt

# Output as Markdown
jina read -u "https://example.com" --output markdown

# Save to file
jina read -u "https://example.com" --output-file result.md
```

#### Batch Processing

```bash
# Read from URL list file
cat > urls.txt << EOF
https://example.com/page1
https://example.com/page2
https://x.com/user/status/123
EOF

jina read --file urls.txt
```

#### Web Search

```bash
# Search for a query
jina search --query "golang latest news"

# Restrict to specific sites
jina search -q "AI developments" --site techcrunch.com --site theverge.com

# Limit results
jina search -q "climate change" --limit 10
```

### Configuration

Config file location: `~/.jina-reader/config.yaml`

```bash
# List all configuration
jina config list

# Set configuration
jina config set timeout 60
jina config set with-generated-alt true

# Get single configuration
jina config get timeout

# Show config file path
jina config path
```

### Configuration Options

| Key | Env Var | Default | Description |
|-----|---------|---------|-------------|
| `api_base_url` | `JINA_API_BASE_URL` | `https://r.jina.ai/` | Read API URL |
| `search_api_url` | `JINA_SEARCH_API_URL` | `https://s.jina.ai/` | Search API URL |
| `default_response_format` | `JINA_RESPONSE_FORMAT` | `markdown` | Response format |
| `default_output_format` | `JINA_OUTPUT_FORMAT` | `json` | Output format |
| `timeout` | `JINA_TIMEOUT` | `30` | Request timeout (seconds) |
| `with_generated_alt` | `JINA_WITH_GENERATED_ALT` | `false` | Enable image captioning |
| `proxy_url` | `JINA_PROXY_URL` | `""` | Proxy server |

**Priority:** CLI args > Env vars > Config file > Defaults

### Output Formats

#### JSON Format (Default)

```json
{
  "success": true,
  "data": {
    "url": "https://example.com",
    "title": "Example Domain",
    "content": "# Example Domain\n\nThis domain is..."
  }
}
```

#### Markdown Format

```bash
jina read -u "https://example.com" --output markdown
```

Output:
```markdown
# Example Domain

**Source**: https://example.com

---

# Example Domain

This domain is for use in illustrative examples...
```

### Advanced Usage

#### Bypass Cache

```bash
jina read -u "https://example.com" --no-cache
```

#### Use Proxy

```bash
jina read -u "https://example.com" --proxy "http://proxy.com:8080"
```

#### CSS Selector Extraction

```bash
# Extract specific element only
jina read -u "https://example.com" --target-selector "article.main"

# Wait for element to load
jina read -u "https://example.com" --wait-for-selector "#content"
```

#### Handle SPA Apps

```bash
# For SPA with hash routing, use POST method
jina read -u "https://example.com/#/route" --post
```

#### Set Cookie

```bash
jina read -u "https://example.com" --cookie "session=abc123"
```

### AI Agent Integration

```bash
# Use with Claude Code or other AI tools
jina read -u "https://x.com/elonmusk/status/123456" --output json

# Search for latest information
jina search -q "golang 1.24 release notes" --output json
```

### Command Reference

```
jina - CLI tool for Jina AI Reader and Search APIs

Usage:
  jina [command]

Available Commands:
  read        Extract and convert content from URLs
  search      Search the web with AI-powered results
  config      Manage configuration
  completion  Generate shell completion
  help        Help about any command

Flags:
  -a, --api-base string   API base URL (overrides config)
  -k, --api-key string    API key (overrides config)
  -o, --output string     Output format: json, markdown (default "json")
  -v, --verbose           Verbose output
  -h, --help              help for jina
      --version           version for jina
```

### Development

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Build
go build -o jina ./cli

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o jina-linux-amd64 ./cli
GOOS=darwin GOARCH=arm64 go build -o jina-darwin-arm64 ./cli
GOOS=windows GOARCH=amd64 go build -o jina-windows-amd64.exe ./cli
```

### License

[MIT License](LICENSE)

### Acknowledgments

- [Jina AI Reader API](https://github.com/jina-ai/reader) - Core API service
- [md2wechat-lite](https://github.com/geekjourneyx/md2wechat-lite) - Architecture reference
