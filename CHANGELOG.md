# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-02-28

### Added
- Initial release of jina-cli
- `read` command for extracting content from URLs
- `search` command for AI-powered web search
- `config` command for configuration management
- JSON and Markdown output formats
- Batch URL processing from file
- Configuration file support (`~/.jina-reader/config.yaml`)
- Environment variable overrides
- Image captioning support with VLM
- Proxy support
- CSS selector-based content extraction
- Cookie forwarding
- POST method for SPA applications
- One-line installation script
- Comprehensive test coverage (70%+)

### Features
- Read URLs and convert to LLM-friendly formats (Markdown/HTML/Text)
- Search the web with automatic content fetching from top 5 results
- Site-restricted search
- Cache bypass support
- Configurable timeout settings
- Sensitive data masking in config display

### Documentation
- Bilingual README (Chinese/English)
- CLAUDE.md with development workflow
- MIT License

[1.0.0]: https://github.com/geekjourneyx/jina-cli/releases/tag/v1.0.0
