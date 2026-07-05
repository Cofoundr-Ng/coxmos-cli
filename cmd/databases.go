package cmd

import (
	"fmt"

	"github.com/Cofoundr-Ng/coxmos-cli/internal/api"
	"github.com/Cofoundr-Ng/coxmos-cli/internal/tui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:     "databases",
	Aliases: []string{"db", "database"},
	Short:   "Manage databases",
	Long:    `Create, list, and delete PostgreSQL and MySQL databases.`,
}

var dbCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		tui.AnimatedTitle(tui.Database.String() + " Create Database")
		fmt.Println()

		var name string
		var dbTypeIdx, kindIdx int

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Database Name").
					Prompt(">").
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("name is required")
						}
						return nil
					}).
					Value(&name),

				huh.NewSelect[int]().
					Title("Type").
					Options(
						huh.NewOption("PostgreSQL", 0),
						huh.NewOption("MySQL", 1),
					).
					Value(&dbTypeIdx),

				huh.NewSelect[int]().
					Title("Kind").
					Description("Isolated = dedicated container, Provisioned = shared server").
					Options(
						huh.NewOption("Isolated (dedicated container)", 0),
						huh.NewOption("Provisioned (shared server)", 1),
					).
					Value(&kindIdx),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		dbTypes := []string{"postgres", "mysql"}
		kinds := []string{"isolated", "provisioned"}

		fmt.Println()
		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Provisioning " + dbTypes[dbTypeIdx] + " database...")

		db, err := client.CreateDatabase(api.CreateDBReq{
			DBType: dbTypes[dbTypeIdx],
			Kind:   kinds[kindIdx],
			Name:   name,
		})
		if err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("create database: %w", err)
		}

		fmt.Print("\r" + tui.CheckMark.String() + "\n")
		fmt.Println()
		fmt.Println(tui.CardStyle.Render(
			tui.LabelStyle.Render("Name:") + " " + tui.ValueStyle.Render(db.Name) + "\n" +
				tui.LabelStyle.Render("Type:") + " " + tui.ValueStyle.Render(db.DBType) + "\n" +
				tui.LabelStyle.Render("Kind:") + " " + tui.ValueStyle.Render(db.Kind) + "\n" +
				tui.LabelStyle.Render("Status:") + " " + tui.SuccessStyle.Render(db.Status) + "\n" +
				tui.LabelStyle.Render("Connection:") + " " + tui.ValueStyle.Render(db.ConnectionString),
		))
		return nil
	},
}

var dbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all databases",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		fmt.Println(tui.Section(tui.Database.String() + " Databases"))
		fmt.Println()

		dbs, err := client.ListDatabases()
		if err != nil {
			return fmt.Errorf("list databases: %w", err)
		}

		if len(dbs) == 0 {
			fmt.Println(tui.DimStyle.Render("No databases created yet."))
			fmt.Println(tui.InfoStyle.Render("Create one:  coxmos databases create"))
			return nil
		}

		for _, db := range dbs {
			fmt.Printf("  %s %s | %s/%s | %s\n",
				tui.Bullet.String(),
				tui.ValueStyle.Render(db.Name),
				tui.InfoStyle.Render(db.DBType),
				tui.InfoStyle.Render(db.Kind),
				tui.SuccessStyle.Render(db.Status),
			)
		}
		fmt.Println()
		return nil
	},
}

var dbDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a database",
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
					Title("Delete database " + id + "?").
					Description("This action cannot be undone.").
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

		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Deleting...")
		if err := client.DeleteDatabase(id); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("delete: %w", err)
		}
		fmt.Print("\r" + tui.CheckMark.String() + " " + tui.SuccessStyle.Render("Database deleted"))
		fmt.Println()
		return nil
	},
}

func init() {
	dbCmd.AddCommand(dbCreateCmd)
	dbCmd.AddCommand(dbListCmd)
	dbCmd.AddCommand(dbDeleteCmd)
	rootCmd.AddCommand(dbCmd)
}
