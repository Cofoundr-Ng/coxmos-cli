package tui

import "github.com/charmbracelet/lipgloss"

var (
	Subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	Spotlight = lipgloss.AdaptiveColor{Light: "#FF6B6B", Dark: "#FF6B6B"}

	Cyan    = lipgloss.Color("#00D4FF")
	Green   = lipgloss.Color("#00FF88")
	Yellow  = lipgloss.Color("#FFD700")
	Red     = lipgloss.Color("#FF4444")
	Purple  = lipgloss.Color("#BB86FC")
	Orange  = lipgloss.Color("#FF8C00")
	White   = lipgloss.Color("#FFFFFF")
	Gray    = lipgloss.Color("#888888")
	DarkBg  = lipgloss.Color("#1A1A2E")
	CardBg  = lipgloss.Color("#16213E")
	Accent  = lipgloss.Color("#0F3460")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Cyan).
			Padding(0, 1).
			MarginBottom(1)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Purple).
			Padding(0, 1)

	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Green)

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Red)

	InfoStyle = lipgloss.NewStyle().
			Foreground(White)

	DimStyle = lipgloss.NewStyle().
			Foreground(Gray).
			Italic(true)

	LabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Yellow).
			MarginRight(1)

	ValueStyle = lipgloss.NewStyle().
			Foreground(White)

	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Accent).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(Cyan)

	ProgressStyle = lipgloss.NewStyle().
			Foreground(Green).
			Background(lipgloss.Color("#333333"))

	BlinkStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	CheckMark = lipgloss.NewStyle().Foreground(Green).SetString("✓")
	CrossMark = lipgloss.NewStyle().Foreground(Red).SetString("✗")
	Bullet    = lipgloss.NewStyle().Foreground(Cyan).SetString("•")
	Arrow     = lipgloss.NewStyle().Foreground(Purple).SetString("→")
	Star      = lipgloss.NewStyle().Foreground(Yellow).SetString("✦")
	Rocket    = lipgloss.NewStyle().Foreground(Orange).SetString("🚀")
	Database  = lipgloss.NewStyle().Foreground(Green).SetString("🗄️")
	RedisIcon = lipgloss.NewStyle().Foreground(Red).SetString("📀")
	MailIcon  = lipgloss.NewStyle().Foreground(Yellow).SetString("📧")
	DNSIcon   = lipgloss.NewStyle().Foreground(Cyan).SetString("🌐")
	KeyIcon   = lipgloss.NewStyle().Foreground(Purple).SetString("🔑")
	LogoutIcon = lipgloss.NewStyle().Foreground(Gray).SetString("🚪")
)

func Logo() string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(Cyan).
		Render(`
  ██████  ██████  ██   ██ ███    ███  ██████  ███████
 ██      ██    ██ ██  ██  ████  ████ ██    ██ ██
 ██      ██    ██ █████   ██ ████ ██ ██    ██ ███████
 ██      ██    ██ ██  ██  ██  ██  ██ ██    ██      ██
  ██████  ██████  ██   ██ ██      ██  ██████  ███████
`)
}

func Section(title string) string {
	line := lipgloss.NewStyle().Foreground(Accent).Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	return line + "\n" + TitleStyle.Render(title) + "\n" + line
}
