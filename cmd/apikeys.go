package cmd

import (
	"fmt"

	"github.com/Cofoundr-Ng/coxmos-cli/internal/api"
	"github.com/Cofoundr-Ng/coxmos-cli/internal/tui"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var apikeysCmd = &cobra.Command{
	Use:     "apikeys",
	Aliases: []string{"api-keys", "keys"},
	Short:   "Manage API keys",
	Long:    `Create, list, and revoke API keys for programmatic access.`,
}

var apikeysCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new API key",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		tui.AnimatedTitle(tui.KeyIcon.String() + " Create API Key")
		fmt.Println()

		var name, service string
		var serviceIdx int
		services := []string{"all", "s3", "database", "deploy", "email", "redis"}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Key Name").
					Prompt(">").
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("name is required")
						}
						return nil
					}).
					Value(&name),

				huh.NewSelect[int]().
					Title("Service Scope").
					Options(
						huh.NewOption("All services", 0),
						huh.NewOption("S3 Storage", 1),
						huh.NewOption("Database", 2),
						huh.NewOption("Deploy", 3),
						huh.NewOption("Email", 4),
						huh.NewOption("Redis", 5),
					).
					Value(&serviceIdx),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		service = services[serviceIdx]

		fmt.Println()
		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Generating API key...")

		key, err := client.CreateAPIKey(api.CreateAPIKeyReq{
			Name:    name,
			Service: service,
		})
		if err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("create api key: %w", err)
		}

		fmt.Print("\r" + tui.CheckMark.String() + "\n")
		fmt.Println()
		fmt.Println(tui.CardStyle.Render(
			tui.LabelStyle.Render("Name:") + " " + tui.ValueStyle.Render(key.Name) + "\n" +
				tui.LabelStyle.Render("Key:") + " " + tui.ValueStyle.Render(key.Key) + "\n" +
				tui.LabelStyle.Render("Secret:") + " " + lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700")).Render(key.Secret) + "\n" +
				tui.LabelStyle.Render("Scope:") + " " + tui.ValueStyle.Render(key.Service),
		))
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).Render("⚠ Save the secret now — it won't be shown again!"))
		return nil
	},
}

var apikeysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		fmt.Println(tui.Section(tui.KeyIcon.String() + " API Keys"))
		fmt.Println()

		keys, err := client.ListAPIKeys()
		if err != nil {
			return fmt.Errorf("list api keys: %w", err)
		}

		if len(keys) == 0 {
			fmt.Println(tui.DimStyle.Render("No API keys."))
			fmt.Println(tui.InfoStyle.Render("Create one:  coxmos apikeys create"))
			return nil
		}

		for _, k := range keys {
			secretHint := ""
			if k.SecretHint != "" {
				secretHint = " [" + k.SecretHint + "]"
			}
			fmt.Printf("  %s %s | %s%s | %s\n",
				tui.Bullet.String(),
				tui.ValueStyle.Render(k.Name),
				tui.InfoStyle.Render(k.Service),
				tui.DimStyle.Render(secretHint),
				tui.DimStyle.Render(k.ID),
			)
		}
		fmt.Println()
		return nil
	},
}

var apikeysDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Revoke an API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		id := args[0]

		var confirm bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Revoke API key " + id + "?").
					Description("Applications using this key will lose access.").
					Value(&confirm),
			),
		)
		if err := form.Run(); err != nil {
			return err
		}
		if !confirm {
			fmt.Println(tui.InfoStyle.Render("Cancelled."))
			return nil
		}

		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Revoking...")
		if err := client.DeleteAPIKey(id); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("revoke: %w", err)
		}
		fmt.Print("\r" + tui.CheckMark.String() + " " + tui.SuccessStyle.Render("API key revoked"))
		fmt.Println()
		return nil
	},
}

func init() {
	apikeysCmd.AddCommand(apikeysCreateCmd)
	apikeysCmd.AddCommand(apikeysListCmd)
	apikeysCmd.AddCommand(apikeysDeleteCmd)
	rootCmd.AddCommand(apikeysCmd)
}
