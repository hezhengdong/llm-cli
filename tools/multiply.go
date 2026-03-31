package tools

import (
	"fmt"

	"github.com/openai/openai-go/v3"
)

type MultiplyTool struct{}

func (t *MultiplyTool) Name() string {
	return "multiply"
}

func (t *MultiplyTool) Definition() openai.ChatCompletionToolUnionParam {
	return openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name: t.Name(),
		Description: openai.String("Multiply two numbers together"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"a": map[string]any{"type": "number"},
				"b": map[string]any{"type": "number"},
			},
			"required": []string{"a", "b"},
		},
	})
}

func (t *MultiplyTool) Execute(args map[string]any) (string, error) {
	a, ok1 := args["a"].(float64)
	b, ok2 := args["b"].(float64)
	if !ok1 || !ok2 {
		return "", fmt.Errorf("invalid arguments")
	}
	return fmt.Sprintf("%f", a*b), nil
}
