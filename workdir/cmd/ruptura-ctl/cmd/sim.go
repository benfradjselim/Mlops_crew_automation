package cmd

import (
	"fmt"

	"github.com/benfradjselim/ruptura/pkg/client"
	"github.com/spf13/cobra"
)

var (
	simWorkload  string
	simDuration  int
)

var simCmd = &cobra.Command{
	Use:   "sim",
	Short: "Control ruptura-sim synthetic load injection",
	Example: `  ruptura-ctl sim inject cascade-failure
  ruptura-ctl sim inject memory-leak --workload production/Deployment/payment-api --duration 120
  ruptura-ctl sim patterns`,
}

var simInjectCmd = &cobra.Command{
	Use:   "inject <pattern>",
	Short: "Inject a synthetic load pattern",
	Long: `Inject a synthetic load pattern into Ruptura to test alerting and action pipelines.

Available patterns:
  memory-leak       Slowly increasing memory pressure
  cascade-failure   Error propagation across service dependencies
  traffic-surge     Sudden spike in throughput and latency
  slow-burn         Gradual multi-signal degradation over time`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := args[0]
		valid := map[string]bool{
			"memory-leak":     true,
			"cascade-failure": true,
			"traffic-surge":   true,
			"slow-burn":       true,
		}
		if !valid[pattern] {
			return fmt.Errorf("unknown pattern %q — valid: memory-leak, cascade-failure, traffic-surge, slow-burn", pattern)
		}

		c := newClient()
		req := client.SimInjectReq{
			Pattern:         pattern,
			Workload:        simWorkload,
			DurationSeconds: simDuration,
		}
		resp, err := c.SimInject(ctx(), req)
		if err != nil {
			return fmt.Errorf("sim inject: %w", err)
		}

		successLine(fmt.Sprintf("Pattern %s injected", cyan(resp.Pattern)))
		if resp.Workload != "" {
			fmt.Printf("  %-16s %s\n", dim("workload"), resp.Workload)
		}
		if resp.Message != "" {
			fmt.Printf("  %-16s %s\n", dim("message"), dim(resp.Message))
		}
		fmt.Println()
		infoLine("Watch the effect: " + cyan("ruptura-ctl status"))
		fmt.Println()
		return nil
	},
}

var simPatternsCmd = &cobra.Command{
	Use:   "patterns",
	Short: "List available simulation patterns",
	RunE: func(cmd *cobra.Command, args []string) error {
		patterns := []struct{ name, desc string }{
			{"memory-leak", "Slowly increasing memory pressure — triggers fatigue accumulation"},
			{"cascade-failure", "Error propagation across service edges — triggers contagion signal"},
			{"traffic-surge", "Sudden throughput spike + latency — triggers stress + pressure"},
			{"slow-burn", "Gradual multi-signal degradation — tests forecast accuracy"},
		}
		fmt.Println()
		fmt.Printf("  %s\n\n", bold("Simulation Patterns"))
		for _, p := range patterns {
			fmt.Printf("  %s\n    %s\n\n", cyan(p.name), dim(p.desc))
		}
		fmt.Printf("  Usage: %s\n\n", cyan("ruptura-ctl sim inject <pattern> [--workload <ref>] [--duration <seconds>]"))
		return nil
	},
}

func init() {
	simCmd.AddCommand(simInjectCmd)
	simCmd.AddCommand(simPatternsCmd)
	rootCmd.AddCommand(simCmd)

	simInjectCmd.Flags().StringVar(&simWorkload, "workload", "", "Target workload ref (namespace/Kind/name). Omit for auto-select.")
	simInjectCmd.Flags().IntVar(&simDuration, "duration", 0, "Duration in seconds (default: pattern default)")
}
