package cmd

import (
	"fmt"

	"github.com/Cofoundr-Ng/coxmos-cli/internal/tui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Platform admin commands",
	Long:  `Administer users, apps, databases, DNS. Superadmin has absolute control.`,
}

// --- users ---

var adminUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "List all users",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		users, err := client.AdminListUsers()
		if err != nil { return fmt.Errorf("list users: %w", err) }
		if len(users) == 0 { fmt.Println(tui.DimStyle.Render("No users")); return nil }
		fmt.Println(tui.Section("Users"))
		fmt.Println()
		for _, u := range users {
			tag := ""
			if u.Role == "superadmin" { tag = tui.Star.String()
			} else if u.Role == "admin" { tag = tui.LabelStyle.Render(" [admin]") }
			sus := ""
			if u.Suspended { sus = tui.ErrorStyle.Render(" (suspended)") }
			fmt.Printf("  %s %s%s%s  %s\n", tui.Bullet.String(),
				tui.ValueStyle.Render(u.Email), tag, sus,
				tui.DimStyle.Render(u.FirstName+" "+u.LastName))
		}
		fmt.Println()
		return nil
	},
}

var adminUserCmd = &cobra.Command{
	Use:   "user <id>",
	Short: "Show user details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		u, err := client.AdminGetUser(args[0])
		if err != nil { return fmt.Errorf("get user: %w", err) }
		fmt.Println(tui.Section("User: " + u.Email))
		fmt.Println()
		fmt.Printf("  %s ID:       %s\n", tui.Bullet.String(), tui.ValueStyle.Render(u.ID))
		fmt.Printf("  %s Role:     %s\n", tui.Bullet.String(), tui.ValueStyle.Render(u.Role))
		fmt.Printf("  %s Email:    %s\n", tui.Bullet.String(), tui.ValueStyle.Render(u.Email))
		fmt.Printf("  %s Name:     %s\n", tui.Bullet.String(), tui.ValueStyle.Render(u.FirstName+" "+u.LastName))
		fmt.Printf("  %s Apps:     %d\n", tui.Bullet.String(), u.AppCount)
		fmt.Printf("  %s DBs:      %d\n", tui.Bullet.String(), u.DatabaseCount)
		fmt.Printf("  %s Redis:    %d\n", tui.Bullet.String(), u.RedisCount)
		fmt.Printf("  %s Verified: %v\n", tui.Bullet.String(), u.Verified)
		fmt.Printf("  %s Suspended:%v\n", tui.Bullet.String(), u.Suspended)
		fmt.Println()
		return nil
	},
}

var adminUserSuspendCmd = &cobra.Command{
	Use:   "suspend <id>",
	Short: "Suspend a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		var confirm bool
		huh.NewForm(huh.NewGroup(huh.NewConfirm().Title("Suspend user " + args[0] + "?").Value(&confirm))).Run()
		if !confirm { fmt.Println(tui.InfoStyle.Render("Cancelled.")); return nil }
		if err := client.AdminSuspendUser(args[0]); err != nil { return fmt.Errorf("suspend: %w", err) }
		fmt.Println(tui.SuccessStyle.Render("✓ User suspended"))
		return nil
	},
}

var adminUserRestoreCmd = &cobra.Command{
	Use:   "restore <id>",
	Short: "Restore a suspended user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		if err := client.AdminRestoreUser(args[0]); err != nil { return fmt.Errorf("restore: %w", err) }
		fmt.Println(tui.SuccessStyle.Render("✓ User restored"))
		return nil
	},
}

var adminUserRemoveCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "Permanently remove a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		var confirm bool
		huh.NewForm(huh.NewGroup(huh.NewConfirm().Title("PERMANENTLY remove user " + args[0] + "?").Description("This cannot be undone.").Value(&confirm))).Run()
		if !confirm { fmt.Println(tui.InfoStyle.Render("Cancelled.")); return nil }
		if err := client.AdminDeleteUser(args[0]); err != nil { return fmt.Errorf("remove: %w", err) }
		fmt.Println(tui.SuccessStyle.Render("✓ User removed"))
		return nil
	},
}

// --- apps ---

var adminAppsCmd = &cobra.Command{
	Use:   "apps",
	Short: "List all apps across all users",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		apps, err := client.AdminListApps()
		if err != nil { return fmt.Errorf("list apps: %w", err) }
		if len(apps) == 0 { fmt.Println(tui.DimStyle.Render("No apps")); return nil }
		fmt.Println(tui.Section("Apps (all users)"))
		fmt.Println()
		for _, a := range apps {
			fmt.Printf("  %s %s  %s  %s  %s\n", tui.Bullet.String(),
				tui.ValueStyle.Render(a.Slug), tui.DimStyle.Render(a.UserEmail),
				tui.LabelStyle.Render(a.Status), tui.DimStyle.Render(a.Framework))
		}
		fmt.Println()
		return nil
	},
}

var adminAppStopCmd = &cobra.Command{
	Use:   "stop <slug>",
	Short: "Stop any app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		if err := client.AdminStopApp(args[0]); err != nil { return fmt.Errorf("stop: %w", err) }
		fmt.Println(tui.SuccessStyle.Render("✓ App stopped"))
		return nil
	},
}

var adminAppStartCmd = &cobra.Command{
	Use:   "start <slug>",
	Short: "Start any app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		if err := client.AdminStartApp(args[0]); err != nil { return fmt.Errorf("start: %w", err) }
		fmt.Println(tui.SuccessStyle.Render("✓ App started"))
		return nil
	},
}

var adminAppRestartCmd = &cobra.Command{
	Use:   "restart <slug>",
	Short: "Restart any app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		if err := client.AdminRestartApp(args[0]); err != nil { return fmt.Errorf("restart: %w", err) }
		fmt.Println(tui.SuccessStyle.Render("✓ App restarted"))
		return nil
	},
}

