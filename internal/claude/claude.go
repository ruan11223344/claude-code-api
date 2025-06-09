package claude

import (
	"claude-code-api/internal/logger"
	"encoding/json"
	"fmt"
	"os/exec"
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
		var formatInstruction string
		switch c.outputFormat {
		case "json":
			formatInstruction = "Please return results in valid JSON format without any other text explanation."
		case "markdown":
			formatInstruction = "Please return results in Markdown format."
		case "list":
			formatInstruction = "Please return results as a list, one item per line."
		case "yaml":
			formatInstruction = "Please return results in YAML format."
		case "csv":
			formatInstruction = "Please return results in CSV format with headers in the first row."
		}
		if formatInstruction != "" {
			fullQuestion = fmt.Sprintf("%s\n\n%s", fullQuestion, formatInstruction)
		}
	}
	
	// Log request details
	logger.Log.Infof("[CLAUDE] Starting request")
	logger.Log.Infof("[CLAUDE] Question length: %d characters", len(question))
	if c.systemPrompt != "" {
		logger.Log.Infof("[CLAUDE] Using system prompt")
	}
	if c.outputFormat != "text" {
		logger.Log.Infof("[CLAUDE] Output format: %s", c.outputFormat)
	}
	
	// Use claude's --print option for non-interactive calls
	cmd := exec.Command("claude", "--print", fullQuestion)
	
	// Execute command and get output
	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)
	
	if err != nil {
		logger.Log.Infof("[CLAUDE] ❌ Request failed after %v: %v", duration, err)
		if len(output) > 0 {
			logger.Log.Infof("[CLAUDE] Error output: %s", string(output))
		}
		return "", fmt.Errorf("claude execution failed: %v", err)
	}
	
	response := strings.TrimSpace(string(output))
	logger.Log.Infof("[CLAUDE] ✅ Request completed in %v", duration)
	logger.Log.Infof("[CLAUDE] Response length: %d characters", len(response))
	
	return response, nil
}

// AskWithOptions sends a question to Claude with custom CLI options
// Options can include:
// - "tools": []string{"Bash", "Edit", "Read"} - allowed tools
// - "disallowed_tools": []string{"Bash"} - disallowed tools  
// - "session_id": "abc123" - resume a specific session
// - "continue": true - continue the most recent conversation
// - "model": "opus" - model to use
// - "output_format": "json" - output format (text, json, stream-json)
// - "debug": true - enable debug mode
// - "working_dir": "/path/to/dir" - working directory for file operations
// - "mcp_config": map or string - MCP server configuration
// - "files": []string{"/path/to/file1", "/path/to/file2"} - files to include
// - "images": []string{"/path/to/image1.png"} - images to include
// - "auto_allow_permissions": true - automatically allow all tool permissions (use with caution)
func (c *Client) AskWithOptions(question string, options map[string]interface{}) (string, error) {
	startTime := time.Now()
	args := []string{}
	
	// Add print flag unless streaming is requested
	if format, ok := options["output_format"].(string); !ok || format != "stream-json" {
		args = append(args, "--print")
	}
	
	// Check if auto-allow permissions is enabled
	if autoAllow, ok := options["auto_allow_permissions"].(bool); ok && autoAllow {
		args = append(args, "--dangerously-skip-permissions")
	}
	
	// Process CLI options
	for key, value := range options {
		switch key {
		case "tools":
			if tools, ok := value.([]string); ok && len(tools) > 0 {
				args = append(args, "--allowedTools", strings.Join(tools, ","))
			}
		case "disallowed_tools":
			if tools, ok := value.([]string); ok && len(tools) > 0 {
				args = append(args, "--disallowedTools", strings.Join(tools, ","))
			}
		case "session_id":
			if sessionID, ok := value.(string); ok && sessionID != "" {
				args = append(args, "--resume", sessionID)
			}
		case "continue":
			if cont, ok := value.(bool); ok && cont {
				args = append(args, "--continue")
			}
		case "model":
			if model, ok := value.(string); ok && model != "" {
				args = append(args, "--model", model)
			}
		case "output_format":
			if format, ok := value.(string); ok && format != "" {
				args = append(args, "--output-format", format)
			}
		case "debug":
			if debug, ok := value.(bool); ok && debug {
				args = append(args, "--debug")
			}
		case "mcp_config":
			if config, ok := value.(string); ok {
				args = append(args, "--mcp-config", config)
			} else if configMap, ok := value.(map[string]interface{}); ok {
				// Convert map to JSON string
				if jsonBytes, err := json.Marshal(configMap); err == nil {
					args = append(args, "--mcp-config", string(jsonBytes))
				}
			}
		case "files", "images":
			// These will be handled in the prompt construction
		case "working_dir":
			// This will be handled by setting cmd.Dir
		}
	}
	
	// Construct the full prompt with files/images if provided
	fullQuestion := question
	
	// Add file contents to the prompt
	if files, ok := options["files"].([]string); ok && len(files) > 0 {
		for _, file := range files {
			fullQuestion += fmt.Sprintf("\n\nFile: %s", file)
		}
	}
	
	// Add image references to the prompt
	if images, ok := options["images"].([]string); ok && len(images) > 0 {
		for _, image := range images {
			fullQuestion += fmt.Sprintf("\n\nImage: %s", image)
		}
	}
	
	// Add system prompt if configured
	if c.systemPrompt != "" {
		fullQuestion = fmt.Sprintf("%s\n\n%s", c.systemPrompt, fullQuestion)
	}
	
	// Add the question/prompt at the end
	args = append(args, fullQuestion)
	
	// Log request details
	logger.Log.Infof("[CLAUDE] Starting request with options")
	logger.Log.Infof("[CLAUDE] CLI args: %v", args)
	
	// Create command
	cmd := exec.Command("claude", args...)
	
	// Set working directory if specified
	if workDir, ok := options["working_dir"].(string); ok {
		cmd.Dir = workDir
		logger.Log.Infof("[CLAUDE] Working directory: %s", workDir)
	}
	
	// Execute command
	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)
	
	if err != nil {
		logger.Log.Infof("[CLAUDE] ❌ Request failed after %v: %v", duration, err)
		if len(output) > 0 {
			logger.Log.Infof("[CLAUDE] Error output: %s", string(output))
		}
		return "", fmt.Errorf("claude execution failed: %v", err)
	}
	
	response := strings.TrimSpace(string(output))
	logger.Log.Infof("[CLAUDE] ✅ Request completed in %v", duration)
	logger.Log.Infof("[CLAUDE] Response length: %d characters", len(response))
	
	return response, nil
}


// Default system prompt for Claude
const defaultSystemPrompt = `You are a professional AI assistant capable of:
1. Understanding and analyzing various types of files (code, documents, images, data, etc.)
2. Providing accurate and professional analysis and suggestions
3. Returning results in the requested format
4. Keeping responses concise and to the point

When processing files, you will:
- Automatically identify file types and use appropriate processing methods
- For code files, understand their functionality and structure
- For image files, describe content and extract relevant information
- For data files, analyze structure and extract key information

Please always maintain a professional, accurate, and helpful attitude.`