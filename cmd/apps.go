package cmd

import (
	"fmt"
	"os"

	"github.com/Cofoundr-Ng/coxmos-cli/internal/api"
	"github.com/Cofoundr-Ng/coxmos-cli/internal/tui"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var appsCmd = &cobra.Command{
	Use:     "apps",
	Aliases: []string{"app", "a"},
	Short:   "Manage edge applications",
	Long:    `Deploy, list, stop, start, restart, and view logs for your applications.`,
}

var appsDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a new application from Git",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		tui.AnimatedTitle(tui.Rocket.String() + " Deploy App")
		fmt.Println()

		var gitURL, branch string
		var frameworkIdx int
		frameworks := []string{"auto-detect", "next.js", "nuxt", "sveltekit", "astro", "remix", "solidstart", "deno", "node.js", "go", "python", "ruby", "php"}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Git Repository URL").
					Description("e.g. https://github.com/user/my-app.git").
					Prompt(">").
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("repository URL is required")
						}
						return nil
					}).
					Value(&gitURL),

				huh.NewInput().
					Title("Branch").
					Prompt(">").
					Value(&branch),

				huh.NewSelect[int]().
					Title("Framework").
					Options(
						huh.NewOption("Auto-detect", 0),
						huh.NewOption("Next.js", 1),
						huh.NewOption("Nuxt", 2),
						huh.NewOption("SvelteKit", 3),
						huh.NewOption("Astro", 4),
						huh.NewOption("Remix", 5),
						huh.NewOption("SolidStart", 6),
						huh.NewOption("Deno", 7),
						huh.NewOption("Node.js", 8),
						huh.NewOption("Go", 9),
						huh.NewOption("Python", 10),
						huh.NewOption("Ruby", 11),
						huh.NewOption("PHP", 12),
					).
					Value(&frameworkIdx),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		if branch == "" {
			branch = "main"
		}

		framework := ""
		if frameworkIdx > 0 {
			framework = frameworks[frameworkIdx]
		}

		fmt.Println()
		fmt.Println(tui.InfoStyle.Render("Deploying " + gitURL + " (" + branch + ")..."))

		res, err := client.DeployApp(api.DeployReq{
			GitURL:    gitURL,
			Branch:    branch,
			Framework: framework,
		})
		if err != nil {
			return fmt.Errorf("deploy: %w", err)
		}

		fmt.Println(tui.CheckMark.String() + " " + tui.SuccessStyle.Render("Deployed!"))
		fmt.Println()
		fmt.Println(tui.CardStyle.Render(
			tui.LabelStyle.Render("URL:") + " " + tui.ValueStyle.Render("https://"+res.URL) + "\n" +
				tui.LabelStyle.Render("Slug:") + " " + tui.ValueStyle.Render(res.Slug) + "\n" +
				tui.LabelStyle.Render("Status:") + " " + tui.SuccessStyle.Render(res.Status) + "\n" +
				tui.LabelStyle.Render("Deployment:") + " " + tui.DimStyle.Render(res.DeploymentID),
		))
		return nil
	},
}

var appsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all applications",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		fmt.Println(tui.Section(tui.Rocket.String() + " Applications"))
		fmt.Println()

		apps, err := client.ListApps()
		if err != nil {
			return fmt.Errorf("list apps: %w", err)
		}

		if len(apps) == 0 {
			fmt.Println(tui.DimStyle.Render("No applications deployed yet."))
			fmt.Println(tui.InfoStyle.Render("Deploy your first app:  coxmos apps deploy"))
			return nil
		}

		for _, app := range apps {
			slug := app.Slug
			if slug == "" {
				slug = app.Name
			}
			statusTxt := app.Status
			statusColor := lipgloss.Color("#00FF88")
			if app.Status == "stopped" {
				statusColor = lipgloss.Color("#FF4444")
			}
			fmt.Printf("  %s %s %s\n",
				tui.Bullet.String(),
				tui.ValueStyle.Render(slug),
				lipgloss.NewStyle().Foreground(statusColor).Render(statusTxt),
			)
		}
		fmt.Println()
		return nil
	},
}

var appsLogsCmd = &cobra.Command{
	Use:   "logs [deployment-id]",
	Short: "Stream deployment build logs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		id := args[0]
		fmt.Println(tui.Section(tui.Rocket.String() + " Build Logs: " + id))
		fmt.Println()

		logs, err := client.GetLogs(id)
		if err != nil {
			return fmt.Errorf("logs: %w", err)
		}

		if len(logs) == 0 {
			fmt.Println(tui.DimStyle.Render("No logs available."))
			return nil
		}

		for _, entry := range logs {
			if entry.Stream == "stderr" {
				fmt.Fprintln(os.Stderr, tui.ErrorStyle.Render(entry.Line))
			} else {
				fmt.Println(tui.InfoStyle.Render(entry.Line))
			}
		}
		return nil
	},
}

var appsStopCmd = &cobra.Command{
	Use:   "stop <slug>",
	Short: "Stop an application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		slug := args[0]
		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Stopping " + slug + "...")
		if err := client.StopApp(slug); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("stop: %w", err)
		}
		fmt.Print("\r" + tui.CheckMark.String() + " " + tui.SuccessStyle.Render(slug+" stopped"))
		fmt.Println()
		return nil
	},
}

var appsStartCmd = &cobra.Command{
	Use:   "start <slug>",
	Short: "Start an application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		slug := args[0]
		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Starting " + slug + "...")
		if err := client.StartApp(slug); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("start: %w", err)
		}
		fmt.Print("\r" + tui.CheckMark.String() + " " + tui.SuccessStyle.Render(slug+" started"))
		fmt.Println()
		return nil
	},
}

var appsRestartCmd = &cobra.Command{
	Use:   "restart <slug>",
	Short: "Restart an application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		slug := args[0]
		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Restarting " + slug + "...")
		if err := client.RestartApp(slug); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("restart: %w", err)
		}
		fmt.Print("\r" + tui.CheckMark.String() + " " + tui.SuccessStyle.Render(slug+" restarted"))
		fmt.Println()
		return nil
	},
}

func init() {
	appsCmd.AddCommand(appsDeployCmd)
	appsCmd.AddCommand(appsListCmd)
	appsCmd.AddCommand(appsLogsCmd)
	appsCmd.AddCommand(appsStopCmd)
	appsCmd.AddCommand(appsStartCmd)
	appsCmd.AddCommand(appsRestartCmd)
	rootCmd.AddCommand(appsCmd)
}
