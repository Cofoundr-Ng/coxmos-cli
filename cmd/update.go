package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/Cofoundr-Ng/coxmos-cli/internal/tui"
	"github.com/spf13/cobra"
)

const updateRepo = "Cofoundr-Ng/coxmos-cli"

func assetName() string {
	arch := runtime.GOARCH
	osName := runtime.GOOS
	switch arch {
	case "x86_64", "amd64":
		arch = "amd64"
	case "aarch64", "arm64":
		arch = "arm64"
	}
	switch osName {
	case "darwin":
		osName = "darwin"
	case "linux":
		osName = "linux"
	}
	return "coxmos-" + osName + "-" + arch
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update coxmos to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		tui.AnimatedTitle("✦ Coxmos Update ✦")
		fmt.Println()

		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("cannot determine executable path: %w", err)
		}

		asset := assetName()

		fmt.Println(tui.InfoStyle.Render("  Checking latest version..."))

		resp, err := http.Get("https://api.github.com/repos/" + updateRepo + "/releases/latest")
		if err != nil {
			return fmt.Errorf("check update: %w", err)
		}
		defer resp.Body.Close()

		var tag string
		body, _ := io.ReadAll(resp.Body)
		for _, line := range strings.Split(string(body), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, `"tag_name"`) {
				tag = strings.Trim(strings.Split(line, ":")[1], ` ",`)
				break
			}
		}
		if tag == "" {
			return fmt.Errorf("could not determine latest version")
		}

		url := "https://github.com/" + updateRepo + "/releases/download/" + tag + "/" + asset
		fmt.Println(tui.InfoStyle.Render("  Latest: " + tag))
		fmt.Println(tui.DimStyle.Render("  Downloading " + asset + "..."))

		tmp, err := os.CreateTemp("", "coxmos-*")
		if err != nil {
			return fmt.Errorf("create temp file: %w", err)
		}
		defer os.Remove(tmp.Name())

		resp, err = http.Get(url)
		if err != nil {
			return fmt.Errorf("download: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("download failed: %s", resp.Status)
		}

		written, err := io.Copy(tmp, resp.Body)
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}
		tmp.Close()

		if err := os.Chmod(tmp.Name(), 0755); err != nil {
			return fmt.Errorf("chmod: %w", err)
		}

		if err := os.Rename(tmp.Name(), exe); err != nil {
			return fmt.Errorf("replace binary at %s: %w", exe, err)
		}

		fmt.Print("\r" + tui.CheckMark.String() + "\n")
		fmt.Println()
		fmt.Println(tui.SuccessStyle.Render(fmt.Sprintf("  Updated to %s (%d bytes)", tag, written)))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
