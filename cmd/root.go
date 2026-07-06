package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Cofoundr-Ng/coxmos-cli/internal/api"
	"github.com/Cofoundr-Ng/coxmos-cli/internal/config"
	"github.com/Cofoundr-Ng/coxmos-cli/internal/tui"
	"github.com/spf13/cobra"
)

var cfg *config.Config
var client *api.Client

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:     "coxmos",
	Version: version,
	Short:   "Coxmos CLI — Deploy at the edge",
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
	RunE: func(cmd *cobra.Command, args []string) error {
		selection, extraArgs := tui.RunMenu()
		return execSelection(selection, extraArgs)
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

func execSelection(selection string, extraArgs []string) error {
	if selection == "" || selection == "exit" {
		return nil
	}

	parts := strings.Split(selection, ":")
	if len(parts) < 2 {
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("executable path: %w", err)
	}

	category, action := parts[0], parts[1]

	var args []string

	switch category {
	case "apps":
		args = []string{"apps", action}
	case "databases":
		args = []string{"databases", action}
	case "redis":
		args = []string{"redis", action}
	case "email":
		args = []string{"email", action}
	case "dns":
		args = []string{"dns", action}
	case "apikeys":
		args = []string{"apikeys", action}
	case "github":
		args = []string{"github", action}
	case "platform":
		args = []string{"platform", action}
	case "account":
		args = []string{action}
	case "system":
		if action == "update" {
			args = []string{"update"}
		}
	}

	args = append(args, extraArgs...)

	cmd := exec.Command(exe, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
