package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

package cmd

import (
"fmt"
"github.com/spf13/cobra"
"github.com/vibhordubey333/k8s-devguardian-ai/internal/scanner"
)

var auditCmd = &cobra.Command{
	Use: "audit",
	Short: "Audits K8s cluster",
	Run : func(cmd *cobra.Command, args []string) {
		fmt.Println("🕵️ Running cluster audit...")
		if err := scanner.ScanCluster(); err != nil {
			fmt.Printf("❌ Error during scan: %v\n", err)
		} else {
			fmt.Println("✅ Scan completed!")
		}
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
}