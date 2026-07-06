package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ArgDef struct {
	Prompt string
}

type MenuItem struct {
	Title       string
	Description string
	Command     string
	Icon        string
	Submenu     []MenuItem
	Args        []ArgDef
}

type menuLayer struct {
	items  []MenuItem
	cursor int
	title  string
}

type MenuModel struct {
	stack      []menuLayer
	quit       bool
	Selected   string
	SelectedArgs []string
	width      int
	height     int
	inputMode  bool
	inputCursor int
	inputs     []string
}

func MainMenuItems() []MenuItem {
	return []MenuItem{
		{
			Title: "Applications", Icon: "🚀",
			Submenu: []MenuItem{
				{Title: "Deploy an App", Description: "From a Git repository", Command: "apps:deploy", Icon: "🚀"},
				{Title: "List Apps", Description: "Show all deployed applications", Command: "apps:list", Icon: "📋"},
				{Title: "View Build Logs", Description: "View deployment build logs", Command: "apps:logs", Icon: "📜", Args: []ArgDef{{Prompt: "Deployment ID:"}}},
				{Title: "Stop an App", Description: "Stop a running application", Command: "apps:stop", Icon: "⏹", Args: []ArgDef{{Prompt: "App slug:"}}},
				{Title: "Start an App", Description: "Start a stopped application", Command: "apps:start", Icon: "▶", Args: []ArgDef{{Prompt: "App slug:"}}},
				{Title: "Restart an App", Description: "Restart an application", Command: "apps:restart", Icon: "🔄", Args: []ArgDef{{Prompt: "App slug:"}}},
			},
		},
		{
			Title: "Databases", Icon: "🗄️",
			Submenu: []MenuItem{
				{Title: "Create a Database", Description: "PostgreSQL or MySQL", Command: "databases:create", Icon: "➕"},
				{Title: "List Databases", Description: "Show all databases", Command: "databases:list", Icon: "📋"},
				{Title: "Delete a Database", Description: "Permanently remove a database", Command: "databases:delete", Icon: "🗑️", Args: []ArgDef{{Prompt: "Database ID:"}}},
			},
		},
		{
			Title: "Redis", Icon: "📀",
			Submenu: []MenuItem{
				{Title: "Create Redis Instance", Description: "Provision a new Redis instance", Command: "redis:create", Icon: "➕"},
				{Title: "List Redis Instances", Description: "Show all Redis instances", Command: "redis:list", Icon: "📋"},
				{Title: "Stop Redis", Description: "Stop a Redis instance", Command: "redis:stop", Icon: "⏹", Args: []ArgDef{{Prompt: "Redis instance ID:"}}},
				{Title: "Start Redis", Description: "Start a Redis instance", Command: "redis:start", Icon: "▶", Args: []ArgDef{{Prompt: "Redis instance ID:"}}},
				{Title: "Restart Redis", Description: "Restart a Redis instance", Command: "redis:restart", Icon: "🔄", Args: []ArgDef{{Prompt: "Redis instance ID:"}}},
			},
		},
		{
			Title: "Email", Icon: "📧",
			Submenu: []MenuItem{
				{Title: "Create Email Account", Description: "Create a new email account", Command: "email:create", Icon: "➕"},
				{Title: "Add Email Domain", Description: "Add and verify a domain for email", Command: "email:domain", Icon: "🌐"},
			},
		},
		{
			Title: "DNS & Domains", Icon: "🌐",
			Submenu: []MenuItem{
				{Title: "Register Domain", Description: "Register a new domain", Command: "dns:register", Icon: "➕", Args: []ArgDef{{Prompt: "Domain:"}}},
				{Title: "Verify Domain", Description: "Verify domain ownership", Command: "dns:verify", Icon: "✅", Args: []ArgDef{{Prompt: "Domain:"}}},
				{Title: "List DNS Records", Description: "Show DNS records for a domain", Command: "dns:records", Icon: "📋", Args: []ArgDef{{Prompt: "Domain:"}}},
				{Title: "Add DKIM Record", Description: "Add DKIM signing record", Command: "dns:dkim", Icon: "🔐", Args: []ArgDef{{Prompt: "Domain:"}}},
				{Title: "Attach Domain to App", Description: "Attach a custom domain to an app", Command: "dns:attach", Icon: "🔗", Args: []ArgDef{{Prompt: "Domain:"}, {Prompt: "App slug:"}}},
				{Title: "Remove Domain", Description: "Remove a registered domain", Command: "dns:remove", Icon: "🗑️", Args: []ArgDef{{Prompt: "Domain:"}}},
			},
		},
		{
			Title: "API Keys", Icon: "🔑",
			Submenu: []MenuItem{
				{Title: "Create API Key", Description: "Generate a new API key", Command: "apikeys:create", Icon: "➕"},
				{Title: "List API Keys", Description: "Show all API keys", Command: "apikeys:list", Icon: "📋"},
				{Title: "Revoke API Key", Description: "Revoke an API key", Command: "apikeys:delete", Icon: "🗑️", Args: []ArgDef{{Prompt: "API Key ID:"}}},
			},
		},
		{
			Title: "GitHub", Icon: "🔗",
			Submenu: []MenuItem{
				{Title: "Show GitHub User", Description: "Show linked GitHub account", Command: "github:user", Icon: "👤"},
				{Title: "List Repos", Description: "List accessible GitHub repos", Command: "github:repos", Icon: "📋"},
				{Title: "Install GitHub App", Description: "Install the GitHub App", Command: "github:install", Icon: "📦"},
				{Title: "Login with GitHub", Description: "Link GitHub via OAuth", Command: "github:login", Icon: "🔑"},
			},
		},
		{
			Title: "Platform", Icon: "🏗️",
			Submenu: []MenuItem{
				{Title: "Health Check", Description: "Show platform health", Command: "platform:health", Icon: "❤️"},
				{Title: "Ping", Description: "Check API connectivity", Command: "platform:ping", Icon: "📡"},
			},
		},
		{
			Title: "Account", Icon: "👤",
			Submenu: []MenuItem{
				{Title: "Login", Description: "Authenticate with Coxmos", Command: "account:login", Icon: "🔐"},
				{Title: "Logout", Description: "Clear saved authentication", Command: "account:logout", Icon: "🚪"},
			},
		},
		{
			Title: "System", Icon: "⚙️",
			Submenu: []MenuItem{
				{Title: "Update CLI", Description: "Update coxmos to the latest version", Command: "system:update", Icon: "📥"},
				{Title: "Exit", Description: "Exit the CLI", Command: "exit", Icon: "❌"},
			},
		},
	}
}

