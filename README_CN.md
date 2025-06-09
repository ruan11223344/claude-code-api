[English](README.md) | [中文](README_CN.md)

# Claude Code API - Claude Code CLI 的 OpenAI 兼容 API

## 🚀 无缝集成，零代码修改

**无需修改一行代码，就能将 Claude Code 的强大功能集成到您现有项目中！** 这个 API 封装为 Claude Code CLI 提供了完美的 OpenAI 兼容接口，让您可以：

- ✅ **即插即用** - 只需更改 base URL 即可替代 OpenAI API
- ✅ **无需修改代码** - 兼容所有 OpenAI 客户端库
- ✅ **访问 Claude Code 文件操作** - 读取、写入、编辑、搜索文件
- ✅ **执行终端命令** - 直接运行 bash 命令
- ✅ **分析图像** - 处理截图和图表
- ✅ **保持对话上下文** - 恢复和继续会话

## 🤔 订阅用不完？这是个普遍问题

很多开发者在订阅了 Claude API 的 $100 或 $200 月度套餐后，发现实际使用量远远达不到订阅额度。这篇文章将帮助你充分利用你的 API 订阅。

## 💡 Claude API 能做什么？

### 1. 内容创作助手
- 批量生成营销文案
- 自动化博客文章创作
- 社交媒体内容规划和生成

### 2. 代码开发辅助
- 代码审查和优化建议
- 自动生成单元测试
- 技术文档自动化编写
- Bug 分析和修复建议

### 3. 数据分析和报告
- 自动化数据分析报告生成
- 商业智能洞察提取
- 竞品分析和市场研究

### 4. 教育和培训
- 个性化学习助手
- 自动生成练习题和答案
- 知识点总结和解释

## 🚀 充分利用订阅的实用建议

### 1. 构建自己的 AI 应用
```python
# 示例：使用 Claude Code API 批量处理文档
import openai

# 指向你的 Claude Code API 服务器
client = openai.OpenAI(
    api_key="your-api-key",
    base_url="http://localhost:8082/v1"
)

def process_documents(documents):
    results = []
    for doc in documents:
        response = client.chat.completions.create(
            model="claude-code",
            messages=[{"role": "user", "content": f"分析以下文档：{doc}"}],
            claude_options={
                "tools": ["Read", "Grep"],  # 启用文件读取工具
                "working_dir": "/path/to/documents"
            }
        )
        results.append(response.choices[0].message.content)
    return results
```

### 2. 团队共享使用
- 将 API 集成到团队工具中
- 构建内部知识库问答系统
- 自动化日常工作流程

### 3. 开源项目贡献
- 开发 Claude API 的开源工具
- 分享你的创新用例
- 参与社区项目

### 4. 个人项目实验
- 尝试不同的 prompt 工程技巧
- 测试模型在特定领域的能力
- 构建个人助理应用

## 📊 成本优化策略

### 1. 优化 Prompt
- 精简提示词，减少 token 消耗
- 使用系统提示词复用
- 批量处理相似请求

### 2. 缓存策略
- 实现响应缓存机制
- 避免重复相同查询
- 定期清理无用缓存

## 🎯 实际案例

### 案例 1：内容营销自动化
一家电商公司使用 Claude API 每月生成 1000+ 产品描述，大幅提升了上新效率。

### 案例 2：技术文档助手
开发团队集成 Claude API，自动生成 API 文档和代码注释，节省 70% 文档编写时间。

### 案例 3：研究助手
学术研究团队使用 Claude API 分析大量文献，自动生成文献综述，研究效率提升 5 倍。

## 🔧 快速开始

### 前置要求

1. **安装 Claude Code CLI（必须）**
   ```bash
   npm install -g @anthropic-ai/claude-code
   ```

2. **确保已登录（必须）**
   ```bash
   # 如果未登录，claude 命令会自动提示登录
   claude
   ```

### 安装与设置

1. **克隆并构建**
   ```bash
   git clone https://github.com/yourusername/claude-code-api
   cd claude-code-api
   go build
   ```

2. **配置环境**
   ```bash
   cp .env.example .env
   # 编辑 .env 文件设置 API_KEY 和其他配置
   ```

3. **运行服务器**
   ```bash
   # 使用环境文件
   ./claude-code-api
   
   # 或使用环境变量
   API_KEY=your-api-key ./claude-code-api
   ```

### API 认证

如果在环境中设置了 `API_KEY`，所有 API 请求必须包含认证：

```bash
curl -X POST http://localhost:8082/v1/chat/completions \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{"model": "claude-code", "messages": [{"role": "user", "content": "Hello"}]}'
```

