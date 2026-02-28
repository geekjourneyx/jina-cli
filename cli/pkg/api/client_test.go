package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("https://r.jina.ai/", "https://s.jina.ai/", "", 30)
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.readAPIURL != "https://r.jina.ai/" {
		t.Errorf("Expected readAPIURL 'https://r.jina.ai/', got %s", client.readAPIURL)
	}
	if client.searchAPIURL != "https://s.jina.ai/" {
		t.Errorf("Expected searchAPIURL 'https://s.jina.ai/', got %s", client.searchAPIURL)
	}
	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", client.httpClient.Timeout)
	}
}

func TestClient_SetTimeout(t *testing.T) {
	client := NewClient("https://r.jina.ai/", "https://s.jina.ai/", "", 30)
	client.SetTimeout(60)
	if client.httpClient.Timeout != 60*time.Second {
		t.Errorf("Expected timeout 60s, got %v", client.httpClient.Timeout)
	}
}

func TestClient_Read_Success(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// 验证请求路径
		if !strings.Contains(r.URL.Path, "https://example.com") {
			t.Errorf("Expected path to contain https://example.com, got %s", r.URL.Path)
		}

		// 验证请求头
		if ua := r.Header.Get("User-Agent"); !strings.Contains(ua, "jina-cli") {
			t.Errorf("Expected User-Agent to contain jina-cli, got %s", ua)
		}

		// 返回成功响应
		w.Header().Set("Content-Type", "text/markdown")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("# Test Content\n\nThis is a test."))
	}))
	defer server.Close()

	// 创建客户端
	client := NewClient(server.URL+"/", server.URL+"/", "", 30)

	// 执行请求
	req := &ReadRequest{
		URL: "https://example.com",
	}
	resp, err := client.Read(req)
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}

	// 验证响应
	if resp.Content != "# Test Content\n\nThis is a test." {
		t.Errorf("Expected content '# Test Content\\n\\nThis is a test.', got %q", resp.Content)
	}
	if resp.URL != "https://example.com" {
		t.Errorf("Expected URL 'https://example.com', got %s", resp.URL)
	}
}

func TestClient_Read_WithHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证自定义头
		if r.Header.Get("X-Respond-With") != "markdown" {
			t.Errorf("Expected X-Respond-With 'markdown', got %s", r.Header.Get("X-Respond-With"))
		}
		if r.Header.Get("X-With-Generated-Alt") != "true" {
			t.Errorf("Expected X-With-Generated-Alt 'true', got %s", r.Header.Get("X-With-Generated-Alt"))
		}
		if r.Header.Get("X-No-Cache") != "true" {
			t.Errorf("Expected X-No-Cache 'true', got %s", r.Header.Get("X-No-Cache"))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("content"))
	}))
	defer server.Close()

	client := NewClient(server.URL+"/", server.URL+"/", "", 30)

	req := &ReadRequest{
		URL:              "https://example.com",
		ResponseFormat:   "markdown",
		WithGeneratedAlt: true,
		NoCache:          true,
	}

	_, err := client.Read(req)
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}
}

func TestClient_Read_POSTMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证 POST 方法
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// 验证 Content-Type
		if ct := r.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
			t.Errorf("Expected Content-Type 'application/x-www-form-urlencoded', got %s", ct)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("content"))
	}))
	defer server.Close()

	client := NewClient(server.URL, server.URL, "", 30)

	req := &ReadRequest{
		URL:        "https://example.com/#/route",
		PostMethod: true,
	}

	_, err := client.Read(req)
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}
}

func TestClient_Read_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	client := NewClient(server.URL+"/", server.URL+"/", "", 30)

	req := &ReadRequest{
		URL: "https://example.com",
	}

	_, err := client.Read(req)
	if err == nil {
		t.Fatal("Expected error for 404 response, got nil")
	}
	if !strings.Contains(err.Error(), "HTTP 错误: 404") {
		t.Errorf("Expected error to contain 'HTTP 错误: 404', got %v", err)
	}
}

func TestClient_Search_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Result 1\nResult 2\nResult 3"))
	}))
	defer server.Close()

	client := NewClient("https://r.jina.ai/", server.URL+"/", "", 30)

	req := &SearchRequest{
		Query: "test query",
	}

	resp, err := client.Search(req)
	if err != nil {
		t.Fatalf("Search() failed: %v", err)
	}

	// 验证响应
	if len(resp.Results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(resp.Results))
	}
	if resp.Query != "test query" {
		t.Errorf("Expected query 'test query', got %s", resp.Query)
	}
}

func TestClient_Search_WithSiteFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证 site 参数
		if sites := r.URL.Query()["site"]; len(sites) != 2 {
			t.Errorf("Expected 2 site parameters, got %d", len(sites))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("filtered results"))
	}))
	defer server.Close()

	client := NewClient("https://r.jina.ai/", server.URL+"/", "", 30)

	req := &SearchRequest{
		Query: "test query",
		Sites: []string{"example.com", "test.com"},
	}

	_, err := client.Search(req)
	if err != nil {
		t.Fatalf("Search() failed: %v", err)
	}
}

func TestClient_Search_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewClient("https://r.jina.ai/", server.URL+"/", "", 30)

	req := &SearchRequest{
		Query: "test query",
	}

	_, err := client.Search(req)
	if err == nil {
		t.Fatal("Expected error for 500 response, got nil")
	}
	if !strings.Contains(err.Error(), "HTTP 错误: 500") {
		t.Errorf("Expected error to contain 'HTTP 错误: 500', got %v", err)
	}
}

func TestParseSearchResults(t *testing.T) {
	tests := []struct {
		name    string
		content string
		format  string
		wantLen int
	}{
		{
			name:    "empty content",
			content: "",
			format:  "",
			wantLen: 0,
		},
		{
			name:    "single line",
			content: "Single result",
			format:  "",
			wantLen: 1,
		},
		{
			name:    "multiple lines",
			content: "Result 1\nResult 2\nResult 3",
			format:  "",
			wantLen: 3,
		},
		{
			name:    "lines with empty lines",
			content: "Result 1\n\nResult 2\n  \nResult 3",
			format:  "",
			wantLen: 3,
		},
		{
			name:    "JSON format (not implemented yet)",
			content: "[]",
			format:  "json",
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := parseSearchResults(tt.content, tt.format)
			if len(results) != tt.wantLen {
				t.Errorf("parseSearchResults() len = %d, want %d", len(results), tt.wantLen)
			}
		})
	}
}
