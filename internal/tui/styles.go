package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

var (
	colorSafe      = lipgloss.Color("#00ff00")
	colorWarning   = lipgloss.Color("#ffaa00")
	colorCritical  = lipgloss.Color("#ff8800")
	colorEmergency = lipgloss.Color("#ff0000")
	colorPurple    = lipgloss.Color("#7c3aed")
	colorGray      = lipgloss.Color("#666666")
	colorWhite     = lipgloss.Color("#ffffff")
	colorGold      = lipgloss.Color("#FFD700")
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPurple).
			Padding(0, 1)

	levelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorGold)

	statusStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			Padding(0, 1)

	errorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorEmergency).
			Padding(0, 1).
			Margin(1, 2)

	errorTextStyle = lipgloss.NewStyle().
			Foreground(colorEmergency).
			Bold(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(colorWhite).
			Background(colorPurple).
			Padding(0, 2).
			Margin(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			Padding(1, 2)

	quotaLabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWhite).
			Padding(0, 1)

	quotaValueStyle = lipgloss.NewStyle().
			Foreground(colorGray)

	emptyStateStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			Italic(true).
			Padding(2, 0)

	setupTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPurple).
			Padding(0, 1)

	setupPromptStyle = lipgloss.NewStyle().
				Foreground(colorWhite)
)

func getColorForPercentage(pct int) color.Color {
	switch {
	case pct >= 95:
		return colorEmergency
	case pct >= 90:
		return colorCritical
	case pct >= 80:
		return colorWarning
	default:
		return colorSafe
	}
}
