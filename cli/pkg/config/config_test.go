package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestLoad_DefaultConfig(t *testing.T) {
	// 使用临时目录
	tmpDir := t.TempDir()
	originalConfigDir := configDir
	originalConfigPath := configPath
	defer func() {
		configDir = originalConfigDir
		configPath = originalConfigPath
	}()

	configDir = tmpDir
	configPath = filepath.Join(tmpDir, ConfigFile)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// 验证默认值
	if cfg.ReadAPIURL != DefaultAPIBaseURL {
		t.Errorf("Expected ReadAPIURL %s, got %s", DefaultAPIBaseURL, cfg.ReadAPIURL)
	}
	if cfg.SearchAPIURL != DefaultSearchAPIURL {
		t.Errorf("Expected SearchAPIURL %s, got %s", DefaultSearchAPIURL, cfg.SearchAPIURL)
	}
	if cfg.DefaultResponseFormat != "markdown" {
		t.Errorf("Expected DefaultResponseFormat 'markdown', got %s", cfg.DefaultResponseFormat)
	}
	if cfg.DefaultOutputFormat != "json" {
		t.Errorf("Expected DefaultOutputFormat 'json', got %s", cfg.DefaultOutputFormat)
	}
	if cfg.Timeout != 30 {
		t.Errorf("Expected Timeout 30, got %d", cfg.Timeout)
	}
	if cfg.WithGeneratedAlt {
		t.Errorf("Expected WithGeneratedAlt false, got %t", cfg.WithGeneratedAlt)
	}
}

