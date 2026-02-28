package output

import (
	"bytes"
	"os"
	"testing"
)

func TestNewJSONOutput(t *testing.T) {
	json := NewJSONOutput(true)
	if json == nil {
		t.Fatal("NewJSONOutput() returned nil")
	}
	if !json.pretty {
		t.Error("Expected pretty=true")
	}
}

func TestJSONOutput_Print(t *testing.T) {
	var buf bytes.Buffer
	// 重定向 stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	json := NewJSONOutput(true)
	data := map[string]interface{}{
		"url":     "https://example.com",
		"content": "test content",
	}

	err := json.Print(data)
	if err != nil {
		t.Fatalf("Print() failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	// 读取输出
	buf.ReadFrom(r)
	output := buf.String()

	// 验证 JSON 格式
	if !contains(output, `"success": true`) {
		t.Error("Expected success=true in output")
	}
	if !contains(output, `"url": "https://example.com"`) {
		t.Error("Expected url in output")
	}
}

func TestJSONOutput_Error(t *testing.T) {
	// 注意：Error() 会调用 os.Exit(1)，所以无法直接测试
	// 这里只验证编译通过
	json := NewJSONOutput(true)
	if json == nil {
		t.Fatal("NewJSONOutput() returned nil")
	}
}

func TestNewMarkdownOutput(t *testing.T) {
	// 无输出文件
	m, err := NewMarkdownOutput("")
	if err != nil {
		t.Fatalf("NewMarkdownOutput() failed: %v", err)
	}
	if m == nil {
		t.Fatal("NewMarkdownOutput() returned nil")
	}
	if m.outputFile != nil {
		t.Error("Expected outputFile to be nil")
	}

	// 临时文件
	tmpFile := os.TempDir() + "/test-output.md"
	m2, err := NewMarkdownOutput(tmpFile)
	if err != nil {
		t.Fatalf("NewMarkdownOutput() with file failed: %v", err)
	}
	if m2 == nil {
		t.Fatal("NewMarkdownOutput() returned nil")
	}
	if m2.outputFile == nil {
		t.Error("Expected outputFile to be set")
	}
	m2.Close()
	os.Remove(tmpFile)
}

func TestMarkdownOutput_Print_Map(t *testing.T) {
	var buf bytes.Buffer
	// 重定向 stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	m, _ := NewMarkdownOutput("")
	data := map[string]interface{}{
		"title":   "Test Title",
		"url":     "https://example.com",
		"content": "This is test content.",
	}

	err := m.Print(data)
	if err != nil {
		t.Fatalf("Print() failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	// 读取输出
	buf.ReadFrom(r)
	output := buf.String()

	// 验证 Markdown 格式
	if !contains(output, "# Test Title") {
		t.Error("Expected title in Markdown format")
	}
	if !contains(output, "**Source**") {
		t.Error("Expected source label")
	}
	if !contains(output, "This is test content.") {
		t.Error("Expected content in output")
	}
}

func TestMarkdownOutput_Print_Slice(t *testing.T) {
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	m, _ := NewMarkdownOutput("")
	data := []interface{}{
		map[string]interface{}{
			"title":   "Result 1",
			"url":     "https://example1.com",
			"content": "Content 1 with some text that is longer than 200 characters to test the truncation feature",
		},
		map[string]interface{}{
			"title":   "Result 2",
			"url":     "https://example2.com",
			"content": "Short content",
		},
	}

	err := m.Print(data)
	if err != nil {
		t.Fatalf("Print() failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	// 读取输出
	buf.ReadFrom(r)
	output := buf.String()

	// 验证 Markdown 格式
	if !contains(output, "## 1. Result 1") {
		t.Error("Expected first result header")
	}
	if !contains(output, "## 2. Result 2") {
		t.Error("Expected second result header")
	}
	if !contains(output, "**URL**:") {
		t.Error("Expected URL label")
	}
}

func TestGetOutput(t *testing.T) {
	tests := []struct {
		name       string
		format     OutputFormat
		outputFile string
		wantNil    bool
	}{
		{
			name:    "JSON format",
			format:  FormatJSON,
			wantNil: false,
		},
		{
			name:    "Markdown format",
			format:  FormatMarkdown,
			wantNil: false,
		},
		{
			name:    "Unknown format defaults to JSON",
			format:  "unknown",
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := GetOutput(tt.format, tt.outputFile)
			if err != nil {
				t.Fatalf("GetOutput() failed: %v", err)
			}
			if (out == nil) != tt.wantNil {
				t.Errorf("GetOutput() = %v, wantNil %v", out, tt.wantNil)
			}
		})
	}
}

func TestSuccess(t *testing.T) {
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Success(map[string]string{"key": "value"})

	w.Close()
	os.Stdout = oldStdout

	buf.ReadFrom(r)
	output := buf.String()

	if !contains(output, `"success": true`) {
		t.Error("Expected success=true")
	}
	if !contains(output, `"key": "value"`) {
		t.Error("Expected key-value pair")
	}
}

func TestPrintSuccess(t *testing.T) {
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintSuccess("test message: %s", "value")

	w.Close()
	os.Stdout = oldStdout

	buf.ReadFrom(r)
	output := buf.String()

	if !contains(output, "test message: value") {
		t.Errorf("Expected 'test message: value', got %s", output)
	}
}

func TestPrintError(t *testing.T) {
	var buf bytes.Buffer
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	PrintError("error: %s", "test error")

	w.Close()
	os.Stderr = oldStderr

	buf.ReadFrom(r)
	output := buf.String()

	if !contains(output, "error: test error") {
		t.Errorf("Expected 'error: test error', got %s", output)
	}
}

// 辅助函数
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
