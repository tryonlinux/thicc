package display

import "github.com/charmbracelet/lipgloss"

var (
	// HeaderStyle for main headers
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("13")).
			MarginBottom(1)

	// InfoStyle for informational text
	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(false)

	// TitleStyle for ASCII art
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)

	// TableBorderStyle for table borders
	TableBorderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("8"))
)
