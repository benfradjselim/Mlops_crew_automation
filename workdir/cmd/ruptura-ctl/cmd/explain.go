package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain <rupture-id>",
	Short: "Get a human-readable narrative explanation of a rupture event",
	Args:  cobra.ExactArgs(1),
	Example: `  ruptura-ctl explain rpt_abc123
  ruptura-ctl explain production/Deployment/payment-api`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := newClient()
		id := args[0]

		narr, err := c.Explain(ctx(), id)
		if err != nil {
			return fmt.Errorf("explain %q: %w", id, err)
		}
		if cfgOutput == "json" {
			return printJSON(narr)
		}

		fmt.Println()
		fmt.Printf("  %s  %s\n", bold("Narrative Explain"), dim(id))
		fmt.Println()

		text := narr.Narrative
		if text == "" {
			text = narr.Summary
		}
		if text == "" {
			fmt.Println(dim("  No narrative available for this rupture."))
			fmt.Println()
			return nil
		}

		// wrap and print with indent
		lines := wrapText(text, 70)
		for _, l := range lines {
			fmt.Printf("  %s\n", l)
		}
		fmt.Println()

		if narr.Summary != "" && narr.Summary != narr.Narrative {
			fmt.Printf("  %s  %s\n\n", bold("Summary:"), strings.TrimSpace(narr.Summary))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
