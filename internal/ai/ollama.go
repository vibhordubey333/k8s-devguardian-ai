package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	
	"github.com/vibhordubey333/k8s-devguardian-ai/internal/auditor"
)

// OllamaExplainer implements the Explainer interface using Ollama
type OllamaExplainer struct {
	ollamaURL string
	modelName string
	client    *http.Client
}

// OllamaRequest represents a request to the Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// OllamaResponse represents a response from the Ollama API
type OllamaResponse struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

// NewOllamaExplainer creates a new Ollama explainer
func NewOllamaExplainer(config ExplainerConfig) (*OllamaExplainer, error) {
	ollamaURL := config.OllamaURL
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434" // Default Ollama URL
	}

	modelName := config.ModelName
	if modelName == "" {
		modelName = "llama2" // Default model
	}

	return &OllamaExplainer{
		ollamaURL: ollamaURL,
		modelName: modelName,
		client: &http.Client{
			Timeout: 60 * time.Second, // Longer timeout for local models
		},
	}, nil
}

// IsAvailable checks if the Ollama API is available
func (e *OllamaExplainer) IsAvailable() (bool, error) {
	// Simple check by making a minimal request
	req := OllamaRequest{
		Model:  e.modelName,
		Prompt: "Hello",
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return false, err
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/generate", e.ollamaURL), bytes.NewBuffer(reqBody))
	if err != nil {
		return false, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(httpReq)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, fmt.Errorf("Ollama API returned status code %d", resp.StatusCode)
}

// ExplainFindings takes audit findings and returns explanations with remediation suggestions
func (e *OllamaExplainer) ExplainFindings(findings []auditor.AuditFinding) ([]FindingExplanation, error) {
	explanations := make([]FindingExplanation, 0, len(findings))

	for _, finding := range findings {
		systemPrompt := "You are a Kubernetes security expert. Provide clear explanations and practical remediation steps for Kubernetes security issues."
		userPrompt := generatePrompt(finding)
		
		fullPrompt := fmt.Sprintf("%s\n\n%s", systemPrompt, userPrompt)
		
		req := OllamaRequest{
			Model:  e.modelName,
			Prompt: fullPrompt,
		}

		reqBody, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}

		httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/generate", e.ollamaURL), bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, err
		}

		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := e.client.Do(httpReq)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var ollamaResp OllamaResponse
		if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
			return nil, err
		}

		if ollamaResp.Error != "" {
			return nil, fmt.Errorf("Ollama API error: %s", ollamaResp.Error)
		}

		explanation, remediation, references := parseAIResponse(ollamaResp.Response)

		explanations = append(explanations, FindingExplanation{
			Finding:     finding,
			Explanation: explanation,
			Remediation: remediation,
			References:  references,
		})
	}

	return explanations, nil
}