func TestLoad_WithConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	originalConfigDir := configDir
	originalConfigPath := configPath
	defer func() {
		configDir = originalConfigDir
		configPath = originalConfigPath
	}()

	configDir = tmpDir
	configPath = filepath.Join(tmpDir, ConfigFile)

	// 创建配置文件
	configContent := `# Test config
api_base_url=https://custom.api.com/
timeout=60
with_generated_alt=true
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.ReadAPIURL != "https://custom.api.com/" {
		t.Errorf("Expected ReadAPIURL 'https://custom.api.com/', got %s", cfg.ReadAPIURL)
	}
	if cfg.Timeout != 60 {
		t.Errorf("Expected Timeout 60, got %d", cfg.Timeout)
	}
	if !cfg.WithGeneratedAlt {
		t.Errorf("Expected WithGeneratedAlt true, got %t", cfg.WithGeneratedAlt)
	}
}

func TestLoad_WithEnvOverride(t *testing.T) {
	tmpDir := t.TempDir()
	originalConfigDir := configDir
	originalConfigPath := configPath
	defer func() {
		configDir = originalConfigDir
		configPath = originalConfigPath
		os.Unsetenv("JINA_TIMEOUT")
		os.Unsetenv("JINA_API_BASE_URL")
	}()

	configDir = tmpDir
	configPath = filepath.Join(tmpDir, ConfigFile)

	// 设置环境变量
	os.Setenv("JINA_TIMEOUT", "120")
	os.Setenv("JINA_API_BASE_URL", "https://env.api.com/")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.ReadAPIURL != "https://env.api.com/" {
		t.Errorf("Expected ReadAPIURL from env 'https://env.api.com/', got %s", cfg.ReadAPIURL)
	}
	if cfg.Timeout != 120 {
		t.Errorf("Expected Timeout from env 120, got %d", cfg.Timeout)
	}
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	originalConfigDir := configDir
	originalConfigPath := configPath
	defer func() {
		configDir = originalConfigDir
		configPath = originalConfigPath
	}()

	configDir = tmpDir
	configPath = filepath.Join(tmpDir, ConfigFile)

	cfg := &Config{
		ReadAPIURL:            "https://test.api.com/",
		SearchAPIURL:          "https://test.search.com/",
		DefaultResponseFormat: "html",
		DefaultOutputFormat:   "markdown",
		Timeout:               90,
		WithGeneratedAlt:      true,
		ProxyURL:              "http://proxy.com:8080",
		CacheTolerance:        "3600",
		APIKey:                "test-api-key-12345",
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// 读取并验证内容
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "api_base_url=https://test.api.com/") {
		t.Error("Config file missing api_base_url")
	}
	if !strings.Contains(content, "timeout=90") {
		t.Error("Config file missing timeout")
	}
	if !strings.Contains(content, "with_generated_alt=true") {
		t.Error("Config file missing with_generated_alt")
	}
	if !strings.Contains(content, "api_key=test-api-key-12345") {
		t.Error("Config file missing api_key")
	}
}

func TestSet(t *testing.T) {
	tmpDir := t.TempDir()
	originalConfigDir := configDir
	originalConfigPath := configPath
	defer func() {
		configDir = originalConfigDir
		configPath = originalConfigPath
	}()

	configDir = tmpDir
	configPath = filepath.Join(tmpDir, ConfigFile)

	tests := []struct {
		name      string
		key       string
		value     string
		expectErr bool
		verify    func(*testing.T, *Config) error
	}{
		{
			name:  "set timeout",
			key:   "timeout",
			value: "45",
			verify: func(t *testing.T, cfg *Config) error {
				if cfg.Timeout != 45 {
					return fmt.Errorf("expected Timeout 45, got %d", cfg.Timeout)
				}
				return nil
			},
		},
		{
			name:  "set with_generated_alt (hyphen not supported, use underscore)",
			key:   "with_generated_alt",
			value: "true",
			verify: func(t *testing.T, cfg *Config) error {
				if !cfg.WithGeneratedAlt {
					return fmt.Errorf("expected WithGeneratedAlt true")
				}
				return nil
			},
		},
		{
			name:  "set with_generated_alt (underscore)",
			key:   "with_generated_alt",
			value: "false",
			verify: func(t *testing.T, cfg *Config) error {
				if cfg.WithGeneratedAlt {
					return fmt.Errorf("expected WithGeneratedAlt false")
				}
				return nil
			},
		},
		{
			name:  "set proxy_url",
			key:   "proxy_url",
			value: "http://proxy.com:8080",
			verify: func(t *testing.T, cfg *Config) error {
				if cfg.ProxyURL != "http://proxy.com:8080" {
					return fmt.Errorf("expected ProxyURL 'http://proxy.com:8080', got %s", cfg.ProxyURL)
				}
				return nil
			},
		},
		{
			name:      "invalid timeout",
			key:       "timeout",
			value:     "invalid",
			expectErr: true,
		},
		{
			name:      "unknown key",
			key:       "unknown_key",
			value:     "value",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Set(tt.key, tt.value)
			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Set() failed: %v", err)
			}
			if tt.verify != nil {
				cfg, err := Load()
				if err != nil {
					t.Fatalf("Load() failed: %v", err)
				}
				if err := tt.verify(t, cfg); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	tmpDir := t.TempDir()
	originalConfigDir := configDir
	originalConfigPath := configPath
	defer func() {
		configDir = originalConfigDir
		configPath = originalConfigPath
	}()

	configDir = tmpDir
	configPath = filepath.Join(tmpDir, ConfigFile)

	// 先设置一些值
	cfg := &Config{
		Timeout:          60,
		WithGeneratedAlt: true,
		APIKey:           "test-key-12345678",
	}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	tests := []struct {
		name      string
		key       string
		expectErr bool
		expected  string
	}{
		{
			name:     "get timeout",
			key:      "timeout",
			expected: "60",
		},
		{
			name:     "get with_generated_alt",
			key:      "with_generated_alt",
			expected: "true",
		},
		{
			name:     "get api-key (hyphen)",
			key:      "api-key",
			expected: "test***5678", // 应该被掩码
		},
		{
			name:     "get api_key (underscore)",
			key:      "api_key",
			expected: "test***5678",
		},
		{
			name:      "unknown key",
			key:       "unknown_key",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := Get(tt.key)
			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("Get() failed: %v", err)
			}
			if value != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, value)
			}
		})
	}
}

func TestList(t *testing.T) {
	tmpDir := t.TempDir()
	originalConfigDir := configDir
	originalConfigPath := configPath
	defer func() {
		configDir = originalConfigDir
		configPath = originalConfigPath
	}()

	configDir = tmpDir
	configPath = filepath.Join(tmpDir, ConfigFile)

	cfg := &Config{
		ReadAPIURL:       "https://test.api.com/",
		Timeout:          45,
		WithGeneratedAlt: true,
		APIKey:           "test-key-12345678",
	}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	list, err := List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	expectedKeys := []string{
		"api_base_url",
		"search_api_url",
		"default_response_format",
		"default_output_format",
		"timeout",
		"with_generated_alt",
		"proxy_url",
		"cache_tolerance",
		"api_key",
	}

	for _, key := range expectedKeys {
		if _, ok := list[key]; !ok {
			t.Errorf("List() missing key: %s", key)
		}
	}

	if list["timeout"] != strconv.Itoa(45) {
		t.Errorf("Expected timeout '45', got %s", list["timeout"])
	}
	if list["with_generated_alt"] != "true" {
		t.Errorf("Expected with_generated_alt 'true', got %s", list["with_generated_alt"])
	}
	if list["api_key"] != "test***5678" {
		t.Errorf("Expected api_key to be masked, got %s", list["api_key"])
	}
}

func TestMaskSensitive(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"short", "***"},
		{"12345678", "***"},
		{"abcd1234efgh", "abcd***efgh"},
		{"test-api-key-12345", "test***2345"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := maskSensitive(tt.input)
			if result != tt.expected {
				t.Errorf("maskSensitive(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseKeyValue(t *testing.T) {
	tests := []struct {
		input     string
		wantKey   string
		wantValue string
		wantOK    bool
	}{
		{"key=value", "key", "value", true},
		{"key=value=with=equals", "key", "value=with=equals", true},
		{"key=", "key", "", true},
		{"=value", "", "value", true},
		{"noequals", "", "", false},
		{"  key  =  value  ", "key", "value", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			key, value, ok := parseKeyValue(tt.input)
			if ok != tt.wantOK {
				t.Errorf("parseKeyValue(%q) ok = %v, want %v", tt.input, ok, tt.wantOK)
				return
			}
			if key != tt.wantKey {
				t.Errorf("parseKeyValue(%q) key = %q, want %q", tt.input, key, tt.wantKey)
			}
			if value != tt.wantValue {
				t.Errorf("parseKeyValue(%q) value = %q, want %q", tt.input, value, tt.wantValue)
			}
		})
	}
}
