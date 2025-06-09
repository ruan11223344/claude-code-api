package claude

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Client represents a Claude AI client
type Client struct {
	systemPrompt string
	outputFormat string // json, text, markdown etc.
}

// Option is a function that configures the Client
type Option func(*Client)

// WithSystemPrompt sets a custom system prompt
func WithSystemPrompt(prompt string) Option {
	return func(c *Client) {
		c.systemPrompt = prompt
	}
}

// WithOutputFormat sets the desired output format
func WithOutputFormat(format string) Option {
	return func(c *Client) {
		c.outputFormat = format
	}
}

// New creates a new Claude client
func New(opts ...Option) *Client {
	c := &Client{
		systemPrompt: defaultSystemPrompt,
		outputFormat: "text",
	}
	
	for _, opt := range opts {
		opt(c)
	}
	
	return c
}

// Ask sends a question to Claude and returns the response
func (c *Client) Ask(question string) (string, error) {
	startTime := time.Now()
	
	// Combine system prompt with question if provided
	fullQuestion := question
	if c.systemPrompt != "" {
		fullQuestion = fmt.Sprintf("%s\n\n%s", c.systemPrompt, question)
	}
	
	// Add output format instruction if specified
	if c.outputFormat != "" && c.outputFormat != "text" {
		formatInstruction := c.getFormatInstruction()
		fullQuestion = fmt.Sprintf("%s\n\n%s", fullQuestion, formatInstruction)
	}
	
	// Log request details
	log.Printf("[CLAUDE] Starting request")
	log.Printf("[CLAUDE] Question length: %d characters", len(question))
	if c.systemPrompt != "" {
		log.Printf("[CLAUDE] Using system prompt")
	}
	if c.outputFormat != "text" {
		log.Printf("[CLAUDE] Output format: %s", c.outputFormat)
	}
	
	// Use claude's --print option for non-interactive calls
	cmd := exec.Command("claude", "--print", fullQuestion)
	
	// Execute command and get output
	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)
	
	if err != nil {
		log.Printf("[CLAUDE] ❌ Request failed after %v: %v", duration, err)
		if len(output) > 0 {
			log.Printf("[CLAUDE] Error output: %s", string(output))
		}
		return "", fmt.Errorf("claude execution failed: %v", err)
	}
	
	response := strings.TrimSpace(string(output))
	log.Printf("[CLAUDE] ✅ Request completed in %v", duration)
	log.Printf("[CLAUDE] Response length: %d characters", len(response))
	
	return response, nil
}

// AskWithFile sends a question along with a file path to Claude
func (c *Client) AskWithFile(question string, filePath string) (string, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("文件不存在: %v", err)
	}
	
	// Convert to absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("获取绝对路径失败: %v", err)
	}
	
	// Claude can handle file paths directly
	fullQuestion := fmt.Sprintf("%s %s", question, absPath)
	
	return c.Ask(fullQuestion)
}

// AskWithFiles sends a question along with multiple file paths to Claude
func (c *Client) AskWithFiles(question string, filePaths []string) (string, error) {
	// Validate all files exist and convert to absolute paths
	var absPaths []string
	for _, filePath := range filePaths {
		if _, err := os.Stat(filePath); err != nil {
			return "", fmt.Errorf("文件不存在 %s: %v", filePath, err)
		}
		
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			return "", fmt.Errorf("获取绝对路径失败 %s: %v", filePath, err)
		}
		absPaths = append(absPaths, absPath)
	}
	
	// Claude can handle multiple file paths directly
	fullQuestion := fmt.Sprintf("%s %s", question, strings.Join(absPaths, " "))
	return c.Ask(fullQuestion)
}

// AskWithImage is a convenience method for asking questions about images
func (c *Client) AskWithImage(question string, imagePath string) (string, error) {
	return c.AskWithFile(question, imagePath)
}

// AnalyzeCode sends code for analysis
func (c *Client) AnalyzeCode(code string, language string) (string, error) {
	question := fmt.Sprintf("请分析以下%s代码:\n```%s\n%s\n```", language, language, code)
	return c.Ask(question)
}

// GenerateCode generates code based on requirements
func (c *Client) GenerateCode(requirements string, language string) (string, error) {
	question := fmt.Sprintf("请用%s语言生成满足以下需求的代码:\n%s", language, requirements)
	return c.Ask(question)
}

// ExecuteTask asks Claude to perform system operations
func (c *Client) ExecuteTask(task string) (string, error) {
	// Claude can execute various system operations directly
	return c.Ask(task)
}

// CreateDirectory asks Claude to create a directory
func (c *Client) CreateDirectory(path string) (string, error) {
	task := fmt.Sprintf("请创建目录: %s", path)
	return c.ExecuteTask(task)
}

// DownloadFile asks Claude to download a file
func (c *Client) DownloadFile(url, savePath string) (string, error) {
	task := fmt.Sprintf("请下载文件 %s 并保存到 %s", url, savePath)
	return c.ExecuteTask(task)
}

// EditFile asks Claude to edit a file
func (c *Client) EditFile(filePath, instructions string) (string, error) {
	task := fmt.Sprintf("请编辑文件 %s，具体要求：%s", filePath, instructions)
	return c.ExecuteTask(task)
}

// CreateFile asks Claude to create a new file with content
func (c *Client) CreateFile(filePath, content string) (string, error) {
	task := fmt.Sprintf("请创建文件 %s，内容如下：\n%s", filePath, content)
	return c.ExecuteTask(task)
}