var (
	menuTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Cyan)

	menuItemStyle = lipgloss.NewStyle().
			Padding(0, 1, 0, 2)

	menuSelectedStyle = lipgloss.NewStyle().
				Padding(0, 1, 0, 1).
				Foreground(DarkBg).
				Background(Cyan)

	menuDescStyle = lipgloss.NewStyle().
			Foreground(Gray).
			Italic(true)

	menuHelpStyle = lipgloss.NewStyle().
			Foreground(Gray)

	menuBackStyle = lipgloss.NewStyle().
			Foreground(Yellow).
			Bold(true)

	inputPromptStyle = lipgloss.NewStyle().
				Foreground(White).
				Bold(true)

	inputTextStyle = lipgloss.NewStyle().
			Foreground(Cyan)

	breadcrumbStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	breadcrumbSepStyle = lipgloss.NewStyle().
				Foreground(Gray)
)

func NewMenuModel() MenuModel {
	m := MenuModel{}
	m.stack = append(m.stack, menuLayer{items: MainMenuItems(), cursor: 0, title: "Coxmos"})
	m.width = 80
	m.height = 24
	return m
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) current() menuLayer {
	return m.stack[len(m.stack)-1]
}

func (m *MenuModel) pushLayer(items []MenuItem, title string) {
	m.stack = append(m.stack, menuLayer{items: items, cursor: 0, title: title})
}

func (m *MenuModel) popLayer() {
	if len(m.stack) > 1 {
		m.stack = m.stack[:len(m.stack)-1]
	}
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.inputMode {
			return m.updateInput(msg)
		}
		return m.updateNavigate(msg)
	}

	return m, nil
}

func (m MenuModel) updateNavigate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	cur := m.current()

	switch msg.String() {
	case "ctrl+c", "q":
		m.quit = true
		m.Selected = "exit"
		return m, tea.Quit

	case "up", "k":
		if cur.cursor > 0 {
			m.stack[len(m.stack)-1].cursor--
		}

	case "down", "j":
		if cur.cursor < len(cur.items)-1 {
			m.stack[len(m.stack)-1].cursor++
		}

	case "enter", "right", "l":
		item := cur.items[cur.cursor]
		if item.Command != "" {
			if len(item.Args) > 0 {
				m.inputMode = true
				m.inputCursor = 0
				m.inputs = make([]string, len(item.Args))
				m.Selected = item.Command
			} else {
				m.Selected = item.Command
				m.quit = true
				return m, tea.Quit
			}
		} else if len(item.Submenu) > 0 {
			m.pushLayer(item.Submenu, item.Title)
		}

	case "esc", "left", "backspace", "h":
		if len(m.stack) > 1 {
			m.popLayer()
		}
	}

	return m, nil
}