var adminAppRemoveCmd = &cobra.Command{
	Use:   "remove <slug>",
	Short: "Remove any app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		var confirm bool
		huh.NewForm(huh.NewGroup(huh.NewConfirm().Title("Remove app " + args[0] + "?").Value(&confirm))).Run()
		if !confirm { fmt.Println(tui.InfoStyle.Render("Cancelled.")); return nil }
		if err := client.AdminDeleteApp(args[0]); err != nil { return fmt.Errorf("remove: %w", err) }
		fmt.Println(tui.SuccessStyle.Render("✓ App removed"))
		return nil
	},
}

// --- databases ---

var adminDatabasesCmd = &cobra.Command{
	Use:   "databases",
	Short: "List all databases across all users",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		dbs, err := client.AdminListDatabases()
		if err != nil { return fmt.Errorf("list databases: %w", err) }
		if len(dbs) == 0 { fmt.Println(tui.DimStyle.Render("No databases")); return nil }
		fmt.Println(tui.Section("Databases (all users)"))
		fmt.Println()
		for _, d := range dbs {
			fmt.Printf("  %s %s  %s  %s/%s  %s\n", tui.Bullet.String(),
				tui.ValueStyle.Render(d.Name), tui.DimStyle.Render(d.UserEmail),
				tui.LabelStyle.Render(d.DBType), tui.DimStyle.Render(d.Kind),
				tui.ValueStyle.Render(d.Status))
		}
		fmt.Println()
		return nil
	},
}

var adminDbRemoveCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "Remove any database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		var confirm bool
		huh.NewForm(huh.NewGroup(huh.NewConfirm().Title("Remove database " + args[0] + "?").Value(&confirm))).Run()
		if !confirm { fmt.Println(tui.InfoStyle.Render("Cancelled.")); return nil }
		if err := client.AdminDeleteDatabase(args[0]); err != nil { return fmt.Errorf("remove: %w", err) }
		fmt.Println(tui.SuccessStyle.Render("✓ Database removed"))
		return nil
	},
}

// --- redis ---

var adminRedisCmd = &cobra.Command{
	Use:   "redis",
	Short: "List all Redis instances across all users",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		instances, err := client.AdminListRedis()
		if err != nil { return fmt.Errorf("list redis: %w", err) }
		if len(instances) == 0 { fmt.Println(tui.DimStyle.Render("No Redis instances")); return nil }
		fmt.Println(tui.Section("Redis (all users)"))
		fmt.Println()
		for _, r := range instances {
			fmt.Printf("  %s %s  %s  %s  %d MB\n", tui.Bullet.String(),
				tui.ValueStyle.Render(r.Name), tui.DimStyle.Render(r.UserEmail),
				tui.LabelStyle.Render(r.Status), r.MemoryMB)
		}
		fmt.Println()
		return nil
	},
}

// --- DNS ---

var adminDomainRegisterCmd = &cobra.Command{
	Use:   "domain <domain>",
	Short: "Register a domain for a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		if err := client.AdminRegisterDomain(args[0]); err != nil { return fmt.Errorf("register: %w", err) }
		fmt.Println(tui.SuccessStyle.Render("✓ Domain registered"))
		return nil
	},
}

var adminDomainRemoveCmd = &cobra.Command{
	Use:   "domain-remove <domain>",
	Short: "Remove any domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		var confirm bool
		huh.NewForm(huh.NewGroup(huh.NewConfirm().Title("Remove domain " + args[0] + "?").Value(&confirm))).Run()
		if !confirm { fmt.Println(tui.InfoStyle.Render("Cancelled.")); return nil }
		if err := client.AdminRemoveDomain(args[0]); err != nil { return fmt.Errorf("remove: %w", err) }
		fmt.Println(tui.SuccessStyle.Render("✓ Domain removed"))
		return nil
	},
}

// --- stats ---

var adminStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Platform statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil { return err }
		stats, err := client.AdminStats()
		if err != nil { return fmt.Errorf("stats: %w", err) }
		fmt.Println(tui.Section("Platform Stats"))
		fmt.Println()
		fmt.Printf("  %s Users:    %d\n", tui.Arrow.String(), stats.Users)
		fmt.Printf("  %s Apps:     %d\n", tui.Arrow.String(), stats.Apps)
		fmt.Printf("  %s Databases: %d\n", tui.Arrow.String(), stats.Databases)
		fmt.Printf("  %s Redis:    %d\n", tui.Arrow.String(), stats.Redis)
		fmt.Printf("  %s Buckets:  %d\n", tui.Arrow.String(), stats.Buckets)
		fmt.Println()
		return nil
	},
}

func init() {
	adminCmd.AddCommand(adminUsersCmd)
	adminCmd.AddCommand(adminUserCmd)
	adminCmd.AddCommand(adminUserSuspendCmd)
	adminCmd.AddCommand(adminUserRestoreCmd)
	adminCmd.AddCommand(adminUserRemoveCmd)
	adminCmd.AddCommand(adminAppsCmd)
	adminCmd.AddCommand(adminAppStopCmd)
	adminCmd.AddCommand(adminAppStartCmd)
	adminCmd.AddCommand(adminAppRestartCmd)
	adminCmd.AddCommand(adminAppRemoveCmd)
	adminCmd.AddCommand(adminDatabasesCmd)
	adminCmd.AddCommand(adminDbRemoveCmd)
	adminCmd.AddCommand(adminRedisCmd)
	adminCmd.AddCommand(adminDomainRegisterCmd)
	adminCmd.AddCommand(adminDomainRemoveCmd)
	adminCmd.AddCommand(adminStatsCmd)
	rootCmd.AddCommand(adminCmd)
}
