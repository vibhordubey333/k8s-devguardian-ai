package ai

import (
	"testing"
	
	"github.com/vibhordubey333/k8s-devguardian-ai/internal/auditor"
)

func TestSimpleExplainer_ExplainFindings(t *testing.T) {
	// Create a simple explainer
	explainer := NewSimpleExplainer()
	
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
			Resource:  "Pod",
			Namespace: "default",
			Name:      "test-pod-3",
			Reason:    "Container 'test-container' uses hostPath volume 'test-volume'",
			Severity:  "Medium",
		},
	}
	
	// Get explanations
	explanations, err := explainer.ExplainFindings(findings)
	
	// Check for errors
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Check that we got the right number of explanations
	if len(explanations) != len(findings) {
		t.Errorf("Expected %d explanations, got %d", len(findings), len(explanations))
	}
	
	// Check that each explanation has the correct finding
	for i, exp := range explanations {
		if exp.Finding.Resource != findings[i].Resource {
			t.Errorf("Expected Resource %s, got %s", findings[i].Resource, exp.Finding.Resource)
		}
		if exp.Finding.Namespace != findings[i].Namespace {
			t.Errorf("Expected Namespace %s, got %s", findings[i].Namespace, exp.Finding.Namespace)
		}
		if exp.Finding.Name != findings[i].Name {
			t.Errorf("Expected Name %s, got %s", findings[i].Name, exp.Finding.Name)
		}
		if exp.Finding.Reason != findings[i].Reason {
			t.Errorf("Expected Reason %s, got %s", findings[i].Reason, exp.Finding.Reason)
		}
		if exp.Finding.Severity != findings[i].Severity {
			t.Errorf("Expected Severity %s, got %s", findings[i].Severity, exp.Finding.Severity)
		}
		
		// Check that each explanation has non-empty explanation and remediation
		if exp.Explanation == "" {
			t.Errorf("Expected non-empty explanation for finding %d", i)
		}
		if exp.Remediation == "" {
			t.Errorf("Expected non-empty remediation for finding %d", i)
		}
		if len(exp.References) == 0 {
			t.Errorf("Expected non-empty references for finding %d", i)
		}
	}
}
