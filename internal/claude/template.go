package claude

import (
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
	
	prompt := fmt.Sprintf(`请分析以下产品信息，并以JSON格式返回：
	{
		"name": "产品名称",
		"brand": "品牌",
		"price": "价格",
		"currency": "货币单位",
		"description": "产品描述",
		"features": ["特征1", "特征2"]
	}
	
	产品信息：
	%s`, string(content))
	
	var productInfo ProductInfo
	if err := c.client.AskJSON(prompt, &productInfo); err != nil {
		return nil, err
	}
	
	return &productInfo, nil
}

// OptimizeSEO analyzes content for SEO optimization
func (c *TemplateClient) OptimizeSEO(content string) (*SEOAnalysis, error) {
	prompt := fmt.Sprintf(`请对以下内容进行SEO分析，并以JSON格式返回：
	{
		"title": "优化后的标题",
		"description": "优化后的描述",
		"keywords": ["关键词1", "关键词2"],
		"score": 85,
		"suggestions": ["建议1", "建议2"]
	}
	
	内容：
	%s`, content)
	
	var seoAnalysis SEOAnalysis
	if err := c.client.AskJSON(prompt, &seoAnalysis); err != nil {
		return nil, err
	}
	
	return &seoAnalysis, nil
}

// UseTemplate applies a predefined template
func (c *TemplateClient) UseTemplate(templateName string, args ...interface{}) (string, error) {
	templates := map[string]string{
		"code_review": `请对以下代码进行审查，包括：
		1. 代码质量分析
		2. 潜在的bug和问题
		3. 性能优化建议
		4. 最佳实践建议
		
		代码文件：%s`,
		
		"workflow": `请创建一个详细的工作流程，包括：
		任务名称：%s
		目标：%s
		步骤：
		%s
		
		请提供详细的执行计划和注意事项。`,
		
		"documentation": `请为以下代码生成文档：
		%s
		
		包括：
		1. 功能说明
		2. 参数说明
		3. 返回值说明
		4. 使用示例`,
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