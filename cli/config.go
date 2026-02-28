package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/geekjourneyx/jina-cli/cli/pkg/config"
	"github.com/geekjourneyx/jina-cli/cli/pkg/output"
)

// ConfigCmd config 命令
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `Manage jina configuration file. Configuration is stored in ~/.jina-reader/config.yaml`,
}

func init() {
	// 添加子命令
	ConfigCmd.AddCommand(configSetCmd)
	ConfigCmd.AddCommand(configGetCmd)
	ConfigCmd.AddCommand(configListCmd)
	ConfigCmd.AddCommand(configPathCmd)
}

// configSetCmd 设置配置
var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long:  `Set a configuration value. The value will be saved to the config file.`,
	Example: `  jina config set api-base "https://r.jina.ai/"
  jina config set timeout 60
  jina config set with-generated-alt true`,
	Args: cobra.ExactArgs(2),
	Run:  runConfigSet,
}

func runConfigSet(cmd *cobra.Command, args []string) {
	key := args[0]
	value := args[1]

	if err := config.Set(key, value); err != nil {
		output.Error(err)
	}

	output.Success(map[string]interface{}{
		"key":   key,
		"value": value,
	})
}

// configGetCmd 获取配置
var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Long:  `Get a configuration value and display it.`,
	Example: `  jina config get api-base
  jina config get timeout`,
	Args: cobra.ExactArgs(1),
	Run:  runConfigGet,
}

func runConfigGet(cmd *cobra.Command, args []string) {
	key := args[0]

	value, err := config.Get(key)
	if err != nil {
		output.Error(err)
	}

	if value == "" {
		fmt.Println("(未设置)")
		return
	}

	fmt.Println(value)
}

// configListCmd 列出所有配置
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration",
	Long:  `List all configuration values.`,
	Run:   runConfigList,
}

func runConfigList(cmd *cobra.Command, args []string) {
	cfgList, err := config.List()
	if err != nil {
		output.Error(err)
	}

	// 以表格形式输出
	fmt.Println("配置列表:")
	fmt.Println("========================================")
	displayOrder := []string{
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

	for _, key := range displayOrder {
		if value, ok := cfgList[key]; ok {
			fmt.Printf("%-25s : %s\n", key, value)
		}
	}
	fmt.Println("========================================")
	fmt.Printf("配置文件路径: %s\n", config.GetConfigPath())
}

// configPathCmd 显示配置文件路径
var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show config file path",
	Long:  `Show the full path to the configuration file.`,
	Run:   runConfigPath,
}

func runConfigPath(cmd *cobra.Command, args []string) {
	path := config.GetConfigPath()

	// 检查文件是否存在
	fileExists := false
	if _, err := os.Stat(path); err == nil {
		fileExists = true
	}

	output.Success(map[string]interface{}{
		"path":   path,
		"exists": fileExists,
	})
}
