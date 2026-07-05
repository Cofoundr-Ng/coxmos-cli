package cmd

import (
	"fmt"

	"github.com/Cofoundr-Ng/coxmos-cli/internal/api"
	"github.com/Cofoundr-Ng/coxmos-cli/internal/tui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var redisCmd = &cobra.Command{
	Use:     "redis",
	Aliases: []string{"r"},
	Short:   "Manage Redis instances",
	Long:    `Create, list, stop, start, and restart managed Redis instances.`,
}

var redisCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Redis instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		tui.AnimatedTitle(tui.RedisIcon.String() + " Create Redis")
		fmt.Println()

		var name string
		var memory int

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Instance Name").
					Prompt(">").
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("name is required")
						}
						return nil
					}).
					Value(&name),

				huh.NewSelect[int]().
					Title("Memory Limit").
					Options(
						huh.NewOption("256 MB", 256),
						huh.NewOption("512 MB", 512),
						huh.NewOption("1 GB", 1024),
						huh.NewOption("2 GB", 2048),
					).
					Value(&memory),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		fmt.Println()
		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Provisioning Redis instance...")

		inst, err := client.CreateRedis(api.CreateRedisReq{
			Name:     name,
			MemoryMB: memory,
		})
		if err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("create redis: %w", err)
		}

		fmt.Print("\r" + tui.CheckMark.String() + "\n")
		fmt.Println()
		fmt.Println(tui.CardStyle.Render(
			tui.LabelStyle.Render("Name:") + " " + tui.ValueStyle.Render(inst.Name) + "\n" +
				tui.LabelStyle.Render("URI:") + " " + tui.ValueStyle.Render(inst.URI) + "\n" +
				tui.LabelStyle.Render("Status:") + " " + tui.SuccessStyle.Render(inst.Status),
		))
		return nil
	},
}

var redisListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Redis instances",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		fmt.Println(tui.Section(tui.RedisIcon.String() + " Redis Instances"))
		fmt.Println()

		instances, err := client.ListRedis()
		if err != nil {
			return fmt.Errorf("list redis: %w", err)
		}

		if len(instances) == 0 {
			fmt.Println(tui.DimStyle.Render("No Redis instances."))
			fmt.Println(tui.InfoStyle.Render("Create one:  coxmos redis create"))
			return nil
		}

		for _, r := range instances {
			cpu := fmt.Sprintf("%.1f%%", r.CPUPercent)
			mem := fmt.Sprintf("%.0f MB", r.MemoryMB)
			fmt.Printf("  %s %s | CPU: %s | RAM: %s | %s\n",
				tui.Bullet.String(),
				tui.ValueStyle.Render(r.Name),
				tui.InfoStyle.Render(cpu),
				tui.InfoStyle.Render(mem),
				tui.SuccessStyle.Render(r.Status),
			)
		}
		fmt.Println()
		return nil
	},
}

var redisStopCmd = &cobra.Command{
	Use:   "stop <id>",
	Short: "Stop a Redis instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		id := args[0]
		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Stopping " + id + "...")
		if err := client.StopRedis(id); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("stop redis: %w", err)
		}
		fmt.Print("\r" + tui.CheckMark.String() + " " + tui.SuccessStyle.Render("Redis stopped"))
		fmt.Println()
		return nil
	},
}

var redisStartCmd = &cobra.Command{
	Use:   "start <id>",
	Short: "Start a Redis instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		id := args[0]
		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Starting " + id + "...")
		if err := client.StartRedis(id); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("start redis: %w", err)
		}
		fmt.Print("\r" + tui.CheckMark.String() + " " + tui.SuccessStyle.Render("Redis started"))
		fmt.Println()
		return nil
	},
}

var redisRestartCmd = &cobra.Command{
	Use:   "restart <id>",
	Short: "Restart a Redis instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(cmd, args); err != nil {
			return err
		}
		id := args[0]
		fmt.Print(tui.SpinnerStyle.Render("⠋") + " Restarting " + id + "...")
		if err := client.RestartRedis(id); err != nil {
			fmt.Print("\r" + tui.CrossMark.String() + "\n")
			return fmt.Errorf("restart redis: %w", err)
		}
		fmt.Print("\r" + tui.CheckMark.String() + " " + tui.SuccessStyle.Render("Redis restarted"))
		fmt.Println()
		return nil
	},
}

func init() {
	redisCmd.AddCommand(redisCreateCmd)
	redisCmd.AddCommand(redisListCmd)
	redisCmd.AddCommand(redisStopCmd)
	redisCmd.AddCommand(redisStartCmd)
	redisCmd.AddCommand(redisRestartCmd)
	rootCmd.AddCommand(redisCmd)
}