func (m MenuModel) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	curItem := m.current().items[m.current().cursor]

	switch msg.String() {
	case "ctrl+c":
		m.quit = true
		m.Selected = "exit"
		return m, tea.Quit

	case "esc":
		m.inputMode = false
		m.inputs = nil
		m.inputCursor = 0

	case "enter":
		if m.inputCursor < len(curItem.Args)-1 {
			m.inputCursor++
		} else {
			m.SelectedArgs = m.inputs
			m.quit = true
			return m, tea.Quit
		}

	case "backspace":
		if len(m.inputs[m.inputCursor]) > 0 {
			runes := []rune(m.inputs[m.inputCursor])
			m.inputs[m.inputCursor] = string(runes[:len(runes)-1])
		}

	default:
		if msg.Type == tea.KeyRunes {
			m.inputs[m.inputCursor] += string(msg.Runes)
		}
	}

	return m, nil
}

func breadcrumb(layers []menuLayer) string {
	parts := make([]string, 0, len(layers))
	for _, l := range layers {
		parts = append(parts, l.title)
	}
	return breadcrumbStyle.Render(strings.Join(parts, " "+breadcrumbSepStyle.Render("›")+" "))
}

func (m MenuModel) View() string {
	if m.inputMode {
		return m.inputView()
	}

	cur := m.current()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(breadcrumb(m.stack))
	b.WriteString("\n\n")

	maxItems := m.height - 6
	start := 0
	if cur.cursor >= maxItems/2 {
		start = cur.cursor - maxItems/2 + 1
	}
	if start+maxItems > len(cur.items) {
		start = len(cur.items) - maxItems
		if start < 0 {
			start = 0
		}
	}

	visible := cur.items[start:]
	if len(visible) > maxItems {
		visible = visible[:maxItems]
	}

	for i, item := range visible {
		idx := start + i
		icon := item.Icon
		if icon == "" {
			icon = " "
		}
		hasSub := ""
		if len(item.Submenu) > 0 {
			hasSub = " ›"
		}
		desc := ""
		if item.Description != "" {
			desc = "  " + menuDescStyle.Render(item.Description)
		}

		line := fmt.Sprintf("%s %s%s%s",
			icon,
			item.Title,
			hasSub,
			desc,
		)

		if idx == cur.cursor {
			b.WriteString(menuSelectedStyle.Render(" " + line))
		} else {
			b.WriteString(menuItemStyle.Render(line))
		}
		b.WriteString("\n")
	}

	if start > 0 {
		b.WriteString(menuItemStyle.Render("  ↑\n"))
	}
	if start+len(visible) < len(cur.items) {
		b.WriteString(menuItemStyle.Render("  ↓\n"))
	}

	b.WriteString("\n")
	var helpParts []string
	if len(m.stack) > 1 {
		helpParts = append(helpParts, menuBackStyle.Render("←")+menuHelpStyle.Render("back"))
	}
	helpParts = append(helpParts, menuHelpStyle.Render("↑↓ navigate"))
	helpParts = append(helpParts, menuHelpStyle.Render("↵ select"))
	helpParts = append(helpParts, menuHelpStyle.Render("q quit"))
	b.WriteString(menuHelpStyle.Render(strings.Join(helpParts, "  ")))
	b.WriteString("\n")

	return b.String()
}

func (m MenuModel) inputView() string {
	curItem := m.current().items[m.current().cursor]
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(menuTitleStyle.Render(curItem.Title))
	b.WriteString("\n\n")

	for i, arg := range curItem.Args {
		value := m.inputs[i]
		displayValue := value
		if i == m.inputCursor {
			displayValue += "█"
		}

		b.WriteString(fmt.Sprintf("  %s%s\n",
			inputPromptStyle.Render(arg.Prompt+" "),
			inputTextStyle.Render(displayValue),
		))

		if i < len(curItem.Args)-1 {
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(menuHelpStyle.Render("↵ confirm  esc cancel  q quit"))
	b.WriteString("\n")

	return b.String()
}

func RunMenu() (string, []string) {
	p := tea.NewProgram(NewMenuModel(), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		return "exit", nil
	}
	model := m.(MenuModel)
	return model.Selected, model.SelectedArgs
}
