package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/geekjourneyx/jina-cli/cli/pkg/api"
	"github.com/geekjourneyx/jina-cli/cli/pkg/output"
)

// ReadCmd read 命令
var ReadCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	Short:   "Extract and convert content from URLs",
	Long:    `Read and extract content from any URL, converting it to LLM-friendly format (Markdown, HTML, or Text).`,
	Example: `  jina read --url "https://example.com"
  jina read -u "https://x.com/user/status/123" --with-alt
  jina read --file urls.txt --output markdown`,
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateReadFlags()
	},
	Run: runRead,
}

var (
	flagReadURL             string
	flagReadFile            string
	flagReadFormat          string
	flagReadTimeout         int
	flagReadWithAlt         bool
	flagReadNoCache         bool
	flagReadProxy           string
	flagReadTargetSelector  string
	flagReadWaitForSelector string
	flagReadCookie          string
	flagReadPostMethod      bool
	flagReadOutputFile      string
)

func init() {
	ReadCmd.Flags().StringVarP(&flagReadURL, "url", "u", "", "URL to read (required if --file not used)")
	ReadCmd.Flags().StringVarP(&flagReadFile, "file", "f", "", "File containing URLs (one per line)")
	ReadCmd.Flags().StringVarP(&flagReadFormat, "format", "F", "", "Response format: markdown, html, text, screenshot (default: markdown)")
	ReadCmd.Flags().IntVarP(&flagReadTimeout, "timeout", "t", 0, "Request timeout in seconds")
	ReadCmd.Flags().BoolVar(&flagReadWithAlt, "with-alt", false, "Enable image captioning with VLM")
	ReadCmd.Flags().BoolVar(&flagReadNoCache, "no-cache", false, "Bypass cache")
	ReadCmd.Flags().StringVar(&flagReadProxy, "proxy", "", "Proxy server URL")
	ReadCmd.Flags().StringVar(&flagReadTargetSelector, "target-selector", "", "CSS selector for content extraction")
	ReadCmd.Flags().StringVar(&flagReadWaitForSelector, "wait-for-selector", "", "CSS selector to wait for")
	ReadCmd.Flags().StringVar(&flagReadCookie, "cookie", "", "Cookie string to forward")
	ReadCmd.Flags().BoolVar(&flagReadPostMethod, "post", false, "Use POST method (for SPA with hash routing)")
	ReadCmd.Flags().StringVarP(&flagReadOutputFile, "output-file", "O", "", "Write output to file instead of stdout")
}

func validateReadFlags() error {
	// 检查 URL 来源
	if flagReadURL == "" && flagReadFile == "" {
		return fmt.Errorf("必须提供 --url 或 --file 参数")
	}
	if flagReadURL != "" && flagReadFile != "" {
		return fmt.Errorf("--url 和 --file 不能同时使用")
	}

	return nil
}

func runRead(cmd *cobra.Command, args []string) {
	// 获取输出格式
	outputFormat := getReadOutputFormat(cmd)

	// 获取 API Base URL（命令行参数优先）
	apiBase := cfg.ReadAPIURL
	if apiBaseFlag, _ := cmd.Parent().PersistentFlags().GetString("api-base"); apiBaseFlag != "" {
		apiBase = apiBaseFlag
	}

	// 获取 API Key（命令行参数优先）
	apiKey := cfg.APIKey
	if apiKeyFlag, _ := cmd.Parent().PersistentFlags().GetString("api-key"); apiKeyFlag != "" {
		apiKey = apiKeyFlag
	}

	// 获取超时时间
	timeout := cfg.Timeout
	if flagReadTimeout > 0 {
		timeout = flagReadTimeout
	}

	// 获取响应格式
	responseFormat := cfg.DefaultResponseFormat
	if flagReadFormat != "" {
		responseFormat = flagReadFormat
	}

	// 创建 API 客户端
	client := api.NewClient(apiBase, cfg.SearchAPIURL, apiKey, timeout)

	// 获取输出处理器
	out, err := output.GetOutput(output.OutputFormat(outputFormat), flagReadOutputFile)
	if err != nil {
		output.Error(err)
	}
	defer func() {
		if closer, ok := out.(interface{ Close() error }); ok {
			closer.Close()
		}
	}()

	// 处理 URL
	if flagReadURL != "" {
		// 单个 URL
		processURL(client, flagReadURL, responseFormat, out)
	} else {
		// 批量处理
		processBatch(client, flagReadFile, responseFormat, out)
	}
}

