package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/timlinux/macaco/internal/api"
	"github.com/timlinux/macaco/internal/config"
	"github.com/timlinux/macaco/internal/game"
	"github.com/timlinux/macaco/internal/stats"
)

// View represents the current screen
type View int

const (
	ViewMenu View = iota
	ViewGame
	ViewStats
	ViewHelp
)

// App is the main TUI application
type App struct {
	cfg       *config.Config
	engine    *game.Engine
	client    *api.Client
	roundType string

	// State
	view         View
	session      *game.Session
	sessionID    string
	matchStatus  game.MatchStatus
	showHint     bool
	hintLevel    int
	sessionStats *stats.SessionStats

	// UI state
	styles     *Styles
	width      int
	height     int
	lastUpdate time.Time
}

// NewApp creates a new TUI application
func NewApp(cfg *config.Config, engine *game.Engine, client *api.Client, roundType string) *App {
	return &App{
		cfg:       cfg,
		engine:    engine,
		client:    client,
		roundType: roundType,
		view:      ViewMenu,
		styles:    NewStyles(GetTheme(cfg.Theme)),
	}
}

// Run starts the TUI application
func (a *App) Run() error {
	p := tea.NewProgram(a, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tickCmd(),
	)
}

// tickMsg is sent periodically to update the timer
type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles messages and updates state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, nil

	case tea.KeyMsg:
		return a.handleKeyPress(msg)

	case tickMsg:
		a.lastUpdate = time.Time(msg)
		return a, tickCmd()

	case completeTaskMsg:
		// Auto-advance to next task
		a.completeTask()
		return a, nil
	}

	return a, nil
}

// handleKeyPress handles keyboard input
func (a *App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Global keys
	switch key {
	case "ctrl+c":
		return a, tea.Quit
	}

	switch a.view {
	case ViewMenu:
		return a.handleMenuKeys(key)
	case ViewGame:
		return a.handleGameKeys(msg)
	case ViewStats:
		return a.handleStatsKeys(key)
	case ViewHelp:
		return a.handleHelpKeys(key)
	}

	return a, nil
}

// handleMenuKeys handles keys in menu view
func (a *App) handleMenuKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "1":
		a.roundType = "beginner"
		a.startGame()
	case "2":
		a.roundType = "intermediate"
		a.startGame()
	case "3":
		a.roundType = "advanced"
		a.startGame()
	case "4":
		a.roundType = "expert"
		a.startGame()
	case "5":
		a.roundType = "mixed"
		a.startGame()
	case "q":
		return a, tea.Quit
	case "?":
		a.view = ViewHelp
	}
	return a, nil
}

// handleGameKeys handles keys in game view
func (a *App) handleGameKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Handle special control keys first
	switch key {
	case "ctrl+r":
		a.resetTask()
		return a, nil
	case "ctrl+s":
		a.skipTask()
		return a, nil
	case "ctrl+h":
		a.toggleHint()
		return a, nil
	case "ctrl+p":
		a.togglePause()
		return a, nil
	case "?":
		a.view = ViewHelp
		return a, nil
	}

	// Don't process keys if paused
	if a.session != nil && a.session.IsPaused() {
		return a, nil
	}

	// Process vim keys
	if a.session != nil {
		// Convert tea key to vim key
		vimKey := a.convertKey(msg)
		if vimKey != "" {
			a.matchStatus = a.session.ProcessKey(vimKey)

			// Check for task completion
			if a.matchStatus == game.MatchComplete {
				return a, tea.Tick(time.Duration(a.cfg.AutoAdvanceDelay)*time.Millisecond, func(t time.Time) tea.Msg {
					return completeTaskMsg{}
				})
			}
		}
	}

	return a, nil
}

// completeTaskMsg is sent when a task should be completed
type completeTaskMsg struct{}

// handleStatsKeys handles keys in stats view
func (a *App) handleStatsKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "enter", " ":
		a.view = ViewMenu
	case "q":
		return a, tea.Quit
	}
	return a, nil
}

// handleHelpKeys handles keys in help view
func (a *App) handleHelpKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc", "q", "?":
		if a.session != nil {
			a.view = ViewGame
		} else {
			a.view = ViewMenu
		}
	}
	return a, nil
}

