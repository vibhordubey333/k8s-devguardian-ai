package scanner

import (
	"context"
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"k8s.io/client-go/rest"
)

func ScanCluster() error {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := loadKubeConfig(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list pods: %w", err)
	}

	fmt.Printf("🔍 Found %d pods in the cluster:\n", len(pods.Items))
	for _, pod := range pods.Items {
		fmt.Printf(" - %s/%s\n", pod.Namespace, pod.Name)
	}

	return nil
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
