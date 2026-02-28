package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/geekjourneyx/jina-cli/cli/pkg/config"
	"github.com/geekjourneyx/jina-cli/cli/pkg/output"
)

// 版本信息（构建时注入）
var (
	version   = "1.0.0"
	buildDate = "unknown"
	gitCommit = "unknown"
)

func main() {
	// 执行根命令
	if err := rootCmd.Execute(); err != nil {
		output.Error(err)
		os.Exit(1)
	}
}

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "jina",
	Short: "CLI tool for Jina AI Reader and Search APIs",
	Long: `jina is a lightweight CLI tool that wraps the Jina AI Reader API.
It converts any URL to LLM-friendly input and searches the web with AI.

Quick Start:
  jina read --url "https://example.com"
  jina search --query "golang latest news"

Get Help:
  jina --help
  jina [command] --help`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

var cfg *config.Config

// 执行命令前的初始化
func initConfig(cmd *cobra.Command, args []string) error {
	var err error
	cfg, err = config.Load()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	return nil
}

func init() {
	// 添加版本标志
	rootCmd.Version = fmt.Sprintf("%s (构建时间: %s, 提交: %s)", version, buildDate, gitCommit)

	// 添加子命令
	rootCmd.AddCommand(ReadCmd)
	rootCmd.AddCommand(SearchCmd)
	rootCmd.AddCommand(ConfigCmd)

	// 持久化标志
	rootCmd.PersistentFlags().StringP("api-base", "a", "", "API base URL (overrides config)")
	rootCmd.PersistentFlags().StringP("api-key", "k", "", "API key (overrides config)")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Output format: json, markdown (default: json)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")

	// 绑定持久化标志到配置
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// 某些命令不需要配置（如 help, version, config set）
		skipConfig := cmd.Name() == "help" || cmd.Name() == "config" || cmd.Name() == "completion"
		if !skipConfig {
			if err := initConfig(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}
