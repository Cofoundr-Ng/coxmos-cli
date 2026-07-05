package cmd

import (
	"fmt"

	"github.com/Cofoundr-Ng/coxmos-cli/internal/tui"
	"github.com/spf13/cobra"
)

var platformCmd = &cobra.Command{
	Use:     "platform",
	Aliases: []string{"infra"},
	Short:   "Platform health and information",
	Long:    `Check platform health, architecture, and runtime stats.`,
}

var platformHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Show platform health",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		var res struct {
			Architecture string                   `json:"architecture"`
			Isolates     map[string][]interface{} `json:"isolates"`
			IsolateCount int                      `json:"isolate_count"`
			VMs          []interface{}            `json:"vms"`
			VMCount      int                      `json:"vm_count"`
		}
		_, err := client.Do("GET", "/platform/health", nil, &res)
		if err != nil {
			return fmt.Errorf("platform health: %w", err)
		}
		fmt.Println(tui.Section("🏗️ Platform Health"))
		fmt.Println()
		fmt.Println(tui.CardStyle.Render(
			tui.LabelStyle.Render("Architecture:") + " " + tui.ValueStyle.Render(res.Architecture) + "\n" +
				tui.LabelStyle.Render("Isolates:") + " " + tui.ValueStyle.Render(fmt.Sprintf("%d", res.IsolateCount)) + "\n" +
				tui.LabelStyle.Render("VMs:") + " " + tui.ValueStyle.Render(fmt.Sprintf("%d", res.VMCount)),
		))
		return nil
	},
}

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Check API connectivity",
	RunE: func(cmd *cobra.Command, args []string) error {
		var res struct {
			Message string `json:"message"`
		}
		_, err := client.Do("GET", "/ping", nil, &res)
		if err != nil {
			return fmt.Errorf("ping: %w", err)
		}
		fmt.Println(tui.CheckMark.String() + " " + tui.SuccessStyle.Render(res.Message))
		return nil
	},
}

func init() {
	platformCmd.AddCommand(platformHealthCmd)
	platformCmd.AddCommand(pingCmd)
	rootCmd.AddCommand(platformCmd)
}
