package cmd

import (
	"fmt"
	"os"

	"github.com/Cofoundr-Ng/coxmos-cli/internal/api"
	"github.com/Cofoundr-Ng/coxmos-cli/internal/config"
	"github.com/Cofoundr-Ng/coxmos-cli/internal/tui"
	"github.com/spf13/cobra"
)

var cfg *config.Config
var client *api.Client

var rootCmd = &cobra.Command{
	Use:   "coxmos",
	Short: "Coxmos CLI — Deploy at the edge",
	Long: tui.Logo() + `
  Coxmos is an edge deployment platform.
  Deploy apps, manage databases, Redis, email, DNS, and more.

  Get started:  coxmos login
  Deploy:       coxmos apps deploy
  Help:         coxmos --help
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}
		client = api.New(config.APIEndpoint)
		if cfg.Token != "" {
			client.Token = cfg.Token
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, tui.ErrorStyle.Render("Error: "+err.Error()))
		os.Exit(1)
	}
}

func requireAuth(cmd *cobra.Command, args []string) error {
	if client.Token == "" {
		fmt.Println(tui.ErrorStyle.Render("You must be logged in. Run: coxmos login"))
		os.Exit(1)
	}
	return nil
}
