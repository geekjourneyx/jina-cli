package main

import (
	"fmt"

	"github.com/geekjourneyx/jina-cli/cli/pkg/api"
	"github.com/geekjourneyx/jina-cli/cli/pkg/output"
	"github.com/spf13/cobra"
)

// SearchCmd search 命令
var SearchCmd = &cobra.Command{
	Use:     "search",
	Aliases: []string{"s"},
	Short:   "Search the web with AI-powered results",
	Long:    `Search the web and return results in LLM-friendly format. Automatically fetches content from top 5 results.`,
	Example: `  jina search --query "golang latest news"
  jina search -q "AI developments" --site techcrunch.com --site theverge.com
  jina search -q "climate change" --limit 10 --output markdown`,
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateSearchFlags()
	},
	Run: runSearch,
}

var (
	flagSearchQuery      string
	flagSearchSites      []string
	flagSearchFormat     string
	flagSearchTimeout    int
	flagSearchLimit      int
	flagSearchOutputFile string
)

func init() {
	SearchCmd.Flags().StringVarP(&flagSearchQuery, "query", "q", "", "Search query (required)")
	SearchCmd.Flags().StringSliceVarP(&flagSearchSites, "site", "s", []string{}, "Restrict to specific domains (repeatable)")
	SearchCmd.Flags().StringVarP(&flagSearchFormat, "format", "F", "", "Response format: markdown, html, text (default: markdown)")
	SearchCmd.Flags().IntVarP(&flagSearchTimeout, "timeout", "t", 0, "Request timeout in seconds")
	SearchCmd.Flags().IntVarP(&flagSearchLimit, "limit", "l", 0, "Max results to return (default: 5)")
	SearchCmd.Flags().StringVarP(&flagSearchOutputFile, "output-file", "O", "", "Write output to file instead of stdout")
}

func validateSearchFlags() error {
	if flagSearchQuery == "" {
		return fmt.Errorf("必须提供 --query 参数")
	}
	return nil
}

func runSearch(cmd *cobra.Command, args []string) {
	// 获取输出格式
	outputFormat := getSearchOutputFormat(cmd)

	// 获取 API Base URL（命令行参数优先）
	searchAPIURL := cfg.SearchAPIURL
	if apiBaseFlag, _ := cmd.Parent().PersistentFlags().GetString("api-base"); apiBaseFlag != "" {
		// 如果用户指定了 api-base，使用它作为 search API URL
		searchAPIURL = apiBaseFlag
	}

	// 获取 API Key
	apiKey := cfg.APIKey
	if apiKeyFlag, _ := cmd.Parent().PersistentFlags().GetString("api-key"); apiKeyFlag != "" {
		apiKey = apiKeyFlag
	}

	// 获取超时时间
	timeout := cfg.Timeout
	if flagSearchTimeout > 0 {
		timeout = flagSearchTimeout
	}

	// 获取响应格式
	responseFormat := cfg.DefaultResponseFormat
	if flagSearchFormat != "" {
		responseFormat = flagSearchFormat
	}

	// 获取结果限制
	limit := 5
	if flagSearchLimit > 0 {
		limit = flagSearchLimit
	}

	// 创建 API 客户端
	client := api.NewClient(cfg.ReadAPIURL, searchAPIURL, apiKey, timeout)

	// 获取输出处理器
	out, err := output.GetOutput(output.OutputFormat(outputFormat), flagSearchOutputFile)
	if err != nil {
		output.Error(err)
	}
	defer func() {
		if closer, ok := out.(interface{ Close() error }); ok {
			closer.Close()
		}
	}()

	// 构建请求
	req := &api.SearchRequest{
		Query:          flagSearchQuery,
		Sites:          flagSearchSites,
		ResponseFormat: responseFormat,
		Timeout:        timeout,
		Limit:          limit,
	}

	// 执行搜索
	resp, err := client.Search(req)
	if err != nil {
		_ = out.Error(err)
		return
	}

	// 构建输出数据
	results := make([]map[string]interface{}, 0, len(resp.Results))
	for _, result := range resp.Results {
		r := map[string]interface{}{
			"content": result.Content,
		}
		if result.Title != "" {
			r["title"] = result.Title
		}
		if result.URL != "" {
			r["url"] = result.URL
		}
		results = append(results, r)
	}

	outputData := map[string]interface{}{
		"query":   resp.Query,
		"results": results,
		"count":   len(results),
	}

	out.Print(outputData)
}

func getSearchOutputFormat(cmd *cobra.Command) string {
	// 持久化标志
	outputFlag, _ := cmd.Parent().PersistentFlags().GetString("output")
	if outputFlag != "" {
		return outputFlag
	}

	// 配置文件
	return cfg.DefaultOutputFormat
}
