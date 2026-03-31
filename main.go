package main

import (
	"context"
	"io"
	"llm-cli/tools"
	"os"
	"path/filepath"
)

func main() {

	// 加载配置
	// 1. 获取系统用户主目录（跨平台支持：Windows/macOS/Linux）
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("获取用户主目录失败：" + err.Error())
	}

	// 2. 拼接完整配置文件路径（推荐用 filepath.Join，自动适配系统分隔符）
	configPath := filepath.Join(homeDir, ".config", "llm-cli", "config.yaml")

	cfg := loadConfig(configPath)

	// 获取输入
	input, err := readInput()
	if err != nil {
		panic("获取输入失败: " + err.Error())
	}

	apiKey := cfg.Providers[cfg.DefaultProvider].APIKey
	baseUrl := cfg.Providers[cfg.DefaultProvider].BaseURL
	model := cfg.Providers[cfg.DefaultProvider].Model
	systemPrompt := cfg.Prompts[cfg.DefaultPrompt]

	bot := NewChatBot(apiKey, baseUrl, model, systemPrompt)

	bot.RegisterTool(&tools.MultiplyTool{})
	bot.RegisterTool(&tools.SearchTool{APIKey: cfg.Tools.Tavily.APIKey})
	bot.RegisterTool(&tools.ShellTool{})

	ctx := context.Background()

	bot.Chat(ctx, input)
}

// readInput 从管道和命令行参数中读取并组合用户的输入内容
func readInput() (string, error) {
	var pipeInput string
	var argInput string

	// 1. 检测并读取管道输入 (stdin)
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	// 按位与操作：如果不是字符设备（即不是直接在终端运行），说明有管道输入或文件重定向
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// 从标准输入读取所有数据
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		pipeInput = string(bytes)
	}

	// 2. 读取命令行参数
	if len(os.Args) > 1 {
		argInput = os.Args[1]
	}

	// 3. 组合输入内容
	finalInput := pipeInput + argInput

	return finalInput, nil
}
