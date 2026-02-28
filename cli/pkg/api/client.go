package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client API 客户端
type Client struct {
	readAPIURL   string
	searchAPIURL string
	apiKey       string
	httpClient   *http.Client
}

// NewClient 创建 API 客户端
func NewClient(readURL, searchURL, apiKey string, timeout int) *Client {
	return &Client{
		readAPIURL:   readURL,
		searchAPIURL: searchURL,
		apiKey:       apiKey,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// SetTimeout 设置请求超时
func (c *Client) SetTimeout(timeout int) {
	c.httpClient.Timeout = time.Duration(timeout) * time.Second
}

// Read 执行 Read API 请求
func (c *Client) Read(req *ReadRequest) (*ReadResponse, error) {
	// 构建请求 URL
	var fullURL string
	var body io.Reader

	if req.PostMethod {
		// POST 方法用于 SPA 带 hash 路由的情况
		fullURL = c.readAPIURL
		formData := url.Values{}
		formData.Set("url", req.URL)
		body = strings.NewReader(formData.Encode())
	} else {
		// GET 方法
		fullURL = c.readAPIURL + "/" + url.PathEscape(req.URL)
	}

	// 创建 HTTP 请求
	var err error
	var httpReq *http.Request

	if req.PostMethod {
		httpReq, err = http.NewRequest("POST", fullURL, body)
		if err != nil {
			return nil, fmt.Errorf("创建请求失败: %w", err)
		}
		httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		httpReq, err = http.NewRequest("GET", fullURL, nil)
		if err != nil {
			return nil, fmt.Errorf("创建请求失败: %w", err)
		}
	}

	// 设置请求头
	c.setCommonHeaders(httpReq)
	c.setRequestHeaders(httpReq, req)

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP 错误: %d, 响应: %s", resp.StatusCode, string(bodyBytes))
	}

	// 读取响应
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 构建响应
	return &ReadResponse{
		Content: string(content),
		URL:     req.URL,
	}, nil
}

// Search 执行 Search API 请求
func (c *Client) Search(req *SearchRequest) (*SearchResponse, error) {
	// 构建查询 URL
	queryParams := url.Values{}
	if len(req.Sites) > 0 {
		for _, site := range req.Sites {
			queryParams.Add("site", site)
		}
	}

	// 编码查询字符串
	queryEncoded := url.QueryEscape(req.Query)
	fullURL := c.searchAPIURL + "/" + queryEncoded
	if queryParams.Encode() != "" {
		fullURL += "?" + queryParams.Encode()
	}

	// 创建 HTTP 请求
	httpReq, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	c.setCommonHeaders(httpReq)
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP 错误: %d, 响应: %s", resp.StatusCode, string(bodyBytes))
	}

	// 读取响应
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析搜索结果（从 JSON 格式或纯文本格式）
	// Search API 返回的是 JSON 格式的结果数组
	// 如果没有指定 Accept: application/json，返回纯文本
	results := parseSearchResults(string(content), req.ResponseFormat)

	return &SearchResponse{
		Query:   req.Query,
		Results: results,
	}, nil
}

// setCommonHeaders 设置通用请求头
func (c *Client) setCommonHeaders(req *http.Request) {
	// 设置 User-Agent
	req.Header.Set("User-Agent", "jina-cli/1.0.0")

	// 如果有 API Key，设置 Authorization 头
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
}

// setRequestHeaders 设置 Read 请求的特定头
func (c *Client) setRequestHeaders(req *http.Request, readReq *ReadRequest) {
	// 设置响应格式头
	if readReq.ResponseFormat != "" {
		req.Header.Set("X-Respond-With", readReq.ResponseFormat)
	}

	// 启用图片描述
	if readReq.WithGeneratedAlt {
		req.Header.Set("X-With-Generated-Alt", "true")
	}

	// 禁用缓存
	if readReq.NoCache {
		req.Header.Set("X-No-Cache", "true")
	}

	// 设置代理
	if readReq.ProxyURL != "" {
		req.Header.Set("X-Proxy-URL", readReq.ProxyURL)
	}

	// 设置目标选择器
	if readReq.TargetSelector != "" {
		req.Header.Set("X-Target-Selector", readReq.TargetSelector)
	}

	// 设置等待选择器
	if readReq.WaitForSelector != "" {
		req.Header.Set("X-Wait-For-Selector", readReq.WaitForSelector)
	}

	// 设置 Cookie
	if readReq.Cookie != "" {
		req.Header.Set("X-Set-Cookie", readReq.Cookie)
	}

	// 自定义请求头
	for k, v := range readReq.Headers {
		req.Header.Set(k, v)
	}
}

// parseSearchResults 解析搜索结果
func parseSearchResults(content string, format string) []SearchResult {
	// 如果是 JSON 格式，需要解析 JSON
	// 如果是纯文本格式，按行分割
	results := []SearchResult{}

	// 简单实现：按行分割，每行是一个结果
	// 实际上 Search API 返回的是 JSON 格式
	// 这里先做一个简单的实现，后续可以根据需要改进

	// 尝试检测是否是 JSON 格式
	if strings.HasPrefix(strings.TrimSpace(content), "[") {
		// TODO: 实现 JSON 解析
		// 目前先返回空结果
		return results
	}

	// 纯文本格式，每行是一个结果
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			results = append(results, SearchResult{
				Content: line,
			})
		}
	}

	return results
}
