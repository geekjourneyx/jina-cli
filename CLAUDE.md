# jina-cli - 项目记忆

## 项目概述

**jina-cli** 是一个轻量级的 CLI 工具，封装 Jina AI Reader API，为 AI Agent 提供网页内容读取和搜索能力。

- **项目名**: jina-cli
- **CLI 名**: jina
- **仓库**: https://github.com/geekjourneyx/jina-cli
- **语言**: Go 1.24+
- **框架**: Cobra
- **核心 API**: Jina AI Reader (https://github.com/jina-ai/reader)

## 项目结构

```
jina-cli/
├── cli/                    # CLI 源码
│   ├── main.go             # 程序入口
│   ├── read.go             # read 命令
│   ├── search.go           # search 命令
│   ├── config.go           # 配置管理命令
│   ├── pkg/                # 内部包
│   │   ├── api/            # API 客户端
│   │   ├── config/         # 配置文件管理
│   │   └── output/         # 输出格式化（JSON/Markdown）
│   └── scripts/            # 安装脚本
├── .github/
│   └── workflows/          # CI/CD 工作流
│       ├── ci.yml          # 持续集成
│       └── release.yml     # 发布流程
├── CLAUDE.md              # 本文件
└── README.md              # 项目说明（双语）
```

## 提交前工作流程

每次代码提交前，必须按顺序执行以下步骤：

### 1. 编译检查

```bash
cd /root/jina-cli
go build -o jina ./cli
```

确保编译无错误。

### 2. 代码检查

```bash
go test ./...
go vet ./...
```

确保测试通过，无 vet 警告。

### 3. 代码覆盖率检查

```bash
go test ./... -cover
```

目标覆盖率：
- 整体: 70%+
- 核心包 (api, config, output): 75%+

### 4. 代码格式化

```bash
go fmt ./...
goimports -w ./cli
```

确保代码格式统一。

### 5. golangci-lint 检查（如果可用）

```bash
golangci-lint run
```

修复所有 lint 问题。

### 6. 功能测试

手动测试核心功能：

```bash
# 测试 read 命令
./jina read --url "https://example.com"
./jina read -u "https://example.com" --output markdown

# 测试 search 命令
./jina search --query "test"

# 测试 config 命令
./jina config list
./jina config set timeout 60
./jina config get timeout
```

### 7. README.md 检查

检查 `README.md` 符合规范：

- [ ] 中英文双语完整
- [ ] 示例代码可运行
- [ ] 安装说明清晰
- [ ] 配置项表格完整
- [ ] 项目结构文档与实际目录结构一致

### 8. Git Commit 规范

提交信息格式：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type**:
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档变更
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具变更

**重要**: Commit 中不能有 claude 账号和邮箱信息！

正确示例：
```
feat(read): add batch URL processing from file

Support reading multiple URLs from a text file,
one URL per line.
```

### 9. 更新 CHANGELOG.md

在 `CHANGELOG.md` 顶部添加新版本：

```markdown
## [1.x.x] - 2025-02-28

### Added
- 新增功能说明

### Changed
- 变更说明

### Fixed
- 修复说明

### Removed
- 移除说明
```

### 10. 发版流程

```bash
# 1. 确认所有检查通过
# 2. 提交变更
git add .
git commit -m "chore(release): prepare for v1.x.x"

# 3. 创建标签
git tag -a v1.x.x -m "Release v1.x.x"

# 4. 推送到远程
git push origin main
git push origin v1.x.x
```

## 技术要点

### 配置管理

- 配置文件位置：`~/.jina-reader/config.yaml`
- 配置优先级：命令行参数 > 环境变量 > 配置文件 > 默认值
- 环境变量前缀：`JINA_`

### API 端点

- Read API: `https://r.jina.ai/`
- Search API: `https://s.jina.ai/`
- 官方文档: https://github.com/jina-ai/reader

### 输出格式

- **JSON**（默认）: 机器可读，适合 AI Agent 解析
- **Markdown**: 人类可读，适合终端查看

### 支持的响应格式

- `markdown` - Markdown 格式（默认）
- `html` - 原始 HTML
- `text` - 纯文本
- `screenshot` - 网页截图

### 零依赖设计

- 除 Cobra 外无其他运行时依赖
- 手动实现 key=value 配置解析（不依赖 YAML 库）
- 单一二进制文件分发

## 命令参考

### read 命令

```bash
jina read --url "https://example.com"
jina read -u "https://example.com" --with-alt --no-cache
jina read --file urls.txt --output markdown
```

关键选项：
- `-u, --url`: URL 地址
- `-f, --file`: URL 文件（批量）
- `-F, --format`: 响应格式
- `-t, --timeout`: 超时时间
- `--with-alt`: 启用图片描述
- `--no-cache`: 禁用缓存
- `--proxy`: 代理服务器
- `--post`: POST 方法（SPA）

### search 命令

```bash
jina search --query "golang news"
jina search -q "AI" --site techcrunch.com --site theverge.com
jina search -q "news" --limit 10
```

关键选项：
- `-q, --query`: 搜索关键词
- `-s, --site`: 站点过滤（可重复）
- `-l, --limit`: 结果数量
- `-F, --format`: 响应格式

### config 命令

```bash
jina config set timeout 60
jina config get timeout
jina config list
jina config path
```

## 开发规范

### 文件命名

- 使用小写字母
- 多单词用下划线分隔：`read_cmd.go`
- 测试文件：`xxx_test.go`

### 包命名

- 全小写，无下划线
- 简短描述性名称
- 与目录名一致

### 错误处理

```go
// 包装错误
return fmt.Errorf("操作失败: %w", err)

// 统一输出
output.Error(err)  // JSON 格式错误输出
```

### 配置解析

```go
// key=value 格式，无需 YAML 库
api_base_url=https://r.jina.ai/
timeout=30
with_generated_alt=false
```

## 相关链接

- Jina AI Reader: https://github.com/jina-ai/reader
- Jina AI 官网: https://jina.ai/reader
- 项目参考: https://github.com/geekjourneyx/md2wechat-lite

## 致谢

本项目架构参考 [md2wechat-lite](https://github.com/geekjourneyx/md2wechat-lite) 的开发规范和 CLI 设计模式。
