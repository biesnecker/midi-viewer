package theme

import "github.com/charmbracelet/lipgloss"

// Theme defines the color scheme for the application
type Theme struct {
	Primary          lipgloss.Color
	Secondary        lipgloss.Color
	Background       lipgloss.Color
	Foreground       lipgloss.Color
	Muted            lipgloss.Color
	Border           lipgloss.Color
	BorderHighlight  lipgloss.Color
	Success          lipgloss.Color
	Warning          lipgloss.Color
	Error            lipgloss.Color
	ModalBackground  lipgloss.Color
	ModalBorder      lipgloss.Color
}

// Dark returns the dark theme
func Dark() Theme {
	return Theme{
		Primary:          lipgloss.Color("#7aa2f7"),
		Secondary:        lipgloss.Color("#bb9af7"),
		Background:       lipgloss.Color("#1a1b26"),
		Foreground:       lipgloss.Color("#c0caf5"),
		Muted:            lipgloss.Color("#565f89"),
		Border:           lipgloss.Color("#414868"),
		BorderHighlight:  lipgloss.Color("#7aa2f7"),
		Success:          lipgloss.Color("#9ece6a"),
		Warning:          lipgloss.Color("#e0af68"),
		Error:            lipgloss.Color("#f7768e"),
		ModalBackground:  lipgloss.Color("#24283b"),
		ModalBorder:      lipgloss.Color("#7aa2f7"),
	}
}

// Light returns the light theme
func Light() Theme {
	return Theme{
		Primary:          lipgloss.Color("#2e7de9"),
		Secondary:        lipgloss.Color("#9854f1"),
		Background:       lipgloss.Color("#e1e2e7"),
		Foreground:       lipgloss.Color("#3760bf"),
		Muted:            lipgloss.Color("#8990b3"),
		Border:           lipgloss.Color("#a8aecb"),
		BorderHighlight:  lipgloss.Color("#2e7de9"),
		Success:          lipgloss.Color("#587539"),
		Warning:          lipgloss.Color("#8c6c3e"),
		Error:            lipgloss.Color("#f52a65"),
		ModalBackground:  lipgloss.Color("#d5d6db"),
		ModalBorder:      lipgloss.Color("#2e7de9"),
	}
}
