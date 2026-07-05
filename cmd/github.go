package cmd

import (
	"fmt"

	"github.com/Cofoundr-Ng/coxmos-cli/internal/tui"
	"github.com/spf13/cobra"
)

var githubCmd = &cobra.Command{
	Use:     "github",
	Aliases: []string{"gh"},
	Short:   "Manage GitHub integration",
	Long:    `Link GitHub account, list repos, manage app installations.`,
}

var githubUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Show linked GitHub user",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		var res struct {
			Username string `json:"username"`
			Avatar   string `json:"avatar"`
			URL      string `json:"url"`
		}
		_, err := client.Do("GET", "/github/user", nil, &res)
		if err != nil {
			return fmt.Errorf("github user: %w", err)
		}
		fmt.Println(tui.CardStyle.Render(
			tui.LabelStyle.Render("Username:") + " " + tui.ValueStyle.Render(res.Username) + "\n" +
				tui.LabelStyle.Render("Profile:") + " " + tui.ValueStyle.Render(res.URL),
		))
		return nil
	},
}

var githubReposCmd = &cobra.Command{
	Use:   "repos",
	Short: "List GitHub repos accessible via app installation",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		var res struct {
			Repos []struct {
				Name     string `json:"name"`
				FullName string `json:"full_name"`
				CloneURL string `json:"clone_url"`
				Private  bool   `json:"private"`
			} `json:"repos"`
		}
		_, err := client.Do("GET", "/github/repos", nil, &res)
		if err != nil {
			return fmt.Errorf("github repos: %w", err)
		}
		if len(res.Repos) == 0 {
			fmt.Println(tui.DimStyle.Render("No repos found. Install the GitHub App first."))
			fmt.Println(tui.InfoStyle.Render("Install:  coxmos github install"))
			return nil
		}
		for _, r := range res.Repos {
			icon := tui.Bullet.String()
			visibility := "public"
			if r.Private {
				visibility = "private"
			}
			fmt.Printf("  %s %s (%s)\n", icon, tui.ValueStyle.Render(r.FullName), tui.DimStyle.Render(visibility))
		}
		fmt.Println()
		return nil
	},
}

var githubInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the GitHub App on an account",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(tui.InfoStyle.Render("Opening GitHub App installation page..."))
		var res struct {
			URL string `json:"url"`
		}
		_, err := client.Do("GET", "/github/install", nil, &res)
		if err != nil {
			// redirect endpoint, just show the URL
			_ = res
		}
		fmt.Println(tui.InfoStyle.Render("Visit: https://github.com/apps/<app-slug>/installations/new"))
		fmt.Println(tui.DimStyle.Render("After installing, run: coxmos github repos"))
		return nil
	},
}

var githubLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Link GitHub account via OAuth",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(tui.InfoStyle.Render("Opening GitHub OAuth login..."))
		fmt.Println(tui.DimStyle.Render("Visit: https://github.com/login/oauth/authorize?client_id=<id>&scope=read:user,user:email,repo"))
		return nil
	},
}

func init() {
	githubCmd.AddCommand(githubUserCmd)
	githubCmd.AddCommand(githubReposCmd)
	githubCmd.AddCommand(githubInstallCmd)
	githubCmd.AddCommand(githubLoginCmd)
	rootCmd.AddCommand(githubCmd)
}
