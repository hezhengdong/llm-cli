package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"llm-cli/tools"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	input, reasoningEnabled, err := readInput()
	if err != nil {
		panic("获取输入失败: " + err.Error())
	}

	apiKey := cfg.Providers[cfg.DefaultProvider].APIKey
	baseUrl := cfg.Providers[cfg.DefaultProvider].BaseURL
	model := cfg.Providers[cfg.DefaultProvider].Model
	systemPrompt := cfg.Prompts[cfg.DefaultPrompt]

	bot := NewChatBot(apiKey, baseUrl, model, systemPrompt)

	bot.reasoningEnabled = reasoningEnabled

	bot.RegisterTool(&tools.MultiplyTool{})
	bot.RegisterTool(&tools.SearchTool{APIKey: cfg.Tools.Tavily.APIKey})
	bot.RegisterTool(&tools.ShellTool{})

	ctx := context.Background()

	bot.Chat(ctx, input)
}

// readInput 从管道和命令行参数中读取并组合用户的输入内容
func readInput() (string, bool, error) {
	var pipeInput string
	var argInput string


	// 1. 定义开关：默认 false（关闭思考）
	reasoning := flag.Bool("r", false, "enable reasoning mode")
	flag.Parse()

	// 2. 检测并读取管道输入 (stdin)
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", false, err
	}

	// 按位与操作：如果不是字符设备（即不是直接在终端运行），说明有管道输入或文件重定向
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// 从标准输入读取所有数据
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", false, err
		}
		pipeInput = string(bytes)
	}

	// 3. 读取命令行参数
	argInput = strings.Join(flag.Args(), " ")

	// 4. 生成当前时间提示语
	currentTime := time.Now()
	// 格式化时间为：年-月-日 时:分:秒
	timeFormat := currentTime.Format("2006-01-02 15:04:05")
	timeTip := fmt.Sprintf("当前时间为 % s，你的预训练数据存在时效性过期风险，若用户咨询时效性较强的问题，需依据该最新时间进行回复。", timeFormat)

	// 5. 组合输入内容
	finalInput := timeTip + "\n" + pipeInput + "\n" + argInput

	return finalInput, *reasoning, nil
}