// convertKey converts a tea key to a vim key string
func (a *App) convertKey(msg tea.KeyMsg) string {
	switch msg.Type {
	case tea.KeyEsc:
		return "esc"
	case tea.KeyEnter:
		return "enter"
	case tea.KeyBackspace:
		return "backspace"
	case tea.KeySpace:
		return " "
	case tea.KeyTab:
		return "\t"
	case tea.KeyCtrlR:
		return "\x12"
	default:
		if msg.Type == tea.KeyRunes {
			return string(msg.Runes)
		}
		return msg.String()
	}
}

// startGame starts a new game session
func (a *App) startGame() {
	if a.engine != nil {
		a.session = a.engine.CreateSession(a.roundType)
		a.sessionID = a.session.ID
	} else if a.client != nil {
		resp, err := a.client.CreateSession(a.roundType)
		if err != nil {
			// Handle error - for now just return to menu
			return
		}
		a.sessionID = resp.SessionID
		// Create local session for display purposes
		// In client mode, we'd need to sync state from server
	}

	a.view = ViewGame
	a.matchStatus = game.MatchNone
	a.showHint = false
	a.hintLevel = 0
}

// completeTask completes the current task and advances
func (a *App) completeTask() {
	if a.session == nil {
		return
	}

	a.session.CompleteTask()
	a.matchStatus = game.MatchNone
	a.showHint = false
	a.hintLevel = 0

	if a.session.IsComplete() {
		a.sessionStats, _ = a.engine.GetSessionStats(a.sessionID)
		a.view = ViewStats
	}
}

// resetTask resets the current task
func (a *App) resetTask() {
	if a.session != nil {
		a.session.ResetTask()
		a.matchStatus = game.MatchNone
	}
}

// skipTask skips the current task
func (a *App) skipTask() {
	if a.session != nil {
		if a.session.SkipTask() {
			a.matchStatus = game.MatchNone
			a.showHint = false
			a.hintLevel = 0

			if a.session.IsComplete() {
				a.sessionStats, _ = a.engine.GetSessionStats(a.sessionID)
				a.view = ViewStats
			}
		}
	}
}

// toggleHint toggles hint display
func (a *App) toggleHint() {
	if a.showHint {
		a.hintLevel++
		if a.hintLevel > 2 {
			a.showHint = false
			a.hintLevel = 0
		}
	} else {
		a.showHint = true
		a.hintLevel = 0
		if a.session != nil {
			a.session.UseHint()
		}
	}
}

// togglePause toggles pause state
func (a *App) togglePause() {
	if a.session != nil {
		if a.session.IsPaused() {
			a.session.Resume()
		} else {
			a.session.Pause()
		}
	}
}

// View renders the application
func (a *App) View() string {
	switch a.view {
	case ViewMenu:
		return a.renderMenu()
	case ViewGame:
		return a.renderGame()
	case ViewStats:
		return a.renderStats()
	case ViewHelp:
		return a.renderHelp()
	default:
		return ""
	}
}

