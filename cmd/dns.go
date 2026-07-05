package cmd

import (
	"fmt"

	"github.com/Cofoundr-Ng/coxmos-cli/internal/tui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:     "dns",
	Aliases: []string{"domain"},
	Short:   "Manage DNS and custom domains",
	Long:    `Register, verify, attach custom domains and manage DNS records.`,
}

var dnsRegisterCmd = &cobra.Command{
	Use:   "register <domain>",
	Short: "Register a new domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		domain := args[0]
		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Registering domain " + domain + "...")
		if err := client.RegisterDomain(domain); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("register domain: %w", err)
		}
		fmt.Print("\r" + tui.CheckMark.String() + " " + tui.SuccessStyle.Render("Domain "+domain+" registered"))
		fmt.Println()
		return nil
	},
}

var dnsVerifyCmd = &cobra.Command{
	Use:   "verify <domain>",
	Short: "Verify domain ownership",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		domain := args[0]

		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Adding verification record...")
		if err := client.VerifyDomain(domain); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("verify: %w", err)
		}
		fmt.Print("\r" + tui.CheckMark.String() + "\n")

		fmt.Println(tui.InfoStyle.Render("  Add this TXT record to your DNS provider:"))
		fmt.Println(tui.DimStyle.Render("  _coxmos-verify." + domain + "  TXT  \"coxmos-verification=\""))

		fmt.Println()
		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Waiting for DNS propagation (checking)...")

		verified, err := client.CheckDomainVerification(domain)
		if err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("check verification: %w", err)
		}
		if verified {
			fmt.Print("\r" + tui.CheckMark.String() + " " + tui.SuccessStyle.Render("Domain verified!"))
		} else {
			fmt.Print("\r" + tui.InfoStyle.Render("Verification record added. Run 'coxmos dns verify "+domain+"' again after DNS propagates."))
		}
		fmt.Println()
		return nil
	},
}

var dnsRecordsCmd = &cobra.Command{
	Use:   "records <domain>",
	Short: "List DNS records for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		domain := args[0]

		records, err := client.ListDNSRecords(domain)
		if err != nil {
			return fmt.Errorf("list records: %w", err)
		}

		if len(records) == 0 {
			fmt.Println(tui.DimStyle.Render("No DNS records found for " + domain))
			return nil
		}

		fmt.Println(tui.Section(tui.DNSIcon.String() + " DNS Records: " + domain))
		fmt.Println()
		for _, r := range records {
			fmt.Printf("  %s %s  %s  %s\n",
				tui.Bullet.String(),
				tui.ValueStyle.Render(r.Type),
				tui.ValueStyle.Render(r.Name),
				tui.DimStyle.Render(r.Value),
			)
		}
		fmt.Println()
		return nil
	},
}

var dnsRemoveCmd = &cobra.Command{
	Use:   "remove <domain>",
	Short: "Remove a registered domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		domain := args[0]

		var confirm bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Remove domain " + domain + "?").
					Description("This will delete all DNS records for the domain.").
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

		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Removing " + domain + "...")
		if err := client.RemoveDomain(domain); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("remove domain: %w", err)
		}
		fmt.Print("\r" + tui.CheckMark.String() + " " + tui.SuccessStyle.Render("Domain removed"))
		fmt.Println()
		return nil
	},
}

var dnsDKIMCmd = &cobra.Command{
	Use:   "dkim <domain>",
	Short: "Add DKIM signing record for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		domain := args[0]
		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Adding DKIM record for " + domain + "...")
		if err := client.AddDKIMRecord(domain); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("dkim: %w", err)
		}
		fmt.Print("\r" + tui.CheckMark.String() + " " + tui.SuccessStyle.Render("DKIM record added for "+domain))
		fmt.Println()
		return nil
	},
}

var dnsAttachCmd = &cobra.Command{
	Use:   "attach <domain> <app-slug>",
	Short: "Attach a custom domain to an app",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		domain, appSlug := args[0], args[1]

		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Attaching " + domain + " to " + appSlug + "...")
		if err := client.AttachDomain(domain, appSlug); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("attach domain: %w", err)
		}
		fmt.Print("\r" + tui.CheckMark.String() + " " + tui.SuccessStyle.Render(domain+" → "+appSlug))
		fmt.Println()
		return nil
	},
}

func init() {
	dnsCmd.AddCommand(dnsRegisterCmd)
	dnsCmd.AddCommand(dnsVerifyCmd)
	dnsCmd.AddCommand(dnsRecordsCmd)
	dnsCmd.AddCommand(dnsRemoveCmd)
	dnsCmd.AddCommand(dnsDKIMCmd)
	dnsCmd.AddCommand(dnsAttachCmd)
	rootCmd.AddCommand(dnsCmd)
}