### 使用 API

1. **基本用法**
   ```python
   import openai
   
   client = openai.OpenAI(
       api_key="your-api-key",
       base_url="http://localhost:8082/v1"
   )
   
   response = client.chat.completions.create(
       model="claude-code",  # 或使用具体模型名
       messages=[{"role": "user", "content": "Hello, Claude!"}]
   )
   ```

2. **使用 Claude 选项 - 详细示例**

   **代码分析与文件操作：**
   ```python
   response = client.chat.completions.create(
       model="claude-code",
       messages=[{
           "role": "user", 
           "content": "请分析代码结构并提出改进建议。检查是否有安全问题并优化性能。"
       }],
       claude_options={
           "tools": ["Read", "Grep", "Edit"],  # 允许读取、搜索和编辑文件
           "working_dir": "/Users/zhang/myproject",
           "files": ["src/main.py", "src/utils.py", "config.json"]
       }
   )
   ```

   **运行测试并修复问题：**
   ```python
   response = client.chat.completions.create(
       model="claude-code",
       messages=[{
           "role": "user", 
           "content": "运行测试套件，找出失败的测试并修复问题"
       }],
       claude_options={
           "tools": ["Bash", "Read", "Edit"],  # 可以执行命令和修改代码
           "working_dir": "/Users/zhang/myproject",
           "model": "opus"  # 使用最强大的模型
       }
   )
   ```

   **图像分析与代码生成：**
   ```python
   response = client.chat.completions.create(
       model="claude-code",
       messages=[{
           "role": "user", 
           "content": "查看这个 UI 设计稿，创建一个匹配的 React 组件"
       }],
       claude_options={
           "tools": ["Read", "Write"],
           "images": ["/Users/zhang/designs/login-page.png"],
           "working_dir": "/Users/zhang/react-app/src/components"
       }
   )
   ```

   **继续之前的会话：**
   ```python
   response = client.chat.completions.create(
       model="claude-code",
       messages=[{
           "role": "user", 
           "content": "继续我们之前开始的重构工作"
       }],
       claude_options={
           "session_id": "abc123",  # 恢复特定会话
           "tools": ["Read", "Edit", "Bash"]
       }
   )
   ```

   **只读调试（限制工具）：**
   ```python
   response = client.chat.completions.create(
       model="claude-code",
       messages=[{
           "role": "user", 
           "content": "调试应用启动时崩溃的原因"
       }],
       claude_options={
           "tools": ["Read", "Grep"],  # 只读访问
           "disallowed_tools": ["Edit", "Write", "Bash"],  # 禁止修改
           "working_dir": "/Users/zhang/production-app",
           "debug": True  # 启用调试输出
       }
   )
   ```

### 支持的 Claude 选项

- **`tools`**：允许使用的工具列表（例如：`["Bash", "Edit", "Read", "Grep", "WebSearch"]`）
- **`disallowed_tools`**：禁用的工具列表
- **`session_id`**：恢复特定会话
- **`continue`**：继续最近的会话
- **`model`**：指定模型（例如："opus"、"sonnet"、"haiku"）
- **`output_format`**：响应格式（"text"、"json"、"stream-json"）
- **`debug`**：启用调试模式
- **`working_dir`**：设置文件操作的工作目录
- **`mcp_config`**：MCP 服务器配置
- **`files`**：要包含在提示中的文件路径列表
- **`images`**：要分析的图像路径列表
- **`auto_allow_permissions`**：跳过所有权限提示（请谨慎使用！）

## ⚠️ 使用注意事项

- **响应延迟**：API 响应有时较慢，适合非实时性任务
- **批量处理**：建议用于批量数据处理、内容生成等场景
- **异步设计**：在应用中采用异步处理模式，提升用户体验

## 💰 投资回报

- **时间节省**：自动化重复性任务
- **质量提升**：AI 辅助决策和创作
- **规模化**：批量处理能力
- **创新机会**：探索新的业务模式

## 🤝 加入社区

- [Claude API Discord](https://discord.gg/anthropic)
- [GitHub 示例项目](https://github.com/anthropics/anthropic-sdk-python)
- [官方文档](https://docs.anthropic.com)

### 💬 微信交流群

扫描二维码加入微信群，与其他开发者交流 Claude API 使用经验：

![微信群二维码](docs/images/wechat-group.jpg)

---

记住：API 订阅不仅是成本，更是对效率和创新的投资。开始构建，让你的订阅发挥最大价值！

#ClaudeAPI #AI开发 #订阅优化