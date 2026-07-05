package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

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
	case "linux":
	default:
		osName = "linux"
	}
	return "coxmos-" + osName + "-" + arch
}

type ghRelease struct {
	TagName string `json:"tag_name"`
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

		if resp.StatusCode == 403 || resp.StatusCode == 429 {
			return fmt.Errorf("GitHub API rate limited. Try again later or set GH_TOKEN")
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("GitHub API returned %s", resp.Status)
		}

		var release ghRelease
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			return fmt.Errorf("parse response: %w", err)
		}
		if release.TagName == "" {
			return fmt.Errorf("could not determine latest version")
		}

		url := "https://github.com/" + updateRepo + "/releases/download/" + release.TagName + "/" + asset
		fmt.Println(tui.InfoStyle.Render("  Latest: " + release.TagName))
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
			return fmt.Errorf("download failed: %s (asset: %s)", resp.Status, asset)
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
			return fmt.Errorf("replace binary at %s: %w\nTry: sudo mv %s %s", exe, err, tmp.Name(), exe)
		}

		fmt.Print("\r" + tui.CheckMark.String() + "\n")
		fmt.Println()
		fmt.Println(tui.SuccessStyle.Render(fmt.Sprintf("  Updated to %s (%d bytes)", release.TagName, written)))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
