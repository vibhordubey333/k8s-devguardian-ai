package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	
	"github.com/vibhordubey333/k8s-devguardian-ai/internal/auditor"
)

// OpenAIExplainer implements the Explainer interface using OpenAI API
type OpenAIExplainer struct {
	apiKey    string
	modelName string
	client    *http.Client
}

// OpenAIRequest represents a request to the OpenAI API
type OpenAIRequest struct {
	Model    string        `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

// OpenAIMessage represents a message in the OpenAI API request
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents a response from the OpenAI API
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// NewOpenAIExplainer creates a new OpenAI explainer
func NewOpenAIExplainer(config ExplainerConfig) (*OpenAIExplainer, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	modelName := config.ModelName
	if modelName == "" {
		modelName = "gpt-3.5-turbo" // Default model
	}

	return &OpenAIExplainer{
		apiKey:    config.APIKey,
		modelName: modelName,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// IsAvailable checks if the OpenAI API is available
func (e *OpenAIExplainer) IsAvailable() (bool, error) {
	// Simple check by making a minimal request
	req := OpenAIRequest{
		Model: e.modelName,
		Messages: []OpenAIMessage{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return false, err
	}

	httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return false, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+e.apiKey)

	resp, err := e.client.Do(httpReq)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, fmt.Errorf("OpenAI API returned status code %d", resp.StatusCode)
}

// ExplainFindings takes audit findings and returns explanations with remediation suggestions
func (e *OpenAIExplainer) ExplainFindings(findings []auditor.AuditFinding) ([]FindingExplanation, error) {
	explanations := make([]FindingExplanation, 0, len(findings))

	for _, finding := range findings {
		prompt := generatePrompt(finding)
		
		req := OpenAIRequest{
			Model: e.modelName,
			Messages: []OpenAIMessage{
				{
					Role:    "system",
					Content: "You are a Kubernetes security expert. Provide clear explanations and practical remediation steps for Kubernetes security issues.",
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
		}

		reqBody, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}

		httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, err
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+e.apiKey)

		resp, err := e.client.Do(httpReq)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var openAIResp OpenAIResponse
		if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
			return nil, err
		}

		if openAIResp.Error != nil {
			return nil, fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
		}

		if len(openAIResp.Choices) == 0 {
			return nil, fmt.Errorf("no response from OpenAI API")
		}

		content := openAIResp.Choices[0].Message.Content
		explanation, remediation, references := parseAIResponse(content)

		explanations = append(explanations, FindingExplanation{
			Finding:     finding,
			Explanation: explanation,
			Remediation: remediation,
			References:  references,
		})
	}

	return explanations, nil
}

// generatePrompt creates a prompt for the OpenAI API based on the finding
func generatePrompt(finding auditor.AuditFinding) string {
	return fmt.Sprintf(`
Explain the following Kubernetes security issue and provide remediation steps:

Resource Type: %s
Namespace: %s
Name: %s
Issue: %s
Severity: %s

Format your response in the following structure:
EXPLANATION: [Detailed explanation of why this is a security issue]
REMEDIATION: [Step-by-step remediation instructions]
REFERENCES: [Comma-separated list of references to security best practices or documentation]
`, finding.Resource, finding.Namespace, finding.Name, finding.Reason, finding.Severity)
}

// parseAIResponse parses the AI response into explanation, remediation, and references
func parseAIResponse(response string) (explanation, remediation string, references []string) {
	// This is a simple parser that assumes the AI follows the requested format
	// In a production environment, you might want to use a more robust parsing method
	
	// Default values in case parsing fails
	explanation = response
	remediation = "No specific remediation steps provided."
	references = []string{"https://kubernetes.io/docs/concepts/security/"}

	// Try to parse the structured response
	lines := bytes.Split([]byte(response), []byte("\n"))
	
	section := ""
	var explanationLines, remediationLines []string
	var refs []string

	for _, line := range lines {
		lineStr := string(line)
		
		if len(lineStr) == 0 {
			continue
		}

		if bytes.HasPrefix(line, []byte("EXPLANATION:")) {
			section = "explanation"
			explanationLines = append(explanationLines, lineStr[len("EXPLANATION:"):])
		} else if bytes.HasPrefix(line, []byte("REMEDIATION:")) {
			section = "remediation"
			remediationLines = append(remediationLines, lineStr[len("REMEDIATION:"):])
		} else if bytes.HasPrefix(line, []byte("REFERENCES:")) {
			section = "references"
			refStr := lineStr[len("REFERENCES:"):]
			for _, ref := range bytes.Split([]byte(refStr), []byte(",")) {
				refs = append(refs, string(bytes.TrimSpace(ref)))
			}
		} else {
			switch section {
			case "explanation":
				explanationLines = append(explanationLines, lineStr)
			case "remediation":
				remediationLines = append(remediationLines, lineStr)
			case "references":
				// Additional references lines are not expected in the current format
			}
		}
	}

	if len(explanationLines) > 0 {
		explanation = string(bytes.Join(bytes.Fields([]byte(fmt.Sprintf("%s", explanationLines))), []byte(" ")))
	}
	
	if len(remediationLines) > 0 {
		remediation = string(bytes.Join(bytes.Fields([]byte(fmt.Sprintf("%s", remediationLines))), []byte(" ")))
	}
	
	if len(refs) > 0 {
		references = refs
	}

	return explanation, remediation, references
}
