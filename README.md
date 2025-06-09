[English](README.md) | [中文](README_CN.md)

# Claude Code API - OpenAI-Compatible API for Claude Code CLI

## 🚀 Seamless Integration, Zero Code Changes

**Integrate Claude Code's powerful capabilities into your existing projects without modifying a single line of code!** This API wrapper provides a perfect OpenAI-compatible interface for Claude Code CLI, allowing you to:

- ✅ **Drop-in replacement** for OpenAI API - just change the base URL
- ✅ **No code modifications needed** - works with all OpenAI client libraries
- ✅ **Access Claude Code's file operations** - read, write, edit, search files
- ✅ **Run terminal commands** - execute bash commands directly
- ✅ **Analyze images** - process screenshots and diagrams
- ✅ **Maintain conversation context** - resume and continue sessions

## 🤔 Can't Use Up Your Subscription? You're Not Alone

Many developers find themselves with unused credits after subscribing to Claude API's $100 or $200 monthly plans. This guide will help you maximize the value of your API subscription.

## 💡 What Can You Do with Claude API?

### 1. Content Creation Assistant
- Bulk generate marketing copy
- Automate blog post creation
- Social media content planning and generation

### 2. Code Development Assistant
- Code review and optimization suggestions
- Automated unit test generation
- Technical documentation automation
- Bug analysis and fix recommendations

### 3. Data Analysis & Reporting
- Automated data analysis report generation
- Business intelligence insight extraction
- Competitive analysis and market research

### 4. Education & Training
- Personalized learning assistants
- Auto-generate practice questions and answers
- Knowledge summarization and explanation

## 🚀 Practical Tips to Maximize Your Subscription

### 1. Build Your Own AI Applications
```python
# Example: Batch document processing using Claude Code API
import openai

# Point to your Claude Code API server
client = openai.OpenAI(
    api_key="your-api-key",
    base_url="http://localhost:8082/v1"
)

def process_documents(documents):
    results = []
    for doc in documents:
        response = client.chat.completions.create(
            model="claude-code",
            messages=[{"role": "user", "content": f"Analyze this document: {doc}"}],
            claude_options={
                "tools": ["Read", "Grep"],  # Enable file reading tools
                "working_dir": "/path/to/documents"
            }
        )
        results.append(response.choices[0].message.content)
    return results
```

### 2. Team Collaboration
- Integrate API into team tools
- Build internal knowledge base Q&A systems
- Automate daily workflows

### 3. Open Source Contributions
- Develop open source tools for Claude API
- Share your innovative use cases
- Participate in community projects

### 4. Personal Project Experiments
- Try different prompt engineering techniques
- Test model capabilities in specific domains
- Build personal assistant applications

## 📊 Cost Optimization Strategies

### 1. Optimize Prompts
- Streamline prompts to reduce token usage
- Reuse system prompts
- Batch process similar requests

### 2. Caching Strategy
- Implement response caching
- Avoid duplicate queries
- Regular cache cleanup

## 🎯 Real-World Examples

### Case 1: Content Marketing Automation
An e-commerce company uses Claude API to generate 1000+ product descriptions monthly, significantly improving product launch efficiency.

### Case 2: Technical Documentation Assistant
Development teams integrate Claude API to auto-generate API docs and code comments, saving 70% of documentation time.

### Case 3: Research Assistant
Academic research teams use Claude API to analyze large volumes of literature, auto-generate literature reviews, increasing research efficiency by 5x.

## 🚀 Advanced Features

This API wrapper supports all Claude Code CLI features through the `claude_options` field:

```python
# Example: Using Claude with file operations
response = client.chat.completions.create(
    model="claude-code",
    messages=[{"role": "user", "content": "Analyze this codebase and suggest improvements"}],
    claude_options={
        "tools": ["Bash", "Edit", "Read", "Grep"],  # Enable specific tools
        "working_dir": "/path/to/project",          # Set working directory
        "session_id": "abc123",                     # Resume a session
        "model": "opus",                            # Use a specific model
        "files": ["config.json", "main.py"],       # Include files in context
        "images": ["diagram.png"]                   # Include images
    }
)
```

### Supported Claude Options

- **`tools`**: List of allowed tools (e.g., `["Bash", "Edit", "Read", "Grep", "WebSearch"]`)
- **`disallowed_tools`**: List of tools to disable
- **`session_id`**: Resume a specific conversation
- **`continue`**: Continue the most recent conversation
- **`model`**: Specify model (e.g., "opus", "sonnet", "haiku")
- **`output_format`**: Response format ("text", "json", "stream-json")
- **`debug`**: Enable debug mode
- **`working_dir`**: Set working directory for file operations
- **`mcp_config`**: MCP server configuration
- **`files`**: List of file paths to include in the prompt
- **`images`**: List of image paths to analyze
- **`auto_allow_permissions`**: Skip all permission prompts (use with caution!)

## 🔧 Quick Start

### Prerequisites

1. **Install Claude Code CLI (Required)**
   ```bash
   npm install -g @anthropic-ai/claude-code
   ```

2. **Ensure You're Logged In (Required)**
   ```bash
   # The claude command will automatically prompt for login if needed
   claude
   ```

### Installation & Setup

1. **Clone and Build**
   ```bash
   git clone https://github.com/yourusername/claude-code-api
   cd claude-code-api
   go build
   ```

2. **Configure Environment**
   ```bash
   cp .env.example .env
   # Edit .env file to set your API_KEY and other configurations
   ```

