package cmd

import (
	"fmt"

	"github.com/Cofoundr-Ng/coxmos-cli/internal/config"
	"github.com/Cofoundr-Ng/coxmos-cli/internal/tui"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear saved authentication",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Clear(); err != nil {
			return fmt.Errorf("logout: %w", err)
		}
		fmt.Println(tui.LogoutIcon.String() + " " + tui.SuccessStyle.Render("Logged out"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
