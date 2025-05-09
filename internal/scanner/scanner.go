package scanner

import (
	"context"
	"fmt"
	"path/filepath"
	"os"

	"github.com/vibhordubey333/k8s-devguardian-ai/internal/auditor"
	"github.com/vibhordubey333/k8s-devguardian-ai/internal/opa"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	corev1 "k8s.io/api/core/v1"
	"gopkg.in/yaml.v2"
)

// ScanCluster scans the Kubernetes cluster and returns findings
func ScanCluster() ([]auditor.AuditFinding, error) {
	var findings []auditor.AuditFinding

	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := loadKubeConfig(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	// Scan pods
	fmt.Println("Scanning pods...")
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	// Audit pods using built-in checks
	podFindings := auditor.AuditPodSecurity(pods.Items)
	findings = append(findings, podFindings...)

	// Scan pods with OPA policies
	for _, pod := range pods.Items {
		podYAML, err := yaml.Marshal(pod)
		if err != nil {
			fmt.Printf("Warning: Failed to marshal pod %s/%s: %v\n", pod.Namespace, pod.Name, err)
			continue
		}

		// Use OPA to evaluate the pod against policies
		// This is a placeholder - in a real implementation, you would load policies from a directory
		policyPath := "internal/policies/privileged_pod.rego"
		if _, err := os.Stat(policyPath); err == nil {
			reasons, err := opa.Evaluate(podYAML, policyPath)
			if err != nil {
				fmt.Printf("Warning: Failed to evaluate pod %s/%s with OPA: %v\n", pod.Namespace, pod.Name, err)
				continue
			}

			for _, reason := range reasons {
				findings = append(findings, auditor.AuditFinding{
					Resource:  "Pod",
					Namespace: pod.Namespace,
					Name:      pod.Name,
					Reason:    reason,
					Severity:  "High", // Default severity, could be determined by policy
				})
			}
		}
	}

	// Scan services for potential security issues
	fmt.Println("Scanning services...")
	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %v", err)
	}

	// Check for NodePort and LoadBalancer services
	for _, svc := range services.Items {
		if svc.Spec.Type == corev1.ServiceTypeNodePort {
			findings = append(findings, auditor.AuditFinding{
				Resource:  "Service",
				Namespace: svc.Namespace,
				Name:      svc.Name,
				Reason:    "Service uses NodePort which exposes ports on all nodes",
				Severity:  "Medium",
			})
		}
		if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
			findings = append(findings, auditor.AuditFinding{
				Resource:  "Service",
				Namespace: svc.Namespace,
				Name:      svc.Name,
				Reason:    "Service uses LoadBalancer which may expose the service to the internet",
				Severity:  "Medium",
			})
		}
	}

	// Scan RBAC roles for excessive permissions
	fmt.Println("Scanning RBAC roles...")
	roles, err := clientset.RbacV1().Roles("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %v", err)
	}

	for _, role := range roles.Items {
		for _, rule := range role.Rules {
			// Check for wildcard resources or verbs
			hasWildcardResources := false
			hasWildcardVerbs := false

			for _, resource := range rule.Resources {
				if resource == "*" {
					hasWildcardResources = true
					break
				}
			}

			for _, verb := range rule.Verbs {
				if verb == "*" {
					hasWildcardVerbs = true
					break
				}
			}

			if hasWildcardResources && hasWildcardVerbs {
				findings = append(findings, auditor.AuditFinding{
					Resource:  "Role",
					Namespace: role.Namespace,
					Name:      role.Name,
					Reason:    "Role has wildcard resources and verbs which grants excessive permissions",
					Severity:  "High",
				})
				break
			}
		}
	}

	// Scan namespaces for PodSecurity settings
	fmt.Println("Scanning namespaces...")
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %v", err)
	}

	for _, ns := range namespaces.Items {
		enforce := ns.Labels["pod-security.kubernetes.io/enforce"]
		if enforce == "" || enforce == "privileged" {
			findings = append(findings, auditor.AuditFinding{
				Resource:  "Namespace",
				Namespace: ns.Name,
				Name:      ns.Name,
				Reason:    "Namespace does not enforce PodSecurity standards or uses 'privileged' level",
				Severity:  "High",
			})
		}
	}

	return findings, nil
}

func loadKubeConfig(kubeConfigPath string) (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: kubeconfig file does not exist: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "Error: failed to load kubeconfig: %v\n", err)
		}
		return nil, err
	}
	return config, nil
}
