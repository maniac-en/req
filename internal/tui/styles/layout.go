package styles

import "github.com/charmbracelet/lipgloss"

var (
	HeaderStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Background(Primary).
			Foreground(TextPrimary).
			Bold(true).
			Align(lipgloss.Center)

	FooterStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(TextSecondary).
			Align(lipgloss.Center)

	ContentStyle = lipgloss.NewStyle().
			Padding(1, 2)

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