// RunCommand asks Claude to execute a shell command
func (c *Client) RunCommand(command string) (string, error) {
	task := fmt.Sprintf("请执行命令: %s", command)
	return c.ExecuteTask(task)
}

// InstallPackage asks Claude to install a package
func (c *Client) InstallPackage(packageName string, packageManager string) (string, error) {
	var task string
	switch packageManager {
	case "npm":
		task = fmt.Sprintf("请使用npm安装包: npm install %s", packageName)
	case "pip":
		task = fmt.Sprintf("请使用pip安装包: pip install %s", packageName)
	case "go":
		task = fmt.Sprintf("请使用go get安装包: go get %s", packageName)
	default:
		task = fmt.Sprintf("请安装包: %s", packageName)
	}
	return c.ExecuteTask(task)
}

// SearchAndReplace asks Claude to search and replace in files
func (c *Client) SearchAndReplace(pattern, replacement, filePattern string) (string, error) {
	task := fmt.Sprintf("请在匹配 %s 的文件中，将 %s 替换为 %s", filePattern, pattern, replacement)
	return c.ExecuteTask(task)
}

// GitOperation asks Claude to perform git operations
func (c *Client) GitOperation(operation string) (string, error) {
	task := fmt.Sprintf("请执行git操作: %s", operation)
	return c.ExecuteTask(task)
}

// ProcessFiles asks Claude to process multiple files with specific instructions
func (c *Client) ProcessFiles(filePattern string, instructions string) (string, error) {
	task := fmt.Sprintf("请处理匹配 %s 的所有文件，具体要求：%s", filePattern, instructions)
	return c.ExecuteTask(task)
}

// AutomateWorkflow asks Claude to execute a complex workflow
func (c *Client) AutomateWorkflow(workflow string) (string, error) {
	task := fmt.Sprintf("请执行以下工作流程：\n%s", workflow)
	return c.ExecuteTask(task)
}

// getFormatInstruction returns format instruction based on output format
func (c *Client) getFormatInstruction() string {
	switch c.outputFormat {
	case "json":
		return "请以有效的JSON格式返回结果，不要包含其他文字说明。"
	case "markdown":
		return "请以Markdown格式返回结果。"
	case "list":
		return "请以列表形式返回结果，每项一行。"
	case "yaml":
		return "请以YAML格式返回结果。"
	case "csv":
		return "请以CSV格式返回结果，第一行为标题。"
	default:
		return ""
	}
}

// AskJSON asks a question and expects a JSON response
func (c *Client) AskJSON(question string, result interface{}) error {
	// Temporarily set output format to JSON
	originalFormat := c.outputFormat
	c.outputFormat = "json"
	defer func() { c.outputFormat = originalFormat }()
	
	response, err := c.Ask(question)
	if err != nil {
		return err
	}
	
	// Parse JSON response
	if err := json.Unmarshal([]byte(response), result); err != nil {
		return fmt.Errorf("解析JSON响应失败: %v", err)
	}
	
	return nil
}

// AskStructured asks a question with a structured template
func (c *Client) AskStructured(template string, data interface{}) (string, error) {
	// Format template with data
	question := fmt.Sprintf(template, data)
	return c.Ask(question)
}

// AnalyzeImage analyzes an image and returns insights
func (c *Client) AnalyzeImage(imagePath string, analysisType string) (string, error) {
	prompts := map[string]string{
		"general":     "请描述这张图片的内容。",
		"products":    "请识别并描述图片中的产品，包括品牌、型号、特征等信息。",
		"text":        "请提取图片中的所有文字内容。",
		"colors":      "请分析图片的主要颜色和配色方案。",
		"composition": "请分析图片的构图和视觉元素。",
	}
	
	prompt, ok := prompts[analysisType]
	if !ok {
		prompt = prompts["general"]
	}
	
	return c.AskWithFile(prompt, imagePath)
}

// ExtractData extracts structured data from files
func (c *Client) ExtractData(filePath string, dataType string) (string, error) {
	prompt := fmt.Sprintf("请从文件中提取%s数据，并以结构化格式返回。", dataType)
	return c.AskWithFile(prompt, filePath)
}

// TranslateFile translates file content to target language
func (c *Client) TranslateFile(filePath string, targetLang string) (string, error) {
	prompt := fmt.Sprintf("请将文件内容翻译成%s，保持原有格式。", targetLang)
	return c.AskWithFile(prompt, filePath)
}

// SummarizeFile creates a summary of file content
func (c *Client) SummarizeFile(filePath string, maxWords int) (string, error) {
	prompt := fmt.Sprintf("请用不超过%d字总结文件的主要内容。", maxWords)
	return c.AskWithFile(prompt, filePath)
}

// CompareFiles compares multiple files and highlights differences
func (c *Client) CompareFiles(filePaths []string) (string, error) {
	prompt := "请比较这些文件的内容，指出主要差异和相似之处。"
	return c.AskWithFiles(prompt, filePaths)
}

// Default system prompt for Claude
const defaultSystemPrompt = `你是一个专业的AI助手，能够：
1. 理解和分析各类文件（代码、文档、图片、数据等）
2. 提供准确、专业的分析和建议
3. 按照要求的格式返回结果
4. 保持回答简洁明了，直接切中要点

在处理文件时，你会：
- 自动识别文件类型并采用合适的处理方式
- 对于代码文件，理解其功能和结构
- 对于图片文件，描述内容并提取相关信息
- 对于数据文件，分析结构并提取关键信息

请始终保持专业、准确和有帮助的态度。`