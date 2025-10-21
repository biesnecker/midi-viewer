package theme

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestDarkTheme(t *testing.T) {
	theme := Dark()

	if theme.Primary == lipgloss.Color("") {
		t.Error("Dark theme Primary color should not be empty")
	}

	if theme.Background == lipgloss.Color("") {
		t.Error("Dark theme Background color should not be empty")
	}

	if theme.Foreground == lipgloss.Color("") {
		t.Error("Dark theme Foreground color should not be empty")
	}
}

func TestLightTheme(t *testing.T) {
	theme := Light()

	if theme.Primary == lipgloss.Color("") {
		t.Error("Light theme Primary color should not be empty")
	}

	if theme.Background == lipgloss.Color("") {
		t.Error("Light theme Background color should not be empty")
	}

	if theme.Foreground == lipgloss.Color("") {
		t.Error("Light theme Foreground color should not be empty")
	}
}

func TestThemesDifferent(t *testing.T) {
	dark := Dark()
	light := Light()

	if dark.Background == light.Background {
		t.Error("Dark and Light themes should have different background colors")
	}

	if dark.Foreground == light.Foreground {
		t.Error("Dark and Light themes should have different foreground colors")
	}
}
