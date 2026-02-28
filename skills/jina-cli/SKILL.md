---
name: jina-cli
description: Reads web content and searches the web using Jina AI Reader API. Use when extracting content from URLs, reading social media posts (X/Twitter), or web searching for current information.
---

# jina - Web Content Reader & Search

CLI tool for reading web content and performing AI-powered web searches.

## Quick start

**Install**:
```bash
curl -fsSL https://raw.githubusercontent.com/geekjourneyx/jina-cli/main/scripts/install.sh | bash
```

**Basic usage**:
```bash
# Read a URL
jina read --url "https://example.com"

# Search the web
jina search --query "golang latest news"
```

## Commands

| Command | Purpose |
|---------|---------|
| `read` | Extract and convert content from URLs to LLM-friendly format |
| `search` | Search the web with AI-powered result processing |
| `config` | Manage settings (set/get/list/path) |

## Read command

Extract content from any URL:

```bash
# Basic read
jina read --url "https://example.com"

# Read with image captioning
jina read -u "https://x.com/user/status/123" --with-alt

# Batch process from file
jina read --file urls.txt

# Output as Markdown
jina read -u "https://example.com" --output markdown

# Save to file
jina read -u "https://example.com" --output-file result.md
```

### Response formats

The API can return content in different formats via `--format`:
- `markdown` - Default, LLM-friendly Markdown
- `html` - Raw HTML
- `text` - Plain text
- `screenshot` - URL to a screenshot

### Advanced options

```bash
# Bypass cache
jina read -u "https://example.com" --no-cache

# Use proxy
jina read -u "https://example.com" --proxy "http://proxy.com:8080"

# CSS selector extraction
jina read -u "https://example.com" --target-selector "article.main"

# Wait for element to load
jina read -u "https://example.com" --wait-for-selector "#content"

# Forward cookies
jina read -u "https://example.com" --cookie "session=abc123"

# POST method for SPA with hash routing
jina read -u "https://example.com/#/route" --post
```

## Search command

Search the web with automatic content fetching from top results:

```bash
# Basic search
jina search --query "golang latest news"

# Restrict to specific sites
jina search -q "AI developments" --site techcrunch.com --site theverge.com

# Limit results
jina search -q "climate change" --limit 10

# Output format
jina search -q "news" --output markdown
```

### Site filtering

Use multiple `--site` flags to restrict search to specific domains:
```bash
jina search -q "startup funding" --site techcrunch.com --site theverge.com --site wired.com
```

## Configuration

Config file: `~/.jina-reader/config.yaml`

**Priority**: Command args > Environment vars > Config file > Defaults

**Environment variables**:
- `JINA_API_BASE_URL` - Read API URL (default: `https://r.jina.ai/`)
- `JINA_SEARCH_API_URL` - Search API URL (default: `https://s.jina.ai/`)
- `JINA_TIMEOUT` - Request timeout in seconds (default: `30`)
- `JINA_WITH_GENERATED_ALT` - Enable image captioning (default: `false`)
- `JINA_OUTPUT_FORMAT` - Output format: json/markdown (default: `json`)
- `JINA_PROXY_URL` - Proxy server URL

**Config commands**:
```bash
# Set configuration
jina config set timeout 60
jina config set with-generated-alt true

# View configuration
jina config list
jina config get timeout
jina config path
```

## Output formats

**JSON format** (default, machine-readable):
```json
{
  "success": true,
  "data": {
    "url": "https://example.com",
    "content": "# Extracted Content\n\n...",
    "title": "Page Title"
  }
}
```

**Markdown format** (human-readable):
```bash
jina read -u "https://example.com" --output markdown
```

## Common use cases

### Reading social media posts

```bash
# X (Twitter) posts
jina read -u "https://x.com/elonmusk/status/123456" --with-alt

# The --with-alt flag enables VLM image captioning for embedded images
```

### Reading articles/blogs

```bash
# Standard article
jina read -u "https://blog.example.com/article"

# With specific format
jina read -u "https://example.com" --format text --output markdown
```

### Research workflows

```bash
# 1. Search for topic
jina search -q "quantum computing 2025" --limit 10

# 2. Read specific results
jina read --file search_results.txt
```

### Batch processing

Create a file with one URL per line:
```bash
cat > urls.txt << EOF
https://example.com/page1
https://example.com/page2
https://x.com/user/status/123
EOF

jina read --file urls.txt --output markdown
```

## Project structure

```
cli/
├── main.go              # Root command
├── read.go              # read command
├── search.go            # search command
├── config.go            # config command
└── pkg/
    ├── api/client.go    # Jina API HTTP client
    ├── config/          # Config file management
    └── output/          # JSON/Markdown formatter
```

## Implementation notes

- **Go 1.24+** required
- **Zero dependencies** except Cobra
- **Single binary** distribution
- Config stored as simple `key=value` format (no YAML library dependency)

For API details: See `cli/pkg/api/client.go`
