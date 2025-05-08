package output

import (
	"github.com/vibhordubey333/k8s-devguardian-ai/internal/ai"
	"github.com/vibhordubey333/k8s-devguardian-ai/internal/auditor"
)

// Format represents the output format
type Format string

const (
	// FormatCLI represents the CLI output format
	FormatCLI Format = "cli"
	// FormatJSON represents the JSON output format
	FormatJSON Format = "json"
	// FormatHTML represents the HTML output format
	FormatHTML Format = "html"
)

// AuditResult represents the result of an audit
type AuditResult struct {
	Findings     []auditor.AuditFinding    // The original findings
	Explanations []ai.FindingExplanation   // The explanations for the findings
	Summary      AuditSummary              // Summary of the audit
}

// AuditSummary represents a summary of the audit
type AuditSummary struct {
	TotalFindings int            // Total number of findings
	BySeverity    map[string]int // Number of findings by severity
	ByResource    map[string]int // Number of findings by resource type
}

// Formatter is the interface that all output formatters must implement
type Formatter interface {
	// Format formats the audit result
	Format(result AuditResult) ([]byte, error)
}

// NewFormatter creates a new formatter based on the provided format
func NewFormatter(format Format) Formatter {
	switch format {
	case FormatJSON:
		return NewJSONFormatter()
	case FormatHTML:
		return NewHTMLFormatter()
	default:
		return NewCLIFormatter()
	}
}

// GenerateSummary generates a summary of the audit result
func GenerateSummary(findings []auditor.AuditFinding) AuditSummary {
	bySeverity := make(map[string]int)
	byResource := make(map[string]int)

	for _, finding := range findings {
		bySeverity[finding.Severity]++
		byResource[finding.Resource]++
	}

	return AuditSummary{
		TotalFindings: len(findings),
		BySeverity:    bySeverity,
		ByResource:    byResource,
	}
}
