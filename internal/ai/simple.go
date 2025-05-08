package ai

import (
	"github.com/vibhordubey333/k8s-devguardian-ai/internal/auditor"
)

// SimpleExplainer implements the Explainer interface without using AI
type SimpleExplainer struct{}

// NewSimpleExplainer creates a new simple explainer
func NewSimpleExplainer() *SimpleExplainer {
	return &SimpleExplainer{}
}

// IsAvailable always returns true for the simple explainer
func (e *SimpleExplainer) IsAvailable() (bool, error) {
	return true, nil
}

// ExplainFindings provides basic explanations based on predefined templates
func (e *SimpleExplainer) ExplainFindings(findings []auditor.AuditFinding) ([]FindingExplanation, error) {
	explanations := make([]FindingExplanation, 0, len(findings))

	for _, finding := range findings {
		var explanation, remediation string
		references := []string{"https://kubernetes.io/docs/concepts/security/"}

		// Provide basic explanations based on the finding reason
		switch {
		case contains(finding.Reason, "privileged"):
			explanation = "Privileged containers have access to all devices on the host, which can lead to security vulnerabilities if compromised."
			remediation = "Remove the privileged flag from the container's securityContext or use a more restrictive security context."
			references = append(references, "https://kubernetes.io/docs/tasks/configure-pod-container/security-context/")

		case contains(finding.Reason, "root"):
			explanation = "Running containers as root (uid 0) gives them elevated permissions, which is a security risk."
			remediation = "Set runAsUser in the container's securityContext to a non-zero value."
			references = append(references, "https://kubernetes.io/docs/concepts/security/pod-security-standards/")

		case contains(finding.Reason, "hostPath"):
			explanation = "hostPath volumes allow pods to access files on the host, which can lead to privilege escalation."
			remediation = "Avoid using hostPath volumes. Consider using more secure volume types like emptyDir, configMap, or PersistentVolumeClaims."
			references = append(references, "https://kubernetes.io/docs/concepts/storage/volumes/")

		default:
			explanation = "This is a security issue that could potentially compromise your Kubernetes cluster."
			remediation = "Review the resource configuration and apply security best practices."
		}

		explanations = append(explanations, FindingExplanation{
			Finding:     finding,
			Explanation: explanation,
			Remediation: remediation,
			References:  references,
		})
	}

	return explanations, nil
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	s, substr = toLower(s), toLower(substr)
	return s != "" && substr != "" && indexOf(s, substr) >= 0
}

// toLower converts a string to lowercase
func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

// indexOf returns the index of the first instance of substr in s, or -1 if substr is not present
func indexOf(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
