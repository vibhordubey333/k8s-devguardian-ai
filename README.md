# K8s DevGuardian AI

![DevGuardian Logo](https://github.com/user-attachments/assets/c9d6e84e-a078-4dec-8fd0-f380840ce977)

![Architecture](https://github.com/user-attachments/assets/2b5f6a82-edb0-4a88-8105-ca8218405bbd)

## Overview

K8s DevGuardian AI is an open-source, AI-powered Kubernetes security auditing CLI tool designed to help DevSecOps teams identify and remediate misconfigurations in real time.

### Key Features

- **Kubernetes Cluster Scanning**: Connect to live Kubernetes clusters and scan resources for security issues
- **Policy-Based Evaluation**: Leverage Open Policy Agent (OPA) and Rego policies to detect misconfigurations
- **AI-Powered Explanations**: Get detailed explanations and remediation suggestions using AI (OpenAI/Ollama)
- **Multiple Output Formats**: View results in CLI, JSON, or HTML formats

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/k8s-devguardian-ai.git
cd k8s-devguardian-ai

# Build the binary
go build -o devguardian

# Move to a directory in your PATH (optional)
mv devguardian /usr/local/bin/
```

## Usage

### Basic Commands

```bash
# Display help information
devguardian --help
devguardian -h

# Display version information
devguardian version

# Display help for a specific command
devguardian audit --help
devguardian audit -h
```

### Audit Commands

```bash
# Run a basic audit on your current Kubernetes context (default output: CLI)
devguardian audit

# Save audit results to a file
devguardian audit --file report.txt
devguardian audit -f report.txt
```

### Output Format Options

```bash
# Output in CLI format (default)
devguardian audit --output cli
devguardian audit -o cli

# Output in JSON format
devguardian audit --output json
devguardian audit -o json

# Output in HTML format
devguardian audit --output html
devguardian audit -o html

# Save JSON output to a file
devguardian audit --output json --file report.json
devguardian audit -o json -f report.json

# Save HTML output to a file
devguardian audit --output html --file report.html
devguardian audit -o html -f report.html
```

### AI Integration Options

```bash
# Use OpenAI for explanations
devguardian audit --ai-provider openai --api-key YOUR_API_KEY
devguardian audit -a openai -k YOUR_API_KEY

# Use OpenAI with a specific model
devguardian audit --ai-provider openai --api-key YOUR_API_KEY --model gpt-4
devguardian audit -a openai -k YOUR_API_KEY -m gpt-4

# Use local Ollama for explanations (default URL: http://localhost:11434)
devguardian audit --ai-provider ollama
devguardian audit -a ollama

# Use Ollama with a custom URL
devguardian audit --ai-provider ollama --ollama-url http://custom-ollama-server:11434
devguardian audit -a ollama -u http://custom-ollama-server:11434

# Use Ollama with a specific model
devguardian audit --ai-provider ollama --model llama2
devguardian audit -a ollama -m llama2
```

### Command Flags Reference

The `audit` command supports the following flags:

| Flag | Short | Description | Default |
|------|-------|-------------|--------|
| `--output` | `-o` | Output format (cli, json, html) | `cli` |
| `--ai-provider` | `-a` | AI provider (openai, ollama) | None (uses simple explainer) |
| `--api-key` | `-k` | API key for OpenAI | None |
| `--model` | `-m` | Model name to use | OpenAI: `gpt-3.5-turbo`, Ollama: `llama2` |
| `--ollama-url` | `-u` | URL for Ollama server | `http://localhost:11434` |
| `--file` | `-f` | Output file path | None (prints to stdout) |
| `--help` | `-h` | Help for audit command | N/A |

### Combined Command Examples

```bash
# Complete audit with OpenAI and JSON output saved to a file
devguardian audit --ai-provider openai --api-key YOUR_API_KEY --output json --file report.json
devguardian audit -a openai -k YOUR_API_KEY -o json -f report.json

# Complete audit with Ollama and HTML output saved to a file
devguardian audit --ai-provider ollama --model llama2 --output html --file report.html
devguardian audit -a ollama -m llama2 -o html -f report.html

# Quick audit with simple explanations and CLI output
devguardian audit
```

## Output Examples

### CLI Output

The CLI output provides a human-readable summary of findings with colorful indicators for severity levels:

```
🔍 KUBERNETES SECURITY AUDIT RESULTS
=====================================

📊 SUMMARY: Found 24 security issues
-------------------------------------
By Severity:
  🔴 Critical: 1
  🟠 High: 4
  🟡 Medium: 19

By Resource Type:
  - Pod: 20
  - Namespace: 4

🛡️ DETAILED FINDINGS
=====================================

FINDING #1: 🔴 CRITICAL
Resource: Pod/kube-system/kube-proxy-tfrnx
Issue: Container 'kube-proxy' is privileged
-------------------------------------
📝 EXPLANATION:
Privileged containers have access to all devices on the host, which can lead to security vulnerabilities if compromised.

🔧 REMEDIATION:
Remove the privileged flag from the container's securityContext or use a more restrictive security context.

📚 REFERENCES:
  - https://kubernetes.io/docs/concepts/security/
  - https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
```

### JSON Output

The JSON output provides a structured format that can be easily parsed by other tools:

```json
{
  "Findings": [
    {
      "Resource": "Pod",
      "Namespace": "kube-system",
      "Name": "kube-proxy-tfrnx",
      "Reason": "Container 'kube-proxy' is privileged",
      "Severity": "Critical"
    }
  ],
  "Explanations": [
    {
      "Finding": {
        "Resource": "Pod",
        "Namespace": "kube-system",
        "Name": "kube-proxy-tfrnx",
        "Reason": "Container 'kube-proxy' is privileged",
        "Severity": "Critical"
      },
      "Explanation": "Privileged containers have access to all devices on the host...",
      "Remediation": "Remove the privileged flag from the container's securityContext...",
      "References": [
        "https://kubernetes.io/docs/concepts/security/",
        "https://kubernetes.io/docs/tasks/configure-pod-container/security-context/"
      ]
    }
  ],
  "Summary": {
    "TotalFindings": 24,
    "BySeverity": {
      "Critical": 1,
      "High": 4,
      "Medium": 19
    },
    "ByResource": {
      "Pod": 20,
      "Namespace": 4
    }
  }
}
```

### HTML Output

The HTML output provides a visually appealing report that can be viewed in a web browser, with color-coded severity levels and expandable sections for detailed information.

## Architecture

K8s DevGuardian AI consists of several components:

1. **CLI Interface**: Built with Cobra for a user-friendly command-line experience
2. **Kubernetes Client**: Connects to your cluster using the standard kubeconfig
3. **Policy Engine**: Uses OPA to evaluate resources against security policies
4. **AI Explainer**: Leverages LLMs to provide human-readable explanations and fix suggestions
5. **Output Generator**: Formats findings in various output formats

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.