// renderMenu renders the menu view
func (a *App) renderMenu() string {
	title := lipgloss.NewStyle().
		Foreground(a.styles.Theme.Primary).
		Bold(true).
		MarginBottom(2).
		Render("MoCaCo - Motion Capture Combatant")

	subtitle := a.styles.Subtitle.Render("Master vim motions through competitive practice")

	menu := strings.Join([]string{
		"",
		"Select a round type:",
		"",
		"  [1] Beginner     - Basic motions and operations",
		"  [2] Intermediate - Counts and text objects",
		"  [3] Advanced     - Complex combinations",
		"  [4] Expert       - Multi-step transformations",
		"  [5] Mixed        - Random difficulty",
		"",
		"  [?] Help",
		"  [q] Quit",
	}, "\n")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		subtitle,
		a.styles.Content.Render(menu),
	)

	return lipgloss.Place(
		a.width, a.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

// renderGame renders the game view
func (a *App) renderGame() string {
	if a.session == nil {
		return "No active session"
	}

	// Header - always at top
	header := a.renderHeader()

	// Footer - always at bottom
	footer := a.renderFooter()

	// Calculate available height for content (between header and footer)
	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)
	contentHeight := a.height - headerHeight - footerHeight - 2 // -2 for spacing

	if contentHeight < 10 {
		contentHeight = 10
	}

	// Build content area
	var content strings.Builder

	// Previous task (dimmed) - only show if we have completed at least one task
	if a.session.CurrentIndex > 0 {
		if prev := a.session.PreviousTask(); prev != nil {
			var prevText string
			if prev.IsMotionTask() {
				prevText = fmt.Sprintf("Previous: %q (cursor moved)", prev.Initial)
			} else {
				prevText = fmt.Sprintf("Previous: %q -> %q", prev.Initial, prev.Desired)
			}
			content.WriteString(a.styles.PreviousTask.Width(a.width).Align(lipgloss.Center).Render(prevText))
			content.WriteString("\n")
		}
	}
	content.WriteString("\n")

	// Current task
	task := a.session.CurrentTask()
	if task != nil {
		// Buffer display - current state with cursor and highlighting
		bufferText := a.session.BufferText()
		cursorIdx := a.session.CursorIndex()

		// Build display buffer with cursor and optional highlighting
		displayBuffer := a.renderBufferWithHighlight(bufferText, cursorIdx, task)

		// Check if this is a motion task (text stays same, only cursor moves)
		isMotionTask := task.IsMotionTask()

		var taskDisplay string
		if isMotionTask {
			// For motion tasks, show:
			// 1. Current buffer with cursor position
			// 2. The same text below (reference) with caret showing target
			targetPos := task.CursorEnd

			// Build caret line: spaces up to target, then ^
			// Must be same length as desired text for alignment
			caretLine := strings.Repeat(" ", targetPos) + "^"

			// Combine reference text and caret as a single block so they stay aligned
			referenceBlock := task.Desired + "\n" + caretLine

			// Instruction for motion task
			instruction := "Move your cursor to the caret (^)"

			taskDisplay = lipgloss.JoinVertical(
				lipgloss.Center,
				displayBuffer,
				"",
				a.styles.CurrentTask.Foreground(a.styles.Theme.Dimmed).Render(referenceBlock),
				"",
				a.styles.Subtitle.Foreground(a.styles.Theme.Dimmed).Render(instruction),
			)
		} else {
			// For non-motion tasks, show initial -> desired transformation
			instruction := "Transform the text above to match the text below"

			// Render desired text with highlighting for insertions
			desiredDisplay := a.renderDesiredWithHighlight(task)

			taskDisplay = lipgloss.JoinVertical(
				lipgloss.Center,
				displayBuffer,
				a.styles.Separator.Render("↓"),
				desiredDisplay,
				"",
				a.styles.Subtitle.Foreground(a.styles.Theme.Dimmed).Render(instruction),
			)
		}

		content.WriteString(lipgloss.Place(a.width, contentHeight-6, lipgloss.Center, lipgloss.Center, taskDisplay))
		content.WriteString("\n")

		// Hint (if enabled)
		if a.showHint {
			hint := a.getHint(task)
			content.WriteString(a.styles.Hint.Width(a.width).Align(lipgloss.Center).Render(hint))
			content.WriteString("\n")
		}
	}

	// Next task preview (dimmed) - show the actual next task
	if next := a.session.NextTask(); next != nil {
		var nextText string
		if next.IsMotionTask() {
			nextText = fmt.Sprintf("Next: %q (move cursor)", next.Initial)
		} else {
			nextText = fmt.Sprintf("Next: %q -> %q", next.Initial, next.Desired)
		}
		content.WriteString(a.styles.NextTask.Width(a.width).Align(lipgloss.Center).Render(nextText))
	}

	// Build final layout: header at top, content in middle, footer at bottom
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content.String(),
		footer,
	)
}

// renderHeader renders the header bar
func (a *App) renderHeader() string {
	task := a.session.CurrentTask()

	// Mode
	mode := a.session.Mode().String()
	modeStyle := a.styles.ModeStyle(mode)
	modeStr := modeStyle.Render(fmt.Sprintf("[%s]", mode))

	// Progress - prominent display like baboon
	progress := fmt.Sprintf("Task %d/%d", a.session.CurrentIndex+1, a.session.TotalTasks)

	// Timer
	elapsed := a.session.ElapsedTime()
	timer := fmt.Sprintf("%02d:%02d", int(elapsed.Minutes()), int(elapsed.Seconds())%60)

	// Category
	category := ""
	if task != nil {
		category = string(task.Category)
	}

	// Paused indicator
	paused := ""
	if a.session.IsPaused() {
		paused = a.styles.StatusError.Render(" [PAUSED] ")
	}

	left := fmt.Sprintf("MoCaCo | %s | %s", a.roundType, category)
	center := progress
	right := fmt.Sprintf("%s %s %s", timer, modeStr, paused)

	// Calculate spacing for three-column layout
	headerWidth := a.width - 4
	leftWidth := lipgloss.Width(left)
	centerWidth := lipgloss.Width(center)
	rightWidth := lipgloss.Width(right)

	// Calculate padding on each side of center
	totalUsed := leftWidth + centerWidth + rightWidth
	remainingSpace := headerWidth - totalUsed
	leftPadding := remainingSpace / 2
	rightPadding := remainingSpace - leftPadding

	if leftPadding < 1 {
		leftPadding = 1
	}
	if rightPadding < 1 {
		rightPadding = 1
	}

	return a.styles.Header.Width(a.width).Render(
		left + strings.Repeat(" ", leftPadding) + center + strings.Repeat(" ", rightPadding) + right,
	)
}

