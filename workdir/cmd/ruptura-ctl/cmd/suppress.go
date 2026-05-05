package cmd

import (
	"fmt"
	"time"

	"github.com/benfradjselim/ruptura/pkg/client"
	"github.com/spf13/cobra"
)

var suppressReason string

var suppressCmd = &cobra.Command{
	Use:   "suppress",
	Short: "Manage maintenance windows (suppress rupture actions during deploys)",
	Example: `  ruptura-ctl suppress create "production/Deployment/payment-api" 30m
  ruptura-ctl suppress create "production/*" 1h --reason "upgrade k8s cluster"
  ruptura-ctl suppress list
  ruptura-ctl suppress delete <id>`,
	RunE: getSuppressionsCmd.RunE,
}

var suppressCreateCmd = &cobra.Command{
	Use:   "create <workload> <duration>",
	Short: "Create a maintenance window",
	Long: `Create a maintenance window to suppress action dispatch during planned deploys.

Duration format: 30m, 1h, 2h30m

Examples:
  ruptura-ctl suppress create "production/Deployment/payment-api" 30m
  ruptura-ctl suppress create "production/*" 1h --reason "rolling upgrade"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		workload := args[0]
		durationStr := args[1]

		dur, err := time.ParseDuration(durationStr)
		if err != nil {
			return fmt.Errorf("invalid duration %q (use 30m, 1h, 2h30m): %w", durationStr, err)
		}

		c := newClient()
		now := time.Now().UTC()
		req := client.CreateSuppressionReq{
			Workload: workload,
			Start:    now,
			End:      now.Add(dur),
			Reason:   suppressReason,
		}
		if req.Reason == "" {
			req.Reason = "manual suppress via ruptura-ctl"
		}

		supp, err := c.CreateSuppression(ctx(), req)
		if err != nil {
			return fmt.Errorf("create suppression: %w", err)
		}

		successLine(fmt.Sprintf("Suppression created: %s", cyan(supp.ID)))
		fmt.Printf("  %-16s %s\n", dim("workload"), supp.Workload)
		fmt.Printf("  %-16s %s → %s\n", dim("window"),
			supp.Start.Format("15:04"),
			supp.End.Format("15:04 (Jan 2)"),
		)
		fmt.Printf("  %-16s %s\n\n", dim("reason"), dim(supp.Reason))
		return nil
	},
}

var suppressListCmd = &cobra.Command{
	Use:   "list",
	Short: "List active maintenance windows",
	RunE:  getSuppressionsCmd.RunE,
}

var suppressDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a maintenance window",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := newClient()
		id := args[0]
		if err := c.DeleteSuppression(ctx(), id); err != nil {
			return fmt.Errorf("delete suppression %q: %w", id, err)
		}
		successLine(fmt.Sprintf("Suppression %s deleted.", cyan(id)))
		fmt.Println()
		return nil
	},
}

func init() {
	suppressCmd.AddCommand(suppressCreateCmd)
	suppressCmd.AddCommand(suppressListCmd)
	suppressCmd.AddCommand(suppressDeleteCmd)
	rootCmd.AddCommand(suppressCmd)

	suppressCreateCmd.Flags().StringVarP(&suppressReason, "reason", "r", "", "Reason for the maintenance window")
}
