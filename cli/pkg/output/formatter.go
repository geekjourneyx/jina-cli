// Package output 提供统一的输出格式化功能。
//
// 所有输出为 JSON 或 Markdown 格式，包含 success 字段表示操作是否成功。
// 成功时包含 data 字段，失败时包含 error 字段。
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// OutputFormat 输出格式类型
type OutputFormat string

const (
	// FormatJSON JSON 格式
	FormatJSON OutputFormat = "json"
	// FormatMarkdown Markdown 格式
	FormatMarkdown OutputFormat = "markdown"
)

// SuccessResponse 成功响应
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Code    string `json:"code,omitempty"`
}

// Output 输出接口
type Output interface {
	Print(data interface{}) error
	Error(err error) error
}

// JSONOutput JSON 输出
type JSONOutput struct {
	pretty bool
}

// NewJSONOutput 创建 JSON 输出
func NewJSONOutput(pretty bool) *JSONOutput {
	return &JSONOutput{pretty: pretty}
}

// Print 输出数据
func (j *JSONOutput) Print(data interface{}) error {
	resp := SuccessResponse{
		Success: true,
		Data:    data,
	}
	return j.printJSON(resp)
}

// Error 输出错误
func (j *JSONOutput) Error(err error) error {
	resp := ErrorResponse{
		Success: false,
		Error:   err.Error(),
	}
	if err := j.printJSON(resp); err != nil {
		return err
	}
	os.Exit(1)
	return nil
}

func (j *JSONOutput) printJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	if j.pretty {
		encoder.SetIndent("", "  ")
	}
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("JSON 编码错误: %w", err)
	}
	return nil
}

// MarkdownOutput Markdown 输出
type MarkdownOutput struct {
	outputFile *os.File
}

// NewMarkdownOutput 创建 Markdown 输出
func NewMarkdownOutput(outputFile string) (*MarkdownOutput, error) {
	m := &MarkdownOutput{}
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			return nil, fmt.Errorf("创建输出文件失败: %w", err)
		}
		m.outputFile = f
	}
	return m, nil
}

// Print 输出数据
func (m *MarkdownOutput) Print(data interface{}) error {
	// 根据数据类型格式化
	switch d := data.(type) {
	case map[string]interface{}:
		m.printMapAsMarkdown(d)
	case []interface{}:
		m.printSliceAsMarkdown(d)
	default:
		fmt.Fprintf(m.getWriter(), "%v\n", data)
	}

	if m.outputFile != nil {
		m.outputFile.Close()
	}
	return nil
}

// Error 输出错误
func (m *MarkdownOutput) Error(err error) error {
	fmt.Fprintf(m.getWriter(), "**Error**: %s\n", err.Error())
	if m.outputFile != nil {
		m.outputFile.Close()
	}
	os.Exit(1)
	return nil
}

func (m *MarkdownOutput) getWriter() io.Writer {
	if m.outputFile != nil {
		return m.outputFile
	}
	return os.Stdout
}

func (m *MarkdownOutput) printMapAsMarkdown(data map[string]interface{}) {
	w := m.getWriter()

	// 尝试提取常见字段
	if title, ok := data["title"].(string); ok {
		fmt.Fprintf(w, "# %s\n\n", title)
	}
	if url, ok := data["url"].(string); ok {
		fmt.Fprintf(w, "**Source**: <%s>\n\n", url)
	}
	if content, ok := data["content"].(string); ok {
		fmt.Fprintf(w, "%s\n", content)
	} else {
		// 通用格式
		for k, v := range data {
			fmt.Fprintf(w, "**%s**: %v\n", k, v)
		}
	}
}

func (m *MarkdownOutput) printSliceAsMarkdown(data []interface{}) {
	w := m.getWriter()
	for i, item := range data {
		if m, ok := item.(map[string]interface{}); ok {
			if title, ok := m["title"].(string); ok {
				fmt.Fprintf(w, "## %d. %s\n", i+1, title)
				if url, ok := m["url"].(string); ok {
					fmt.Fprintf(w, "**URL**: <%s>\n", url)
				}
				if content, ok := m["content"].(string); ok {
					// 只显示前 200 字符
					if len(content) > 200 {
						content = content[:200] + "..."
					}
					fmt.Fprintf(w, "%s\n", content)
				}
				fmt.Fprintln(w)
			}
		}
	}
}

// Close 关闭输出文件
func (m *MarkdownOutput) Close() error {
	if m.outputFile != nil {
		return m.outputFile.Close()
	}
	return nil
}

// Success 输出成功响应（JSON 格式，兼容旧代码）
func Success(data interface{}) {
	resp := SuccessResponse{
		Success: true,
		Data:    data,
	}
	printJSON(resp)
}

// Error 输出错误响应（JSON 格式，兼容旧代码）
func Error(err error) {
	resp := ErrorResponse{
		Success: false,
		Error:   err.Error(),
	}
	printJSON(resp)
	os.Exit(1)
}

// ErrorWithCode 输出带错误码的错误响应（JSON 格式，兼容旧代码）
func ErrorWithCode(code, message string) {
	resp := ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	}
	printJSON(resp)
	os.Exit(1)
}

// printJSON 打印 JSON（兼容旧代码）
func printJSON(v interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "JSON 编码错误: %v\n", err)
		os.Exit(1)
	}
}

// PrintSuccess 打印成功消息（兼容函数）
func PrintSuccess(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}

// PrintError 打印错误消息到 stderr（兼容函数）
func PrintError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

// GetOutput 获取输出处理器
func GetOutput(format OutputFormat, outputFile string) (Output, error) {
	switch format {
	case FormatMarkdown:
		return NewMarkdownOutput(outputFile)
	case FormatJSON:
		return NewJSONOutput(true), nil
	default:
		return NewJSONOutput(true), nil
	}
}