// renderFooter renders the footer bar
func (a *App) renderFooter() string {
	hints := []string{
		a.styles.HelpKey.Render("Ctrl+R") + " Reset",
		a.styles.HelpKey.Render("Ctrl+S") + " Skip",
		a.styles.HelpKey.Render("Ctrl+H") + " Hint",
		a.styles.HelpKey.Render("Ctrl+P") + " Pause",
		a.styles.HelpKey.Render("?") + " Help",
	}

	return a.styles.Footer.Width(a.width).Render(strings.Join(hints, "  |  "))
}

// getHint returns the appropriate hint based on level
func (a *App) getHint(task *game.Task) string {
	switch a.hintLevel {
	case 0:
		return fmt.Sprintf("This is a %s task", task.Category)
	case 1:
		return task.Hint
	case 2:
		return fmt.Sprintf("Optimal solution: %s", task.OptimalKeys)
	default:
		return ""
	}
}

// renderStats renders the statistics view
func (a *App) renderStats() string {
	if a.sessionStats == nil {
		return "No statistics available"
	}

	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Foreground(a.styles.Theme.Success).
		Bold(true).
		MarginBottom(1).
		Render("ROUND COMPLETE!")

	b.WriteString(lipgloss.Place(a.width, 3, lipgloss.Center, lipgloss.Center, title))
	b.WriteString("\n")

	// Grade
	grade := a.styles.GradeStyle(a.sessionStats.Grade).Render(
		fmt.Sprintf("Grade: %s", a.sessionStats.Grade),
	)
	b.WriteString(lipgloss.Place(a.width, 1, lipgloss.Center, lipgloss.Center, grade))
	b.WriteString("\n\n")

	// Summary
	summary := fmt.Sprintf(
		"Tasks: %d/%d  |  Time: %s  |  Efficiency: %.1f%%",
		a.sessionStats.TasksCompleted,
		a.sessionStats.TasksAttempted,
		formatDuration(time.Duration(a.sessionStats.TotalTimeMs)*time.Millisecond),
		a.sessionStats.AvgEfficiency,
	)
	b.WriteString(lipgloss.Place(a.width, 1, lipgloss.Center, lipgloss.Center, summary))
	b.WriteString("\n\n")

	// Category breakdown
	b.WriteString(a.styles.Title.Render("Category Breakdown"))
	b.WriteString("\n")

	categories := []string{"motion", "delete", "change", "insert", "visual", "complex"}
	for _, cat := range categories {
		if cs, ok := a.sessionStats.CategoryStats[cat]; ok {
			efficiency := cs.TotalEfficiency / float64(max(cs.TasksAttempted, 1))
			bar := renderProgressBar(efficiency, 20)
			line := fmt.Sprintf("  %-10s %s %.0f%%", cat, bar, efficiency)
			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(a.styles.Hint.Render("Press ENTER to continue"))

	return b.String()
}

// renderHelp renders the help view
func (a *App) renderHelp() string {
	var b strings.Builder

	title := a.styles.Title.Render("MoCaCo Help")
	b.WriteString(lipgloss.Place(a.width, 3, lipgloss.Center, lipgloss.Center, title))
	b.WriteString("\n")

	helpText := `
GAME CONTROLS
  Ctrl+R    Reset current task
  Ctrl+S    Skip current task
  Ctrl+H    Show/cycle hints
  Ctrl+P    Pause/resume timer
  Ctrl+C    Quit

VIM BASICS
  h/j/k/l   Move cursor left/down/up/right
  w/b/e     Word motions
  0/$       Line start/end
  i/a       Insert before/after cursor
  I/A       Insert at line start/end
  o/O       Open line below/above
  d/c/y     Delete/change/yank operators
  x         Delete character
  r         Replace character
  u         Undo
  ESC       Return to normal mode

TEXT OBJECTS
  iw/aw     Inner/around word
  i"/a"     Inner/around quotes
  i(/a(     Inner/around parentheses

FIND MOTIONS
  f{char}   Find character forward
  t{char}   Until character forward
  F{char}   Find character backward
  T{char}   Until character backward

Press ESC or ? to close this help
`
	b.WriteString(helpText)

	return b.String()
}

// Helper functions

// renderBufferWithHighlight renders the buffer text with cursor and character-level highlighting
// Red = characters to delete, Orange = characters to change, Green = cursor target (motion tasks)
func (a *App) renderBufferWithHighlight(text string, cursorIdx int, task *game.Task) string {
	runes := []rune(text)
	statusStyle := a.styles.BufferStyle(a.matchStatus.String())

	// Get character-level highlights if buffer matches initial (hasn't been modified yet)
	if text == task.Initial {
		highlights := task.GetCharHighlights()

		// Define ANSI color codes for inline styling (avoids lipgloss per-char issues)
		// Red for delete, Orange/Yellow for change, Green for target
		const (
			colorReset  = "\033[0m"
			colorRed    = "\033[1;31m" // Bold red
			colorOrange = "\033[1;33m" // Bold yellow/orange
			colorGreen  = "\033[1;32;4m" // Bold green underlined
		)

		// Build the display character by character using ANSI codes
		var result strings.Builder
		for i, r := range runes {
			charStr := string(r)

			// Handle cursor position - show block cursor
			if i == cursorIdx {
				charStr = "█"
			}

			// Apply appropriate color based on highlight type
			if i < len(highlights) {
				switch highlights[i] {
				case game.HighlightDelete:
					result.WriteString(colorRed)
					result.WriteString(charStr)
					result.WriteString(colorReset)
				case game.HighlightChange:
					result.WriteString(colorOrange)
					result.WriteString(charStr)
					result.WriteString(colorReset)
				case game.HighlightTarget:
					result.WriteString(colorGreen)
					result.WriteString(charStr)
					result.WriteString(colorReset)
				default:
					result.WriteString(charStr)
				}
			} else {
				result.WriteString(charStr)
			}
		}

		// Handle cursor at end of text
		if cursorIdx >= len(runes) {
			result.WriteString("█")
		}

		return result.String()
	}

	// Buffer has been modified - no highlighting, just show with cursor
	var displayBuffer string
	if cursorIdx >= 0 && cursorIdx < len(runes) {
		displayBuffer = string(runes[:cursorIdx]) + "█" + string(runes[cursorIdx+1:])
	} else if cursorIdx >= len(runes) && len(runes) > 0 {
		displayBuffer = text + "█"
	} else {
		displayBuffer = text
	}

	return statusStyle.Render(displayBuffer)
}

// renderDesiredWithHighlight renders the desired text with highlighting
// White/bright = characters that need to be inserted, Green = base color
func (a *App) renderDesiredWithHighlight(task *game.Task) string {
	runes := []rune(task.Desired)
	highlights := task.GetDesiredHighlights()

	// Define ANSI color codes
	const (
		colorReset  = "\033[0m"
		colorGreen  = "\033[32m"       // Green for base desired text
		colorWhite  = "\033[1;37m"     // Bold white for insertions
		colorOrange = "\033[1;33m"     // Bold orange for changes
	)

	var result strings.Builder
	for i, r := range runes {
		charStr := string(r)

		if i < len(highlights) {
			switch highlights[i] {
			case game.HighlightInsert:
				result.WriteString(colorWhite)
				result.WriteString(charStr)
				result.WriteString(colorReset)
			case game.HighlightChange:
				result.WriteString(colorOrange)
				result.WriteString(charStr)
				result.WriteString(colorReset)
			default:
				result.WriteString(colorGreen)
				result.WriteString(charStr)
				result.WriteString(colorReset)
			}
		} else {
			result.WriteString(colorGreen)
			result.WriteString(charStr)
			result.WriteString(colorReset)
		}
	}

	return result.String()
}

func formatDuration(d time.Duration) string {
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", m, s)
}

func renderProgressBar(value float64, width int) string {
	filled := int(value / 100 * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	empty := width - filled
	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
