package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Primary colors - Warm & Earthy
	Primary   = lipgloss.Color("95")  // Muted reddish-brown (e.g., rust)
	Secondary = lipgloss.Color("101") // Soft olive green
	Success   = lipgloss.Color("107") // Earthy sage green
	Warning   = lipgloss.Color("172") // Warm goldenrod/ochre
	Error     = lipgloss.Color("160") // Deep muted red

	// Text colors
	TextPrimary   = lipgloss.Color("254") // Off-white/cream
	TextSecondary = lipgloss.Color("246") // Medium warm gray
	TextMuted     = lipgloss.Color("241") // Darker warm gray

	// Background colors
	BackgroundPrimary   = lipgloss.Color("235") // Very dark brown-gray
	BackgroundSecondary = lipgloss.Color("238") // Dark brown-gray
)
