package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme represents a color theme
type Theme struct {
	Name string

	// Backgrounds
	Background     lipgloss.Color
	TaskBackground lipgloss.Color

	// Text colors
	Foreground     lipgloss.Color
	Dimmed         lipgloss.Color
	Bright         lipgloss.Color

	// Status colors
	Success   lipgloss.Color
	Error     lipgloss.Color
	Warning   lipgloss.Color
	Info      lipgloss.Color

	// Accent colors
	Primary   lipgloss.Color
	Secondary lipgloss.Color
}

var (
	// DarkTheme is the default dark color theme
	DarkTheme = Theme{
		Name:           "dark",
		Background:     lipgloss.Color("#0a0e14"),
		TaskBackground: lipgloss.Color("#1a1f29"),
		Foreground:     lipgloss.Color("#b3b1ad"),
		Dimmed:         lipgloss.Color("#6b7280"),
		Bright:         lipgloss.Color("#ffffff"),
		Success:        lipgloss.Color("#10b981"),
		Error:          lipgloss.Color("#ef4444"),
		Warning:        lipgloss.Color("#f59e0b"),
		Info:           lipgloss.Color("#0ea5e9"),
		Primary:        lipgloss.Color("#6366f1"),
		Secondary:      lipgloss.Color("#8b5cf6"),
	}

	// LightTheme is a light color theme
	LightTheme = Theme{
		Name:           "light",
		Background:     lipgloss.Color("#f9fafb"),
		TaskBackground: lipgloss.Color("#ffffff"),
		Foreground:     lipgloss.Color("#1f2937"),
		Dimmed:         lipgloss.Color("#9ca3af"),
		Bright:         lipgloss.Color("#000000"),
		Success:        lipgloss.Color("#059669"),
		Error:          lipgloss.Color("#dc2626"),
		Warning:        lipgloss.Color("#d97706"),
		Info:           lipgloss.Color("#0284c7"),
		Primary:        lipgloss.Color("#4f46e5"),
		Secondary:      lipgloss.Color("#7c3aed"),
	}

	// HighContrastTheme is a high contrast theme
	HighContrastTheme = Theme{
		Name:           "high-contrast",
		Background:     lipgloss.Color("#000000"),
		TaskBackground: lipgloss.Color("#000000"),
		Foreground:     lipgloss.Color("#ffffff"),
		Dimmed:         lipgloss.Color("#808080"),
		Bright:         lipgloss.Color("#ffffff"),
		Success:        lipgloss.Color("#00ff00"),
		Error:          lipgloss.Color("#ff0000"),
		Warning:        lipgloss.Color("#ffff00"),
		Info:           lipgloss.Color("#00ffff"),
		Primary:        lipgloss.Color("#ff00ff"),
		Secondary:      lipgloss.Color("#00ffff"),
	}
)

// GetTheme returns a theme by name
func GetTheme(name string) Theme {
	switch name {
	case "light":
		return LightTheme
	case "high-contrast":
		return HighContrastTheme
	default:
		return DarkTheme
	}
}

// Styles holds all the lipgloss styles for the TUI
type Styles struct {
	Theme Theme

	// Layout styles
	App      lipgloss.Style
	Header   lipgloss.Style
	Footer   lipgloss.Style
	Content  lipgloss.Style

	// Task display styles
	PreviousTask lipgloss.Style
	CurrentTask  lipgloss.Style
	NextTask     lipgloss.Style
	Separator    lipgloss.Style

	// Status styles
	StatusNormal   lipgloss.Style
	StatusInsert   lipgloss.Style
	StatusVisual   lipgloss.Style
	StatusComplete lipgloss.Style
	StatusError    lipgloss.Style
	StatusProgress lipgloss.Style

	// Text styles
	Title      lipgloss.Style
	Subtitle   lipgloss.Style
	Label      lipgloss.Style
	Value      lipgloss.Style
	Hint       lipgloss.Style

	// Stats styles
	Grade       lipgloss.Style
	StatLabel   lipgloss.Style
	StatValue   lipgloss.Style
	ProgressBar lipgloss.Style

	// Help styles
	HelpKey  lipgloss.Style
	HelpDesc lipgloss.Style
}

