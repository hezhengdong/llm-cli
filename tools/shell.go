package tools

import (
	"bufio"
	"fmt"
	"llm-cli/utils"
	"os"
	"os/exec"
	"strings"

	"github.com/openai/openai-go/v3"
)

type ShellTool struct{}

func (t *ShellTool) Name() string {
	return "shell"
}

func (t *ShellTool) Definition() openai.ChatCompletionToolUnionParam {
	return openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        t.Name(),
		Description: openai.String("Execute a shell command on the host machine. USE WITH CAUTION."),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"command": map[string]any{
					"type":        "string",
					"description": "The bash/shell command to execute",
				},
			},
			"required": []string{"command"},
		},
	})
}

func (t *ShellTool) Execute(args map[string]any) (string, error) {
	command, ok := args["command"].(string)
	if !ok {
		return "", fmt.Errorf("invalid argument: command must be a string")
	}

	// 询问人类确认
	utils.GrayPrintf("[请求执行命令 `%s`] 是否允许执行？(y/N): ", command)

	// 直接从终端读取用户输入，避开管道输入
	reader := bufio.NewReader(os.Stdin)
	// 注意：如果是管道输入模式运行，这里获取交互式输入可能会被阻塞，实际生产中可能需要读取 /dev/tty
	confirmation, _ := reader.ReadString('\n')
	confirmation = strings.TrimSpace(strings.ToLower(confirmation))

	if confirmation != "y" && confirmation != "yes" {
		return "User denied the execution of the command.", nil
	}

	// 执行命令
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()

	utils.GrayPrintf("%s", string(out))

	if err != nil {
		return fmt.Sprintf("Execution failed: %v\nOutput: %s", err, string(out)), nil
	}
	return string(out), nil
}
