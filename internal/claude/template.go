package claude

import (
	"encoding/json"
	"fmt"
	"os"
)

// TemplateClient handles template-based interactions
type TemplateClient struct {
	client *Client
}

// ProductInfo represents product information
type ProductInfo struct {
	Name        string   `json:"name"`
	Brand       string   `json:"brand"`
	Price       string   `json:"price"`
	Currency    string   `json:"currency"`
	Description string   `json:"description"`
	Features    []string `json:"features"`
}

// SEOAnalysis represents SEO analysis results
type SEOAnalysis struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Score       int      `json:"score"`
	Suggestions []string `json:"suggestions"`
}

// NewTemplateClient creates a new template client
func NewTemplateClient() *TemplateClient {
	return &TemplateClient{
		client: New(WithOutputFormat("json")),
	}
}

// AnalyzeProduct analyzes product information from a file
func (c *TemplateClient) AnalyzeProduct(filePath string) (*ProductInfo, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	
	prompt := fmt.Sprintf(`Please analyze the following product information and return it in JSON format:
	{
		"name": "product name",
		"brand": "brand",
		"price": "price",
		"currency": "currency unit",
		"description": "product description",
		"features": ["feature1", "feature2"]
	}
	
	Product information:
	%s`, string(content))
	
	// Use AskWithOptions with JSON output format
	response, err := c.client.AskWithOptions(prompt, map[string]interface{}{
		"output_format": "json",
	})
	if err != nil {
		return nil, err
	}
	
	// Parse JSON response
	var productInfo ProductInfo
	if err := json.Unmarshal([]byte(response), &productInfo); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}
	
	return &productInfo, nil
}

// OptimizeSEO analyzes content for SEO optimization
func (c *TemplateClient) OptimizeSEO(content string) (*SEOAnalysis, error) {
	prompt := fmt.Sprintf(`Please analyze the following content for SEO and return it in JSON format:
	{
		"title": "optimized title",
		"description": "optimized description",
		"keywords": ["keyword1", "keyword2"],
		"score": 85,
		"suggestions": ["suggestion1", "suggestion2"]
	}
	
	Content:
	%s`, content)
	
	// Use AskWithOptions with JSON output format
	response, err := c.client.AskWithOptions(prompt, map[string]interface{}{
		"output_format": "json",
	})
	if err != nil {
		return nil, err
	}
	
	// Parse JSON response
	var seoAnalysis SEOAnalysis
	if err := json.Unmarshal([]byte(response), &seoAnalysis); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}
	
	return &seoAnalysis, nil
}

// UseTemplate applies a predefined template
func (c *TemplateClient) UseTemplate(templateName string, args ...interface{}) (string, error) {
	templates := map[string]string{
		"code_review": `Please review the following code, including:
		1. Code quality analysis
		2. Potential bugs and issues
		3. Performance optimization suggestions
		4. Best practice recommendations
		
		Code file: %s`,
		
		"workflow": `Please create a detailed workflow, including:
		Task name: %s
		Objective: %s
		Steps:
		%s
		
		Please provide a detailed execution plan and considerations.`,
		
		"documentation": `Please generate documentation for the following code:
		%s
		
		Including:
		1. Function description
		2. Parameter description
		3. Return value description
		4. Usage examples`,
	}
	
	template, exists := templates[templateName]
	if !exists {
		return "", fmt.Errorf("template %s not found", templateName)
	}
	
	prompt := fmt.Sprintf(template, args...)
	
	// For code_review template, read file if argument is a file path
	if templateName == "code_review" && len(args) > 0 {
		if filePath, ok := args[0].(string); ok {
			if content, err := os.ReadFile(filePath); err == nil {
				prompt = fmt.Sprintf(templates[templateName], string(content))
			}
		}
	}
	
	return c.client.Ask(prompt)
}