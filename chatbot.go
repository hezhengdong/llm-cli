package main

import (
	"context"
	"encoding/json"
	"fmt"
	"llm-cli/tools"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type ChatBot struct {
	client *openai.Client
	model string
	messages []openai.ChatCompletionMessageParamUnion
	tools    map[string]tools.Tool
}

func NewChatBot(apiKey, baseUrl, model, systemPrompt string) *ChatBot {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseUrl),
	)

	// 初始化时先把 System Prompt 塞进对话列表
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
	}

	bot := &ChatBot{
		client:   &client,
		model:    model,
		messages: messages,
		tools:    make(map[string]tools.Tool),
	}

	return bot
}

func (b *ChatBot) RegisterTool(tool tools.Tool) {
	b.tools[tool.Name()] = tool
}

func (b *ChatBot) getToolDefinitions() []openai.ChatCompletionToolUnionParam {
    var tools []openai.ChatCompletionToolUnionParam
    for _, t := range b.tools {
        tools = append(tools, t.Definition()) // 工具定义，这个倒是简单好理解
    }
    return tools
}

func (b *ChatBot) llmNode(ctx context.Context) (openai.ChatCompletionMessageParamUnion, []openai.ChatCompletionMessageToolCallUnion) {
	params := openai.ChatCompletionNewParams{
		Messages: b.messages, // 这里的 b.messages 已经是最新状态
		Seed:     openai.Int(0),
		Model:    b.model,
	}

	// 绑定工具，每次 LLM 都必须绑定工具
	if tools := b.getToolDefinitions(); len(tools) > 0 {
		params.Tools = tools
	}

	stream := b.client.Chat.Completions.NewStreaming(ctx, params)
	acc := openai.ChatCompletionAccumulator{}
	var replyBuilder strings.Builder

	fmt.Print("[🤖 LLM 回复] ")
	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			content := chunk.Choices[0].Delta.Content
			fmt.Print(content)
			replyBuilder.WriteString(content)
		}
	}

	if replyBuilder.Len() > 0 {
		fmt.Printf("\n\n")
	}

	if err := stream.Err(); err != nil {
		panic(err.Error())
	}

	return acc.Choices[0].Message.ToParam(), acc.Choices[0].Message.ToolCalls
}

func (b *ChatBot) toolNode(toolCalls []openai.ChatCompletionMessageToolCallUnion) []openai.ChatCompletionMessageParamUnion {
	var toolMessages []openai.ChatCompletionMessageParamUnion
	for _, tc := range toolCalls {
		toolName := tc.Function.Name
		argsJSON := tc.Function.Arguments

		tool, exists := b.tools[toolName]

		// LLM 发生幻觉，调用了不存在的工具
		if !exists {
			errMsg := fmt.Sprintf("System Error: Tool '%s' not found.", toolName)
			toolMessages = append(toolMessages, openai.ToolMessage(errMsg, tc.ID))
			continue
		}

		// 解释 LLM 传来的 JSON 参数
		var args map[string]any
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			// 这里面是解析失败的情况，把报错内容反馈给 LLM，让 LLM 自行判断下一步该怎么做
			errMsg := fmt.Sprintf("System Error parsing arguments: %v", err)
			toolMessages = append(toolMessages, openai.ToolMessage(errMsg, tc.ID))
			continue
		}

		fmt.Printf("[🛠️ LLM 正在调用 %s 工具] {toolName: %s, argsJSON: %s}\n\n", toolName, toolName, argsJSON)

		// 根据参数执行 Execute 函数
		result, err := tool.Execute(args)
		if err != nil {
			// 工具报错也应返给 LLM，让 LLM 知道执行失败并尝试自我修复
			result = fmt.Sprintf("Execution error: %v", err)
		}

		fmt.Printf("\n[🛠️ 工具执行完毕] \n%s\n\n", result)
		toolMessages = append(toolMessages, openai.ToolMessage(result, tc.ID))
	}

	return toolMessages
}

func (b *ChatBot) Chat(ctx context.Context, input string) {
	const maxIterations = 999

	// 添加用户给的提示词
	b.messages = append(b.messages, openai.UserMessage(input))

	for iter := 0; iter < maxIterations; iter++ {
		// 调用 LLM 节点
		param, toolCalls := b.llmNode(ctx)
		// 将 AI 的输出添加到 messages 中
		b.messages = append(b.messages, param)

		// 没有工具
		if len(toolCalls) == 0 {
			return
		}

		toolMessages := b.toolNode(toolCalls)

		b.messages = append(b.messages, toolMessages...)
	}
}
