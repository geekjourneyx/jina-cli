// Package api 提供 Jina Reader API 服务的 HTTP 客户端。
//
// 客户端支持两种主要操作：
//   - Read: 从 URL 读取内容并转换为 LLM 友好格式
//   - Search: 搜索网络并返回 LLM 友好的结果
//
// 使用方法：
//
//	client := api.NewClient(readURL, searchURL, apiKey, timeout)
//	resp, err := client.Read(req)
package api

// ReadRequest Read 请求
type ReadRequest struct {
	URL              string
	Method           string // GET or POST
	ResponseFormat   string // markdown, html, text, screenshot
	Headers          map[string]string
	Timeout          int // 超时时间（秒）
	NoCache          bool
	ProxyURL         string
	TargetSelector   string
	WaitForSelector  string
	Cookie           string
	WithGeneratedAlt bool
	PostMethod       bool // 使用 POST 方法（用于 SPA）
}

// ReadResponse Read 响应
type ReadResponse struct {
	Content string
	URL     string
	Title   string
}

// SearchRequest Search 请求
type SearchRequest struct {
	Query          string
	Sites          []string
	ResponseFormat string
	Headers        map[string]string
	Timeout        int
	Limit          int
}

// SearchResult 单个搜索结果
type SearchResult struct {
	Title   string
	URL     string
	Content string
}

// SearchResponse Search 响应
type SearchResponse struct {
	Query   string
	Results []SearchResult
}
