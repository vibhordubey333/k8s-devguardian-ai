package ai

import (
	"github.com/vibhordubey333/k8s-devguardian-ai/internal/auditor"
)

// ExplainerConfig holds configuration for AI explainers
type ExplainerConfig struct {
	Provider  string // "openai" or "ollama"
	APIKey    string // API key for OpenAI
	ModelName string // Model name to use
	OllamaURL string // URL for Ollama server
}

// Explainer is the interface that all AI explainers must implement
type Explainer interface {
	// ExplainFindings takes audit findings and returns explanations with remediation suggestions
	ExplainFindings(findings []auditor.AuditFinding) ([]FindingExplanation, error)
	
	// IsAvailable checks if the AI service is available
	IsAvailable() (bool, error)
}

// FindingExplanation represents an explanation for a security finding
type FindingExplanation struct {
	Finding     auditor.AuditFinding // The original finding
	Explanation string               // Human-readable explanation
	Remediation string               // Suggested remediation steps
	References  []string             // References to security best practices or documentation
}

// NewExplainer creates a new explainer based on the provided configuration
func NewExplainer(config ExplainerConfig) (Explainer, error) {
	switch config.Provider {
	case "openai":
		return NewOpenAIExplainer(config)
	case "ollama":
		return NewOllamaExplainer(config)
	default:
		// Default to a simple explainer that doesn't use AI
		return NewSimpleExplainer(), nil
	}
}
