package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/openai/openai-go/v3"
)

type SearchTool struct{
	APIKey string
}

func (t *SearchTool) Name() string {
	return "search"
}

// Definition 定义工具的函数签名，供 LLM 使用
// LLM 会根据需求自动提取几个关键搜索词组成 query
func (t *SearchTool) Definition() openai.ChatCompletionToolUnionParam {
	return openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        t.Name(),
		Description: openai.String("使用关键词搜索互联网，获取最新、最准确的信息。建议 LLM 根据用户需求提取 2-5 个核心关键词组成简洁的 query 参数（例如：原神 最新版本 更新内容）"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "搜索关键词或短语",
				},
			},
			"required": []string{"query"},
		},
	})
}

// 构造 Tavily 请求（仅设置最佳默认配置，其他参数走 API 默认值）
type tavilySearchRequest struct {
	Query         string `json:"query"`
	SearchDepth   string `json:"search_depth"`
	MaxResults    int    `json:"max_results"`
	IncludeAnswer bool   `json:"include_answer"`
	Topic         string `json:"topic"`
}

// tavilyResponse 用于解析 Tavily API 返回的 JSON 结果
type tavilySearchResponse struct {
	Results[]struct {
		Title   string `json:"title"`
		URL     string `json:"url"`
		Content string `json:"content"`
	} `json:"results"`
}

// Execute 执行搜索并返回格式化结果给 LLM
// 内部调用 Tavily API，自动组装「AI 总结 + 搜索结果列表」作为返回值
func (t *SearchTool) Execute(args map[string]any) (string, error) {
	// 1. 提取 query 参数
	query, ok := args["query"].(string)
	if !ok {
		return "", fmt.Errorf("无效的参数: query 必须是字符串")
	}

	reqBody := tavilySearchRequest{
		Query:         query,
		SearchDepth:   "basic",   // 搜索质量（advanced 与 basic 有什么区别？）
		MaxResults:    5,         // 推荐数量
		IncludeAnswer: false,     // 不需要 Tavily 的 AI 总结，让本地 LLM 负责总结即可
		Topic:         "general", // 通用搜索
	}

	payloadBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("构建搜索请求失败: %v", err)
	}

	// 4. 构造 HTTP 请求
	url := "https://api.tavily.com/search"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("创建 HTTP 请求失败: %v", err)
	}

	// 5. 设置必要的请求头
	req.Header.Add("Authorization", "Bearer " + t.APIKey)
	req.Header.Add("Content-Type", "application/json")

	// 6. 发送请求
	client := &http.Client{
    	Timeout: 20 * time.Second, // 防止长时间阻塞
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求 Tavily API 失败: %v", err)
	}
	defer resp.Body.Close()

	// 7. 处理异常的 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Tavily 搜索失败，状态码: %d, 详情: %s", resp.StatusCode, string(body))
	}

	// 8. 解析响应结果
	var apiRes tavilySearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		return "", fmt.Errorf("解析搜索结果失败: %v", err)
	}

	// 9. 将结果组装成易于 LLM 阅读的文本格式
	if len(apiRes.Results) == 0 {
		return fmt.Errorf("未找到关于 [%s] 的任何结果", query).Error(), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("以下是关于 [%s] 的最新搜索结果：\n\n", query))
	for i, item := range apiRes.Results {
		sb.WriteString(fmt.Sprintf("结果 %d:\n", i+1))
		sb.WriteString(fmt.Sprintf("- 标题: %s\n", item.Title))
		sb.WriteString(fmt.Sprintf("- 链接: %s\n", item.URL))
		sb.WriteString(fmt.Sprintf("- 内容: %s\n\n", item.Content))
	}

	// 10. 返回组装好的文本给 LLM
	return sb.String(), nil
}
