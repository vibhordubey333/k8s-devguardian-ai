package output

import (
	"bytes"
	"fmt"
	"strings"
)

// CLIFormatter implements the Formatter interface for CLI output
type CLIFormatter struct{}

// NewCLIFormatter creates a new CLI formatter
func NewCLIFormatter() *CLIFormatter {
	return &CLIFormatter{}
}

// Format formats the audit result for CLI output
func (f *CLIFormatter) Format(result AuditResult) ([]byte, error) {
	var buf bytes.Buffer

	// Print header
	buf.WriteString("\n")
	buf.WriteString("🔍 KUBERNETES SECURITY AUDIT RESULTS\n")
	buf.WriteString("=====================================\n\n")

	// Print summary
	buf.WriteString(fmt.Sprintf("📊 SUMMARY: Found %d security issues\n", result.Summary.TotalFindings))
	buf.WriteString("-------------------------------------\n")
	
	// Print findings by severity
	buf.WriteString("By Severity:\n")
	for severity, count := range result.Summary.BySeverity {
		var icon string
		switch strings.ToLower(severity) {
		case "critical":
			icon = "🔴"
		case "high":
			icon = "🟠"
		case "medium":
			icon = "🟡"
		case "low":
			icon = "🟢"
		default:
			icon = "⚪"
		}
		buf.WriteString(fmt.Sprintf("  %s %s: %d\n", icon, severity, count))
	}
	
	// Print findings by resource
	buf.WriteString("\nBy Resource Type:\n")
	for resource, count := range result.Summary.ByResource {
		buf.WriteString(fmt.Sprintf("  - %s: %d\n", resource, count))
	}
	
	buf.WriteString("\n")

	// Print detailed findings
	buf.WriteString("🛡️ DETAILED FINDINGS\n")
	buf.WriteString("=====================================\n\n")

	for i, exp := range result.Explanations {
		finding := exp.Finding
		
		// Determine severity icon
		var severityIcon string
		switch strings.ToLower(finding.Severity) {
		case "critical":
			severityIcon = "🔴 CRITICAL"
		case "high":
			severityIcon = "🟠 HIGH"
		case "medium":
			severityIcon = "🟡 MEDIUM"
		case "low":
			severityIcon = "🟢 LOW"
		default:
			severityIcon = "⚪ " + finding.Severity
		}
		
		// Print finding header
		buf.WriteString(fmt.Sprintf("FINDING #%d: %s\n", i+1, severityIcon))
		buf.WriteString(fmt.Sprintf("Resource: %s/%s/%s\n", finding.Resource, finding.Namespace, finding.Name))
		buf.WriteString(fmt.Sprintf("Issue: %s\n", finding.Reason))
		buf.WriteString("-------------------------------------\n")
		
		// Print explanation
		buf.WriteString("📝 EXPLANATION:\n")
		buf.WriteString(exp.Explanation)
		buf.WriteString("\n\n")
		
		// Print remediation
		buf.WriteString("🔧 REMEDIATION:\n")
		buf.WriteString(exp.Remediation)
		buf.WriteString("\n\n")
		
		// Print references
		buf.WriteString("📚 REFERENCES:\n")
		for _, ref := range exp.References {
			buf.WriteString(fmt.Sprintf("  - %s\n", ref))
		}
		
		buf.WriteString("\n=====================================\n\n")
	}

	return buf.Bytes(), nil
}
