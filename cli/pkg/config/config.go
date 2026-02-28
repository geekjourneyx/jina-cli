// Package config 提供配置文件管理功能。
//
// 配置文件位置: ~/.jina-reader/config.yaml
//
// 支持的配置项:
//   - api_base_url: Read API 基础 URL
//   - search_api_url: Search API URL
//   - default_response_format: 默认响应格式 (markdown/html/text/screenshot)
//   - default_output_format: 默认输出格式 (json/markdown)
//   - timeout: 请求超时时间（秒）
//   - with_generated_alt: 启用图片描述
//   - proxy_url: 代理服务器 URL
//   - cache_tolerance: 缓存容忍度（秒）
//   - api_key: API 密钥
//
// 配置优先级: 命令行参数 > 环境变量 > 配置文件 > 默认值
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

const (
	// DefaultAPIBaseURL 默认 Read API 基础 URL
	DefaultAPIBaseURL = "https://r.jina.ai/"
	// DefaultSearchAPIURL 默认 Search API URL
	DefaultSearchAPIURL = "https://s.jina.ai/"
	// ConfigDir 配置目录名
	ConfigDir = ".jina-reader"
	// ConfigFile 配置文件名
	ConfigFile = "config.yaml"
)

var (
	// configPath 配置文件完整路径
	configPath string
	// configDir 配置目录路径
	configDir string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	configDir = filepath.Join(homeDir, ConfigDir)
	configPath = filepath.Join(configDir, ConfigFile)
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() string {
	return configPath
}

// GetConfigDir 获取配置目录路径
func GetConfigDir() string {
	return configDir
}

// Config 应用配置
type Config struct {
	ReadAPIURL            string
	SearchAPIURL          string
	DefaultResponseFormat string
	DefaultOutputFormat   string
	Timeout               int
	WithGeneratedAlt      bool
	ProxyURL              string
	CacheTolerance        string
	APIKey                string
}

// Load 从配置文件加载配置
func Load() (*Config, error) {
	cfg := &Config{
		ReadAPIURL:            DefaultAPIBaseURL,
		SearchAPIURL:          DefaultSearchAPIURL,
		DefaultResponseFormat: "markdown",
		DefaultOutputFormat:   "json",
		Timeout:               30,
		WithGeneratedAlt:      false,
	}

	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return applyEnvOverrides(cfg), nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 简单的 key=value 解析（为了保持轻量，不依赖 yaml 库）
	// 格式: key=value
	lines := splitLines(data)
	for _, line := range lines {
		line = trimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}
		key, value, ok := parseKeyValue(line)
		if !ok {
			continue
		}
		switch key {
		case "api_base_url":
			cfg.ReadAPIURL = value
		case "search_api_url":
			cfg.SearchAPIURL = value
		case "default_response_format":
			cfg.DefaultResponseFormat = value
		case "default_output_format":
			cfg.DefaultOutputFormat = value
		case "timeout":
			if timeout, err := strconv.Atoi(value); err == nil {
				cfg.Timeout = timeout
			}
		case "with_generated_alt":
			cfg.WithGeneratedAlt = strings.ToLower(value) == "true" || value == "1"
		case "proxy_url":
			cfg.ProxyURL = value
		case "cache_tolerance":
			cfg.CacheTolerance = value
		case "api_key":
			cfg.APIKey = value
		}
	}

	return applyEnvOverrides(cfg), nil
}

// applyEnvOverrides 应用环境变量覆盖
func applyEnvOverrides(cfg *Config) *Config {
	if v := os.Getenv("JINA_API_BASE_URL"); v != "" {
		cfg.ReadAPIURL = v
	}
	if v := os.Getenv("JINA_SEARCH_API_URL"); v != "" {
		cfg.SearchAPIURL = v
	}
	if v := os.Getenv("JINA_RESPONSE_FORMAT"); v != "" {
		cfg.DefaultResponseFormat = v
	}
	if v := os.Getenv("JINA_OUTPUT_FORMAT"); v != "" {
		cfg.DefaultOutputFormat = v
	}
	if v := os.Getenv("JINA_TIMEOUT"); v != "" {
		if timeout, err := strconv.Atoi(v); err == nil {
			cfg.Timeout = timeout
		}
	}
	if v := os.Getenv("JINA_WITH_GENERATED_ALT"); v != "" {
		cfg.WithGeneratedAlt = strings.ToLower(v) == "true" || v == "1"
	}
	if v := os.Getenv("JINA_PROXY_URL"); v != "" {
		cfg.ProxyURL = v
	}
	if v := os.Getenv("JINA_CACHE_TOLERANCE"); v != "" {
		cfg.CacheTolerance = v
	}
	if v := os.Getenv("JINA_API_KEY"); v != "" {
		cfg.APIKey = v
	}
	return cfg
}

