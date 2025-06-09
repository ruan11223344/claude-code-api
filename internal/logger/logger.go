package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// Initialize initializes the logger with configuration from environment variables
func Initialize() {
	Log = logrus.New()
	
	// Set log level from environment
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	Log.SetLevel(level)
	
	// Set output format
	if os.Getenv("LOG_FORMAT") == "json" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		// Use custom text formatter for better readability
		Log.SetFormatter(&CustomTextFormatter{
			TimestampFormat: "15:04:05",
			FullTimestamp:   true,
		})
	}
	
	// Configure output
	logToFile := os.Getenv("LOG_TO_FILE")
	if logToFile == "true" || logToFile == "1" {
		// Get max log files (default: 7 days)
		maxFiles := 7
		if maxFilesEnv := os.Getenv("LOG_MAX_FILES"); maxFilesEnv != "" {
			if parsed, err := fmt.Sscanf(maxFilesEnv, "%d", &maxFiles); err == nil && parsed == 1 {
				// Successfully parsed
			}
		}
		
		// Create multi-writer with rotating file
		writer, err := MultiWriter(os.Stdout, "logs", maxFiles)
		if err != nil {
			Log.WithError(err).Error("Failed to create rotating log writer")
			Log.SetOutput(os.Stdout)
		} else {
			Log.SetOutput(writer)
			Log.WithFields(logrus.Fields{
				"directory": "logs",
				"max_files": maxFiles,
			}).Info("Rotating file logging enabled")
		}
	} else {
		// Only log to stdout
		Log.SetOutput(os.Stdout)
	}
}

// init creates a basic logger for use before Initialize is called
func init() {
	Log = logrus.New()
	Log.SetLevel(logrus.InfoLevel)
	Log.SetFormatter(&CustomTextFormatter{
		TimestampFormat: "15:04:05",
		FullTimestamp:   true,
	})
	Log.SetOutput(os.Stdout)
}

// CustomTextFormatter formats logs in a clean, readable format
type CustomTextFormatter struct {
	TimestampFormat string
	FullTimestamp   bool
}

func (f *CustomTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)
	
	// Color codes for different levels
	var levelColor string
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = "\033[37m" // White
	case logrus.InfoLevel:
		levelColor = "\033[36m" // Cyan
	case logrus.WarnLevel:
		levelColor = "\033[33m" // Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = "\033[31m" // Red
	default:
		levelColor = "\033[0m" // Default
	}
	
	// Format the log entry
	var logLine string
	level := strings.ToUpper(entry.Level.String())
	
	// Check if we have fields
	if len(entry.Data) > 0 {
		// Format with fields
		fields := ""
		for k, v := range entry.Data {
			fields += " " + k + "=" + formatValue(v)
		}
		logLine = levelColor + "[" + timestamp + "] " + level + "\033[0m " + entry.Message + fields + "\n"
	} else {
		// Simple format without fields
		logLine = levelColor + "[" + timestamp + "] " + level + "\033[0m " + entry.Message + "\n"
	}
	
	return []byte(logLine), nil
}

func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		// Truncate long strings
		if len(val) > 100 {
			return "\"" + val[:97] + "...\""
		}
		return "\"" + val + "\""
	default:
		return fmt.Sprintf("%v", val)
	}
}

// SanitizeRequest removes sensitive information from request logs
func SanitizeRequest(body string) string {
	// List of sensitive fields to mask
	sensitiveFields := []string{"api_key", "apiKey", "password", "token", "secret", "authorization"}
	
	sanitized := body
	for _, field := range sensitiveFields {
		// Simple regex-like replacement
		if strings.Contains(strings.ToLower(sanitized), strings.ToLower(field)) {
			// This is a simple implementation - in production, use proper JSON parsing
			sanitized = maskSensitiveField(sanitized, field)
		}
	}
	
	// Truncate very long bodies
	if len(sanitized) > 1000 {
		sanitized = sanitized[:997] + "..."
	}
	
	return sanitized
}

func maskSensitiveField(content, field string) string {
	// This is a simplified version - in production, parse JSON properly
	// For now, just replace common patterns
	lowerContent := strings.ToLower(content)
	lowerField := strings.ToLower(field)
	
	idx := strings.Index(lowerContent, lowerField)
	if idx != -1 {
		// Find the value after the field
		valueStart := idx + len(field)
		// Skip to the actual value (past : and quotes)
		for valueStart < len(content) && (content[valueStart] == ':' || content[valueStart] == ' ' || content[valueStart] == '"') {
			valueStart++
		}
		
		// Find the end of the value
		valueEnd := valueStart
		inQuotes := valueStart > 0 && content[valueStart-1] == '"'
		for valueEnd < len(content) {
			if inQuotes && content[valueEnd] == '"' {
				break
			} else if !inQuotes && (content[valueEnd] == ',' || content[valueEnd] == '}' || content[valueEnd] == ' ') {
				break
			}
			valueEnd++
		}
		
		// Replace the value with asterisks
		if valueEnd > valueStart {
			masked := content[:valueStart] + "********" + content[valueEnd:]
			return masked
		}
	}
	
	return content
}