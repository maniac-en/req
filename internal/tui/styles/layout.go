package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Window and layout styles
	WindowHeaderStyle = lipgloss.NewStyle().
				Padding(1, 2).
				Background(Primary).
				Foreground(TextPrimary).
				Bold(true).
				Align(lipgloss.Center)

	WindowContentStyle = lipgloss.NewStyle().
				Padding(1, 2)

	WindowBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("15")) // White border

	AppBrandingStyle = lipgloss.NewStyle().
				Align(lipgloss.Center).
				Foreground(lipgloss.Color("230")). // Soft cream
				Bold(true).
				Padding(1, 4).
				Margin(1, 0)

	WindowFooterStyle = lipgloss.NewStyle().
				Padding(0, 2).
				Foreground(TextSecondary).
				Align(lipgloss.Center)

	// List and item styles
	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	SelectedListItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(Secondary)

	TitleStyle = lipgloss.NewStyle().
			MarginLeft(2).
			MarginBottom(1).
			Foreground(Primary).
			Bold(true)

	SidebarStyle = lipgloss.NewStyle().
			BorderRight(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(Secondary)

	MainContentStyle = lipgloss.NewStyle().
				PaddingLeft(2)
)
