package auditor

import (
	"context"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	corev1 "k8s.io/api/core/v1"
)

type AuditFinding struct {
	Resource string
	Namespace string
	Name string
	Reason string
	Severity string
}

// AuditPodSecurity runs basic checks on pods and containers
func AuditPodSecurity(pods []corev1.Pod) []AuditFinding {
	findings := []AuditFinding{}
	for _, pod := range pods {
		for _, c := range pod.Spec.Containers {
			if c.SecurityContext != nil {
				if c.SecurityContext.RunAsUser != nil && *c.SecurityContext.RunAsUser == 0 {
					findings = append(findings, AuditFinding{
						Resource:  "Pod",
						Namespace: pod.Namespace,
						Name:      pod.Name,
						Reason:    fmt.Sprintf("Container '%s' runs as root user (uid 0)", c.Name),
						Severity:  "High",
					})
				}
				if c.SecurityContext.Privileged != nil && *c.SecurityContext.Privileged {
					findings = append(findings, AuditFinding{
						Resource:  "Pod",
						Namespace: pod.Namespace,
						Name:      pod.Name,
						Reason:    fmt.Sprintf("Container '%s' is privileged", c.Name),
						Severity:  "Critical",
					})
				}
			}
			// Check for hostPath mounts
			for _, v := range pod.Spec.Volumes {
				if v.HostPath != nil {
					findings = append(findings, AuditFinding{
						Resource:  "Pod",
						Namespace: pod.Namespace,
						Name:      pod.Name,
						Reason:    fmt.Sprintf("Container '%s' uses hostPath volume '%s'", c.Name, v.Name),
						Severity:  "Medium",
					})
				}
			}
		}
	}
	return findings
}

// EvaluateWithOPA evaluates a pod against a given rego policy
func EvaluateWithOPA(pod corev1.Pod, regoModule string) ([]AuditFinding, error) {
	ctx := context.Background()
	query := rego.New(
		rego.Query("data.k8s.pod.violation"),
		rego.Module("policy.rego", regoModule),
		rego.Input(pod),
	)
	results, err := query.Eval(ctx)
	if err != nil {
		return nil, err
	}
	findings := []AuditFinding{}
	for _, result := range results {
		for _, exp := range result.Expressions {
			if arr, ok := exp.Value.([]interface{}); ok {
				for _, f := range arr {
					if m, ok := f.(map[string]interface{}); ok {
						findings = append(findings, AuditFinding{
							Resource:  "Pod",
							Namespace: pod.Namespace,
							Name:      pod.Name,
							Reason:    fmt.Sprintf("OPA: %v", m["reason"]),
							Severity:  fmt.Sprintf("%v", m["severity"]),
						})
					}
				}
			}
		}
	}
	return findings, nil
}
