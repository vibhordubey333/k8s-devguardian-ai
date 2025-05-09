package output

import (
	"testing"
	
	"github.com/vibhordubey333/k8s-devguardian-ai/internal/ai"
	"github.com/vibhordubey333/k8s-devguardian-ai/internal/auditor"
)

func TestNewFormatter(t *testing.T) {
	// Test CLI formatter
	formatter := NewFormatter(FormatCLI)
	if _, ok := formatter.(*CLIFormatter); !ok {
		t.Errorf("Expected CLIFormatter, got %T", formatter)
	}
	
	// Test JSON formatter
	formatter = NewFormatter(FormatJSON)
	if _, ok := formatter.(*JSONFormatter); !ok {
		t.Errorf("Expected JSONFormatter, got %T", formatter)
	}
	
	// Test HTML formatter
	formatter = NewFormatter(FormatHTML)
	if _, ok := formatter.(*HTMLFormatter); !ok {
		t.Errorf("Expected HTMLFormatter, got %T", formatter)
	}
	
	// Test default formatter
	formatter = NewFormatter("invalid")
	if _, ok := formatter.(*CLIFormatter); !ok {
		t.Errorf("Expected CLIFormatter as default, got %T", formatter)
	}
}

func TestGenerateSummary(t *testing.T) {
	// Create some test findings
	findings := []auditor.AuditFinding{
		{
			Resource:  "Pod",
			Namespace: "default",
			Name:      "test-pod",
			Reason:    "Container 'test-container' is privileged",
			Severity:  "Critical",
		},
		{
			Resource:  "Pod",
			Namespace: "default",
			Name:      "test-pod-2",
			Reason:    "Container 'test-container' runs as root user (uid 0)",
			Severity:  "High",
		},
		{
			Resource:  "Service",
			Namespace: "default",
			Name:      "test-service",
			Reason:    "Service uses NodePort which exposes ports on all nodes",
			Severity:  "Medium",
		},
	}
	
	// Generate summary
	summary := GenerateSummary(findings)
	
	// Check total findings
	if summary.TotalFindings != len(findings) {
		t.Errorf("Expected %d total findings, got %d", len(findings), summary.TotalFindings)
	}
	
	// Check findings by severity
	if summary.BySeverity["Critical"] != 1 {
		t.Errorf("Expected 1 Critical finding, got %d", summary.BySeverity["Critical"])
	}
	if summary.BySeverity["High"] != 1 {
		t.Errorf("Expected 1 High finding, got %d", summary.BySeverity["High"])
	}
	if summary.BySeverity["Medium"] != 1 {
		t.Errorf("Expected 1 Medium finding, got %d", summary.BySeverity["Medium"])
	}
	
	// Check findings by resource
	if summary.ByResource["Pod"] != 2 {
		t.Errorf("Expected 2 Pod findings, got %d", summary.ByResource["Pod"])
	}
	if summary.ByResource["Service"] != 1 {
		t.Errorf("Expected 1 Service finding, got %d", summary.ByResource["Service"])
	}
}

func TestJSONFormatter_Format(t *testing.T) {
	// Create a JSON formatter
	formatter := NewJSONFormatter()
	
	// Create a test audit result
	result := AuditResult{
		Findings: []auditor.AuditFinding{
			{
				Resource:  "Pod",
				Namespace: "default",
				Name:      "test-pod",
				Reason:    "Container 'test-container' is privileged",
				Severity:  "Critical",
			},
		},
		Explanations: []ai.FindingExplanation{
			{
				Finding: auditor.AuditFinding{
					Resource:  "Pod",
					Namespace: "default",
					Name:      "test-pod",
					Reason:    "Container 'test-container' is privileged",
					Severity:  "Critical",
				},
				Explanation: "Test explanation",
				Remediation: "Test remediation",
				References:  []string{"https://example.com"},
			},
		},
		Summary: AuditSummary{
			TotalFindings: 1,
			BySeverity:    map[string]int{"Critical": 1},
			ByResource:    map[string]int{"Pod": 1},
		},
	}
	
	// Format the result
	bytes, err := formatter.Format(result)
	
	// Check for errors
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Check that we got some output
	if len(bytes) == 0 {
		t.Errorf("Expected non-empty output")
	}
}
