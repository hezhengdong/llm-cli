package tools

import "github.com/openai/openai-go/v3"

// 定义标准工具接口
type Tool interface {
	Name() string
	Definition() openai.ChatCompletionToolUnionParam
	Execute(args map[string]any) (string, error)
}