// NewStyles creates a new Styles instance with the given theme
func NewStyles(theme Theme) *Styles {
	return &Styles{
		Theme: theme,

		App: lipgloss.NewStyle().
			Background(theme.Background),

		Header: lipgloss.NewStyle().
			Foreground(theme.Foreground).
			Background(lipgloss.Color("#111827")).
			Padding(0, 2).
			Bold(true),

		Footer: lipgloss.NewStyle().
			Foreground(theme.Dimmed).
			Background(lipgloss.Color("#111827")).
			Padding(0, 2),

		Content: lipgloss.NewStyle().
			Padding(1, 2),

		PreviousTask: lipgloss.NewStyle().
			Foreground(theme.Dimmed).
			Align(lipgloss.Center),

		CurrentTask: lipgloss.NewStyle().
			Foreground(theme.Bright).
			Background(theme.TaskBackground).
			Padding(2, 4).
			Align(lipgloss.Center).
			Bold(true),

		NextTask: lipgloss.NewStyle().
			Foreground(theme.Dimmed).
			Align(lipgloss.Center),

		Separator: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true),

		StatusNormal: lipgloss.NewStyle().
			Foreground(theme.Info).
			Bold(true),

		StatusInsert: lipgloss.NewStyle().
			Foreground(theme.Success).
			Bold(true),

		StatusVisual: lipgloss.NewStyle().
			Foreground(theme.Secondary).
			Bold(true),

		StatusComplete: lipgloss.NewStyle().
			Foreground(theme.Success).
			Bold(true),

		StatusError: lipgloss.NewStyle().
			Foreground(theme.Error).
			Bold(true),

		StatusProgress: lipgloss.NewStyle().
			Foreground(theme.Warning).
			Bold(true),

		Title: lipgloss.NewStyle().
			Foreground(theme.Bright).
			Bold(true).
			MarginBottom(1),

		Subtitle: lipgloss.NewStyle().
			Foreground(theme.Foreground).
			MarginBottom(1),

		Label: lipgloss.NewStyle().
			Foreground(theme.Dimmed),

		Value: lipgloss.NewStyle().
			Foreground(theme.Foreground),

		Hint: lipgloss.NewStyle().
			Foreground(theme.Dimmed).
			Italic(true),

		Grade: lipgloss.NewStyle().
			Bold(true).
			Padding(0, 2),

		StatLabel: lipgloss.NewStyle().
			Foreground(theme.Dimmed).
			Width(20),

		StatValue: lipgloss.NewStyle().
			Foreground(theme.Foreground).
			Bold(true),

		ProgressBar: lipgloss.NewStyle().
			Foreground(theme.Primary),

		HelpKey: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true),

		HelpDesc: lipgloss.NewStyle().
			Foreground(theme.Dimmed),
	}
}

// BufferStyle returns the appropriate style for buffer text based on match status
func (s *Styles) BufferStyle(status string) lipgloss.Style {
	switch status {
	case "complete":
		return s.CurrentTask.Copy().Foreground(s.Theme.Success)
	case "in_progress":
		return s.CurrentTask.Copy().Foreground(s.Theme.Warning)
	default:
		return s.CurrentTask
	}
}

// GradeStyle returns the appropriate style for a grade
func (s *Styles) GradeStyle(grade string) lipgloss.Style {
	var color lipgloss.Color
	switch grade {
	case "S":
		color = lipgloss.Color("#fbbf24") // Gold
	case "A":
		color = s.Theme.Success
	case "B":
		color = s.Theme.Info
	case "C":
		color = s.Theme.Warning
	case "D":
		color = lipgloss.Color("#f97316") // Orange
	default:
		color = s.Theme.Error
	}
	return s.Grade.Copy().Foreground(color)
}

// ModeStyle returns the appropriate style for a vim mode
func (s *Styles) ModeStyle(mode string) lipgloss.Style {
	switch mode {
	case "INSERT":
		return s.StatusInsert
	case "VISUAL", "V-LINE", "V-BLOCK":
		return s.StatusVisual
	default:
		return s.StatusNormal
	}
}
