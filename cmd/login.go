package cmd

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Cofoundr-Ng/coxmos-cli/internal/config"
	"github.com/Cofoundr-Ng/coxmos-cli/internal/tui"
	"github.com/spf13/cobra"
)

func randHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Coxmos",
	Long:  `Device-code SSO: CLI shows a code, you enter it on the web dashboard.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tui.AnimatedTitle("✦ Welcome to Coxmos ✦")
		fmt.Println()

		raw := randHex(4) // 8 hex chars
		hash := sha256.Sum256([]byte(raw))
		hashStr := hex.EncodeToString(hash[:])

		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Registering code...")
		if err := client.RegisterDeviceCode(hashStr, raw[:4]); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("register code: %w", err)
		}
		fmt.Print("\r" + tui.CheckMark.String() + "\n")
		fmt.Println()

		fmt.Println(tui.InfoStyle.Render("Open in your browser:"))
		fmt.Println(tui.TitleStyle.Render("  https://coxmos.app/cli-login"))
		fmt.Println()
		fmt.Println(tui.InfoStyle.Render("Enter code:"))
		fmt.Println(tui.BlinkStyle.Render("  " + raw))
		fmt.Println()

		spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		si := 0
		for {
			fmt.Print("\r" + tui.SpinnerStyle.Render(spinner[si%len(spinner)]) + " Waiting for authentication...")
			si++

			time.Sleep(2 * time.Second)

			res, err := client.PollDeviceCode(raw)
			if err != nil {
				continue
			}
			if res.Status == "pending" {
				continue
			}

			cfg.Token = res.Token
			cfg.User.Email = res.User.Email
			name := res.User.FirstName
			if res.User.LastName != "" {
				name += " " + res.User.LastName
			}
			cfg.User.Name = name
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			fmt.Print("\r" + tui.CheckMark.String() + " Authenticated!\n")
			fmt.Println()
			fmt.Println(tui.SuccessStyle.Render("✓ Logged in as " + res.User.Email))
			client.Token = res.Token
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
