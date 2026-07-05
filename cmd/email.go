package cmd

import (
	"fmt"

	"github.com/Cofoundr-Ng/coxmos-cli/internal/api"
	"github.com/Cofoundr-Ng/coxmos-cli/internal/tui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var emailCmd = &cobra.Command{
	Use:     "email",
	Aliases: []string{"mail"},
	Short:   "Manage email accounts and domains",
	Long:    `Create email accounts and add/verify email domains via Mailu.`,
}

var emailCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an email account",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		tui.AnimatedTitle(tui.MailIcon.String() + " Create Email Account")
		fmt.Println()

		var email, password, displayName string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Email Address").
					Prompt(">").
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("email is required")
						}
						return nil
					}).
					Value(&email),

				huh.NewInput().
					Title("Password").
					Prompt(">").
					EchoMode(huh.EchoModePassword).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("password is required")
						}
						return nil
					}).
					Value(&password),

				huh.NewInput().
					Title("Display Name").
					Prompt(">").
					Value(&displayName),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		fmt.Println()
		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Creating email account...")

		acc, err := client.CreateEmailAccount(api.CreateEmailAccountReq{
			Email:       email,
			Password:    password,
			DisplayName: displayName,
		})
		if err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("create email: %w", err)
		}

		fmt.Print("\r" + tui.CheckMark.String() + "\n")
		fmt.Println()
		fmt.Println(tui.CardStyle.Render(
			tui.LabelStyle.Render("Email:") + " " + tui.ValueStyle.Render(acc.Email) + "\n" +
				tui.LabelStyle.Render("Name:") + " " + tui.ValueStyle.Render(acc.DisplayName),
		))
		return nil
	},
}

var emailDomainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Add an email domain",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		tui.AnimatedTitle(tui.MailIcon.String() + " Add Email Domain")
		fmt.Println()

		var domain string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Domain").
					Prompt(">").
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("domain is required")
						}
						return nil
					}).
					Value(&domain),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		fmt.Println()
		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Adding domain " + domain + "...")
		if err := client.AddEmailDomain(domain); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("add domain: %w", err)
		}

		fmt.Print("\r" + tui.CheckMark.String() + "\n")
		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Verifying domain and generating DKIM...")
		if err := client.VerifyEmailDomain(domain); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("verify domain: %w", err)
		}
		fmt.Print("\r" + tui.CheckMark.String() + " " + tui.SuccessStyle.Render("Domain "+domain+" is ready for email"))
		fmt.Println()
		return nil
	},
}

func init() {
	emailCmd.AddCommand(emailCreateCmd)
	emailCmd.AddCommand(emailDomainCmd)
	rootCmd.AddCommand(emailCmd)
}