func processURL(client *api.Client, url, responseFormat string, out output.Output) {
	req := &api.ReadRequest{
		URL:              url,
		Method:           "GET",
		ResponseFormat:   responseFormat,
		WithGeneratedAlt: flagReadWithAlt,
		NoCache:          flagReadNoCache,
		ProxyURL:         flagReadProxy,
		TargetSelector:   flagReadTargetSelector,
		WaitForSelector:  flagReadWaitForSelector,
		Cookie:           flagReadCookie,
		PostMethod:       flagReadPostMethod,
	}

	resp, err := client.Read(req)
	if err != nil {
		out.Error(err)
		return
	}

	// 构建输出数据
	result := map[string]interface{}{
		"url":     resp.URL,
		"content": resp.Content,
	}

	// 尝试从内容中提取标题（对于 Markdown 格式）
	if responseFormat == "markdown" || responseFormat == "" {
		if title := extractTitle(resp.Content); title != "" {
			result["title"] = title
		}
	}

	out.Print(result)
}

func processBatch(client *api.Client, filename, responseFormat string, out output.Output) {
	// 读取文件
	content, err := os.ReadFile(filename)
	if err != nil {
		out.Error(fmt.Errorf("读取文件失败: %w", err))
		return
	}

	// 解析 URL 列表
	urls := parseURLList(string(content))
	if len(urls) == 0 {
		out.Error(fmt.Errorf("文件中没有找到有效的 URL"))
		return
	}

	// 获取输出格式用于显示进度
	outputFormat := getReadOutputFormat(nil)
	if outputFormat == "markdown" {
		fmt.Fprintf(os.Stderr, "正在处理 %d 个 URL...\n", len(urls))
	}

	// 处理每个 URL
	results := make([]map[string]interface{}, 0, len(urls))
	for i, url := range urls {
		if outputFormat == "markdown" {
			fmt.Fprintf(os.Stderr, "正在处理 [%d/%d]: %s\n", i+1, len(urls), url)
		}

		req := &api.ReadRequest{
			URL:              url,
			Method:           "GET",
			ResponseFormat:   responseFormat,
			WithGeneratedAlt: flagReadWithAlt,
			NoCache:          flagReadNoCache,
			ProxyURL:         flagReadProxy,
		}

		resp, err := client.Read(req)
		if err != nil {
			result := map[string]interface{}{
				"url":   url,
				"error": err.Error(),
			}
			results = append(results, result)
			continue
		}

		result := map[string]interface{}{
			"url":     resp.URL,
			"content": resp.Content,
		}

		// 尝试提取标题
		if title := extractTitle(resp.Content); title != "" {
			result["title"] = title
		}

		results = append(results, result)
	}

	// 输出结果
	out.Print(results)
}

func parseURLList(content string) []string {
	lines := strings.Split(content, "\n")
	urls := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 跳过空行和注释行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		urls = append(urls, line)
	}

	return urls
}

func extractTitle(content string) string {
	// 尝试从 Markdown 内容中提取标题
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return ""
}

func getReadOutputFormat(cmd *cobra.Command) string {
	// 持久化标志
	var outputFlag string
	if cmd != nil {
		outputFlag, _ = cmd.Parent().PersistentFlags().GetString("output")
	}

	if outputFlag != "" {
		return outputFlag
	}

	// 配置文件
	return cfg.DefaultOutputFormat
}
