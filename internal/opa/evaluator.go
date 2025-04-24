package opa

import (
	"context"
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/rego"
	"gopkg.in/yaml.v2"
)

func Evaluate(resourceYAML []byte, policyPath string) ([]string, error) {
	var input map[string]interface{}
	err := yaml.Unmarshal(resourceYAML, &input)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	policy, err := os.ReadFile(policyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy: %w", err)
	}

	ctx := context.Background()
	query, err := rego.New(
		rego.Query("data.devguardian.k8s.deny"),
		rego.Module("policy.rego", string(policy)),
	).PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("rego compile error: %w", err)
	}

	results, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, fmt.Errorf("rego eval error: %w", err)
	}

	var reasons []string
	for _, result := range results {
		for _, expr := range result.Expressions {
			for _, val := range expr.Value.([]interface{}) {
				reasons = append(reasons, val.(string))
			}
		}
	}
	return reasons, nil
}