// Save 保存配置到文件
func Save(cfg *Config) error {
	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 构建配置内容
	var content string
	content += "# jina-reader 配置文件\n"
	content += "# 可通过环境变量覆盖（优先级更高）\n"
	content += "#\n"
	content += "# 配置项说明：\n"
	content += "#   api_base_url             - Read API 基础 URL（默认：https://r.jina.ai/）\n"
	content += "#   search_api_url           - Search API URL（默认：https://s.jina.ai/）\n"
	content += "#   default_response_format   - 响应格式（默认：markdown）\n"
	content += "#     可选值: markdown, html, text, screenshot\n"
	content += "#   default_output_format    - 输出格式（默认：json）\n"
	content += "#     可选值: json, markdown\n"
	content += "#   timeout                  - 请求超时时间，单位：秒（默认：30）\n"
	content += "#   with_generated_alt       - 启用图片描述（默认：false）\n"
	content += "#   proxy_url                - 代理服务器 URL\n"
	content += "#   cache_tolerance          - 缓存容忍度，单位：秒\n"
	content += "#   api_key                  - API 密钥（如果需要）\n"
	content += "#\n\n"

	if cfg.ReadAPIURL != "" && cfg.ReadAPIURL != DefaultAPIBaseURL {
		content += fmt.Sprintf("api_base_url=%s\n", cfg.ReadAPIURL)
	}
	if cfg.SearchAPIURL != "" && cfg.SearchAPIURL != DefaultSearchAPIURL {
		content += fmt.Sprintf("search_api_url=%s\n", cfg.SearchAPIURL)
	}
	if cfg.DefaultResponseFormat != "" && cfg.DefaultResponseFormat != "markdown" {
		content += fmt.Sprintf("default_response_format=%s\n", cfg.DefaultResponseFormat)
	}
	if cfg.DefaultOutputFormat != "" && cfg.DefaultOutputFormat != "json" {
		content += fmt.Sprintf("default_output_format=%s\n", cfg.DefaultOutputFormat)
	}
	if cfg.Timeout != 30 {
		content += fmt.Sprintf("timeout=%d\n", cfg.Timeout)
	}
	if cfg.WithGeneratedAlt {
		content += fmt.Sprintf("with_generated_alt=%t\n", cfg.WithGeneratedAlt)
	}
	if cfg.ProxyURL != "" {
		content += fmt.Sprintf("proxy_url=%s\n", cfg.ProxyURL)
	}
	if cfg.CacheTolerance != "" {
		content += fmt.Sprintf("cache_tolerance=%s\n", cfg.CacheTolerance)
	}
	if cfg.APIKey != "" {
		content += fmt.Sprintf("api_key=%s\n", cfg.APIKey)
	}

	// 写入文件
	if err := os.WriteFile(configPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// Set 设置单个配置项
func Set(key, value string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	// 支持下划线和连字符两种格式
	normalizedKey := strings.ReplaceAll(key, "-", "_")

	switch normalizedKey {
	case "api_base_url":
		cfg.ReadAPIURL = value
	case "search_api_url":
		cfg.SearchAPIURL = value
	case "default_response_format":
		cfg.DefaultResponseFormat = value
	case "default_output_format":
		cfg.DefaultOutputFormat = value
	case "timeout":
		if timeout, err := strconv.Atoi(value); err != nil {
			return fmt.Errorf("无效的超时值: %s", value)
		} else {
			cfg.Timeout = timeout
		}
	case "with_generated_alt":
		cfg.WithGeneratedAlt = strings.ToLower(value) == "true" || value == "1"
	case "proxy_url":
		cfg.ProxyURL = value
	case "cache_tolerance":
		cfg.CacheTolerance = value
	case "api_key":
		cfg.APIKey = value
	default:
		return fmt.Errorf("未知的配置项: %s", key)
	}

	return Save(cfg)
}

// Get 获取单个配置项
func Get(key string) (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", err
	}

	// 支持下划线和连字符两种格式
	normalizedKey := strings.ReplaceAll(key, "-", "_")

	switch normalizedKey {
	case "api_base_url":
		return cfg.ReadAPIURL, nil
	case "search_api_url":
		return cfg.SearchAPIURL, nil
	case "default_response_format":
		return cfg.DefaultResponseFormat, nil
	case "default_output_format":
		return cfg.DefaultOutputFormat, nil
	case "timeout":
		return strconv.Itoa(cfg.Timeout), nil
	case "with_generated_alt":
		return strconv.FormatBool(cfg.WithGeneratedAlt), nil
	case "proxy_url":
		return cfg.ProxyURL, nil
	case "cache_tolerance":
		return cfg.CacheTolerance, nil
	case "api_key":
		if cfg.APIKey == "" {
			return "", nil
		}
		return maskSensitive(cfg.APIKey), nil
	default:
		return "", fmt.Errorf("未知的配置项: %s", key)
	}
}

// List 列出所有配置项
func List() (map[string]string, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	result["api_base_url"] = cfg.ReadAPIURL
	result["search_api_url"] = cfg.SearchAPIURL
	result["default_response_format"] = cfg.DefaultResponseFormat
	result["default_output_format"] = cfg.DefaultOutputFormat
	result["timeout"] = strconv.Itoa(cfg.Timeout)
	result["with_generated_alt"] = strconv.FormatBool(cfg.WithGeneratedAlt)
	result["proxy_url"] = cfg.ProxyURL
	result["cache_tolerance"] = cfg.CacheTolerance
	if cfg.APIKey != "" {
		result["api_key"] = maskSensitive(cfg.APIKey)
	} else {
		result["api_key"] = ""
	}

	return result, nil
}

// maskSensitive 掩码敏感信息
func maskSensitive(s string) string {
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "***" + s[len(s)-4:]
}

// 简单的字符串处理辅助函数
func splitLines(data []byte) []string {
	var lines []string
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, string(data[start:i]))
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, string(data[start:]))
	}
	return lines
}

func trimSpace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\r') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

func parseKeyValue(line string) (string, string, bool) {
	idx := slices.Index([]byte(line), '=')
	if idx == -1 {
		return "", "", false
	}
	key := trimSpace(line[:idx])
	value := trimSpace(line[idx+1:])
	return key, value, true
}
