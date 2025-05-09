package output

import (
	"encoding/json"
)

// JSONFormatter implements the Formatter interface for JSON output
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Format formats the audit result as JSON
func (f *JSONFormatter) Format(result AuditResult) ([]byte, error) {
	return json.MarshalIndent(result, "", "  ")
}
