# llm-cli

`llm-cli` 是一款跨平台的本地大语言模型终端对话工具。支持在任意目录下通过 `llm` 命令快速唤起！

## 操作指南

### 安装

#### macOS / Linux

我们在终端提供了一键安装脚本，自动识别您的系统架构、处理 macOS 拦截机制并初始化所需配置：

```bash
# 执行一键安装脚本
curl -sL https://raw.githubusercontent.com/hezhengdong/llm-cli/main/scripts/install.sh | bash
```

*(注意：安装过程中因需要将二进制文件移动至 `/usr/local/bin`，可能需要您输入系统管理员密码授权。)*

#### Windows

对于 Windows 用户，建议使用 **Git Bash** 或者 **WSL** 执行上述指令。执行完毕后，系统会在您用户目录的 `~/bin` 下生成可执行文件，请确保您的系统环境变量 `PATH` 中包含该目录。

### 配置文件

安装成功后，脚本会自动检查并创建配置文件目录：`~/.config/llm-cli/config.yaml`。若该文件已存在，则**不会进行任何覆盖操作**，确保您的数据安全。

初次安装后，请打开该文件，并填入您真实的 `API-KEY`。

```bash
# 您可以使用 vim 或 nano 编辑配置
vim ~/.config/llm-cli/config.yaml
```

### 更新

若发布了新版本，您可以直接在终端输入以下命令进行无痛升级。更新过程不会覆写或清除您的 `config.yaml` 个人配置。

```bash
curl -sL https://raw.githubusercontent.com/hezhengdong/llm-cli/main/scripts/update.sh | bash
```

### 卸载

如果您不再需要使用该程序，可通过以下脚本自动卸载。卸载过程中会询问您是否保留 `config.yaml` 配置文件。

```bash
curl -sL https://raw.githubusercontent.com/hezhengdong/llm-cli/main/scripts/uninstall.sh | bash
```

## TODO

继续完善文档

---

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
- [x] CI DI
- [ ] SQLite 持久化（感觉没有必要）
- [ ] Python 解释器、MCP、Skills（不是很想手搓了）

其他：

- 向量数据库
    - 根据 [Claude Code Doesn't Index Your Codebase](https://vadim.blog/claude-code-no-indexing) 这篇文章所述，使用类似 grep 的终端工具即可很好的完成本地搜索任务。
- 自行处理标准输入
    - 尝试过，容易遇到 bug，而且这属于体力活，直接使用系统终端输入即可。