3. **Run the Server**
   ```bash
   # With environment file
   ./claude-code-api
   
   # Or with environment variables
   API_KEY=your-api-key ./claude-code-api
   ```

### API Authentication

If `API_KEY` is set in your environment, all API requests must include authentication:

```bash
curl -X POST http://localhost:8082/v1/chat/completions \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{"model": "claude-3-sonnet-20240229", "messages": [{"role": "user", "content": "Hello"}]}'
```

### Using the API

#### 🎉 Works with Your Existing Code!

**If you're already using OpenAI's API, just change the base URL - that's it!**

1. **Basic Usage - No Code Changes Required**
   ```python
   import openai
   
   # Just change the base_url - everything else stays the same!
   client = openai.OpenAI(
       api_key="your-api-key",
       base_url="http://localhost:8082/v1"  # ← Only change needed!
   )
   
   # Use it exactly like OpenAI API
   response = client.chat.completions.create(
       model="claude-code",
       messages=[{"role": "user", "content": "Hello, Claude!"}]
   )
   ```

   **Works with any OpenAI-compatible library:**
   ```javascript
   // Node.js / JavaScript
   const OpenAI = require('openai');
   const client = new OpenAI({
     apiKey: 'your-api-key',
     baseURL: 'http://localhost:8082/v1'
   });
   ```

   ```csharp
   // C# / .NET
   var client = new OpenAIClient(
       new Uri("http://localhost:8082/v1"),
       new AzureKeyCredential("your-api-key")
   );
   ```

2. **With Claude Options - Examples**

   **Code Analysis with File Operations:**
   ```python
   response = client.chat.completions.create(
       model="claude-code",
       messages=[{
           "role": "user", 
           "content": "Please analyze the code structure and suggest improvements. Check for any security issues and optimize the performance."
       }],
       claude_options={
           "tools": ["Read", "Grep", "Edit"],  # Allow reading, searching, and editing files
           "working_dir": "/Users/john/myproject",
           "files": ["src/main.py", "src/utils.py", "config.json"]
       }
   )
   ```

   **Running Tests and Fixing Issues:**
   ```python
   response = client.chat.completions.create(
       model="claude-code",
       messages=[{
           "role": "user", 
           "content": "Run the test suite, identify failing tests, and fix the issues"
       }],
       claude_options={
           "tools": ["Bash", "Read", "Edit"],  # Can run commands and fix code
           "working_dir": "/Users/john/myproject",
           "model": "opus"  # Use the most capable model
       }
   )
   ```

   **Image Analysis with Code Generation:**
   ```python
   response = client.chat.completions.create(
       model="claude-code",
       messages=[{
           "role": "user", 
           "content": "Look at this UI design mockup and create a React component that matches it"
       }],
       claude_options={
           "tools": ["Read", "Write"],
           "images": ["/Users/john/designs/login-page.png"],
           "working_dir": "/Users/john/react-app/src/components"
       }
   )
   ```

   **Continue Previous Session:**
   ```python
   response = client.chat.completions.create(
       model="claude-code",
       messages=[{
           "role": "user", 
           "content": "Continue working on the refactoring we started earlier"
       }],
       claude_options={
           "session_id": "abc123",  # Resume specific session
           "tools": ["Read", "Edit", "Bash"]
       }
   )
   ```

   **Debugging with Restricted Tools:**
   ```python
   response = client.chat.completions.create(
       model="claude-code",
       messages=[{
           "role": "user", 
           "content": "Debug why the application is crashing on startup"
       }],
       claude_options={
           "tools": ["Read", "Grep"],  # Read-only access
           "disallowed_tools": ["Edit", "Write", "Bash"],  # Prevent modifications
           "working_dir": "/Users/john/production-app",
           "debug": True  # Enable debug output
       }
   )
   ```

   **Auto-approve All Operations (Use with Caution!):**
   ```python
   response = client.chat.completions.create(
       model="claude-code",
       messages=[{
           "role": "user", 
           "content": "Refactor the entire codebase to use async/await"
       }],
       claude_options={
           "tools": ["Read", "Edit", "Write", "Bash"],
           "working_dir": "/Users/john/myproject",
           "auto_allow_permissions": True  # Skip ALL permission prompts!
       }
   )
   ```

## ⚠️ Important Considerations

- **Response Latency**: API responses can be slow at times, suitable for non-real-time tasks
- **Batch Processing**: Recommended for batch data processing, content generation scenarios
- **Asynchronous Design**: Use async processing patterns in your applications for better user experience

## 💰 Return on Investment

- **Time Savings**: Automate repetitive tasks
- **Quality Improvement**: AI-assisted decision making and creation
- **Scalability**: Batch processing capabilities
- **Innovation Opportunities**: Explore new business models

## 🤝 Join the Community

- [Claude API Discord](https://discord.gg/anthropic)
- [GitHub Examples](https://github.com/anthropics/anthropic-sdk-python)
- [Official Documentation](https://docs.anthropic.com)

## 💖 Support This Project

If you find this project helpful, consider supporting its development:

### Donate USDT (TRC20)
![USDT Donation QR Code](docs/images/usdt-tron-donation.png)

**USDT (TRC20) Address:**
```
TBQcKHqDtYw17KrZvkEPeK3TLP96SunmgX
```

---

Remember: Your API subscription isn't just a cost—it's an investment in efficiency and innovation. Start building and maximize your subscription value!

#ClaudeAPI #AIDevelopment #SubscriptionOptimization