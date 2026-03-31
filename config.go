package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ToolsConfig struct {
	Tavily struct {
		APIKey string `yaml:"api_key"`
	} `yaml:"tavily"`
}

type ProviderConfig struct {
    BaseURL string `yaml:"base_url"`
    APIKey  string `yaml:"api_key"`
    Model   string `yaml:"model"`
}

type Config struct {
    DefaultProvider string                    `yaml:"default_provider"`
    DefaultPrompt   string                    `yaml:"default_prompt"`
    Providers       map[string]ProviderConfig `yaml:"providers"`
    Prompts         map[string]string         `yaml:"prompts"`
    Tools           ToolsConfig               `yaml:"tools"`
}

// loadConfig 读取并解析配置文件
func loadConfig(configPath string) Config {

	// 读取 config.yaml 文件
    yamlFile, err := os.ReadFile(configPath)
    if err != nil {
        panic(fmt.Sprintf("读取配置文件失败 (%s): %v\n尝试创建配置文件或者检查路径。", configPath, err))
    }

    var cfg Config
    err = yaml.Unmarshal(yamlFile, &cfg)
    if err != nil {
        panic("解析配置文件失败: " + err.Error())
    }

    return cfg
}