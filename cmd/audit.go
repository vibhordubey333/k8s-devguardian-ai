package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vibhordubey333/k8s-devguardian-ai/internal/ai"
	"github.com/vibhordubey333/k8s-devguardian-ai/internal/output"
	"github.com/vibhordubey333/k8s-devguardian-ai/internal/scanner"
	"os"
)

var (
	outputFormat string
	aiProvider   string
	apiKey       string
	modelName    string
	ollamaURL    string
	outputFile   string
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audits K8s cluster",
	Long:  `Performs a security audit on your Kubernetes cluster and provides AI-powered explanations and remediation suggestions.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🕵️ Running cluster audit...")

		// Scan the cluster
		findings, err := scanner.ScanCluster()
		if err != nil {
			fmt.Printf("❌ Error during scan: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Scan completed! Found %d potential security issues.\n", len(findings))

		// If no findings, exit early
		if len(findings) == 0 {
			fmt.Println("🎉 No security issues found!")
			return
		}

		// Initialize AI explainer
		fmt.Println("🧠 Analyzing findings with AI...")
		explainer, err := ai.NewExplainer(ai.ExplainerConfig{
			Provider:  aiProvider,
			APIKey:    apiKey,
			ModelName: modelName,
			OllamaURL: ollamaURL,
		})
		if err != nil {
			fmt.Printf("⚠️ Warning: Could not initialize AI explainer: %v\n", err)
			fmt.Println("⚠️ Continuing with basic explanations...")
			explainer = ai.NewSimpleExplainer()
		}

		// Get explanations for findings
		explanations, err := explainer.ExplainFindings(findings)
		if err != nil {
			fmt.Printf("⚠️ Warning: Error getting AI explanations: %v\n", err)
			fmt.Println("⚠️ Continuing with basic explanations...")
			explainer = ai.NewSimpleExplainer()
			explanations, _ = explainer.ExplainFindings(findings)
		}

		// Format the output
		fmt.Println("📊 Generating report...")
		formatter := output.NewFormatter(output.Format(outputFormat))
		result := output.AuditResult{
			Findings:     findings,
			Explanations: explanations,
			Summary:      output.GenerateSummary(findings),
		}

		reportBytes, err := formatter.Format(result)
		if err != nil {
			fmt.Printf("❌ Error formatting report: %v\n", err)
			os.Exit(1)
		}

		// Output the report
		if outputFile != "" {
			err := os.WriteFile(outputFile, reportBytes, 0644)
			if err != nil {
				fmt.Printf("❌ Error writing report to file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Report saved to %s\n", outputFile)
		} else {
			fmt.Println(string(reportBytes))
		}
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)

	// Add flags
	auditCmd.Flags().StringVarP(&outputFormat, "output", "o", "cli", "Output format (cli, json, html)")
	auditCmd.Flags().StringVarP(&aiProvider, "ai-provider", "a", "", "AI provider (openai, ollama)")
	auditCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "API key for OpenAI")
	auditCmd.Flags().StringVarP(&modelName, "model", "m", "", "Model name to use")
	auditCmd.Flags().StringVarP(&ollamaURL, "ollama-url", "u", "http://localhost:11434", "URL for Ollama server")
	auditCmd.Flags().StringVarP(&outputFile, "file", "f", "", "Output file path")
}
