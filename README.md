一个终端运行的 LLM 工具

- 熟悉一下基本的 [OpenAI Chat Completions API](https://developers.openai.com/api/reference/resources/chat/subresources/completions/methods/create)。
- 熟悉 LLM 的基本概念，比如流式输出、函数调用。

参考链接

- [You Should Write An Agent - fly.io](https://fly.io/blog/everyone-write-an-agent/)
- [Language models on the command-line](https://simonwillison.net/2024/Jun/17/cli-language-models/)
- [openai-go SDK 的四个 Chat Completions API 示例](https://github.com/openai/openai-go/tree/main/examples)

功能：

- [x] 跑通最基本的流式输出。
- [x] yaml 文件接耦配置 apiKey、baseUrl、model 等配置。
- [x] 接收管道与命令行输入。
- [x] 添加工具调用，乘法工具、搜索工具、终端执行工具
- [ ] CI DI
- [ ] SQLite 持久化（感觉没有必要）
- [ ] Python 解释器、MCP、Skills（不是很想手搓了）

其他：

- 向量数据库
    - 根据 [Claude Code Doesn't Index Your Codebase](https://vadim.blog/claude-code-no-indexing) 这篇文章所述，使用类似 grep 的终端工具即可很好的完成本地搜索任务。
- 自行处理标准输入
    - 尝试过，容易遇到 bug，而且这属于体力活，直接使用系统终端输入即可。

