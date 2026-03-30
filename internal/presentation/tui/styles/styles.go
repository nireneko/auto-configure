package styles

import "github.com/charmbracelet/lipgloss"

var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6B6B")).
		MarginBottom(1)

	Subtitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))

	Selected = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD93D")).
		Bold(true)

	Unselected = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC"))

	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6BCB77")).
		Bold(true)

	Failure = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Bold(true)

	Warning = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD93D"))

	Hint = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Italic(true)

	Border = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(1, 2)
)
