package game

import (
	"time"

	"github.com/google/uuid"
	"github.com/timlinux/macaco/internal/vim"
)

// SessionState represents the state of a game session
type SessionState int

const (
	SessionStateActive SessionState = iota
	SessionStatePaused
	SessionStateCompleted
)

// Session represents an active game session
type Session struct {
	ID             string       `json:"session_id"`
	RoundType      string       `json:"round_type"`
	Tasks          []*Task      `json:"tasks"`
	CurrentIndex   int          `json:"current_task_index"`
	State          SessionState `json:"state"`
	StartedAt      time.Time    `json:"started_at"`
	CompletedAt    *time.Time   `json:"completed_at,omitempty"`
	TaskResults    []TaskResult `json:"task_results"`
	TotalTasks     int          `json:"total_tasks"`
	SkipsRemaining int          `json:"skips_remaining"`

	// Runtime state (not serialized)
	engine       *vim.Engine
	taskStart    time.Time
	keystrokes   int
	keysUsed     string
	hintsUsed    int
	resets       int
	isPaused     bool
	pauseStart   time.Time
	pausedTime   time.Duration
}

// TaskResult stores the result of a completed task
type TaskResult struct {
	TaskID           string        `json:"task_id"`
	Category         TaskCategory  `json:"category"`
	Difficulty       int           `json:"difficulty"`
	TimeMs           int64         `json:"time_ms"`
	Keystrokes       int           `json:"keystrokes"`
	OptimalKeystrokes int          `json:"optimal_keystrokes"`
	Efficiency       float64       `json:"efficiency"`
	Success          bool          `json:"success"`
	KeysUsed         string        `json:"keys_used"`
	Resets           int           `json:"resets"`
	HintsUsed        int           `json:"hints_used"`
	CompletedAt      time.Time     `json:"completed_at"`
}

// NewSession creates a new game session
func NewSession(roundType string, tasks []*Task) *Session {
	return &Session{
		ID:             uuid.New().String(),
		RoundType:      roundType,
		Tasks:          tasks,
		CurrentIndex:   0,
		State:          SessionStateActive,
		StartedAt:      time.Now(),
		TaskResults:    make([]TaskResult, 0, len(tasks)),
		TotalTasks:     len(tasks),
		SkipsRemaining: 5,
	}
}

// CurrentTask returns the current task
func (s *Session) CurrentTask() *Task {
	if s.CurrentIndex >= 0 && s.CurrentIndex < len(s.Tasks) {
		return s.Tasks[s.CurrentIndex]
	}
	return nil
}

// PreviousTask returns the previous task (if any)
func (s *Session) PreviousTask() *Task {
	if s.CurrentIndex > 0 && s.CurrentIndex <= len(s.Tasks) {
		return s.Tasks[s.CurrentIndex-1]
	}
	return nil
}

// NextTask returns the next task (if any)
func (s *Session) NextTask() *Task {
	if s.CurrentIndex < len(s.Tasks)-1 {
		return s.Tasks[s.CurrentIndex+1]
	}
	return nil
}

// StartTask initializes the current task
func (s *Session) StartTask() {
	task := s.CurrentTask()
	if task == nil {
		return
	}

	s.engine = vim.NewEngine(task.Initial)
	s.engine.SetCursorIndex(task.CursorStart)
	s.taskStart = time.Now()
	s.keystrokes = 0
	s.keysUsed = ""
	s.hintsUsed = 0
	s.resets = 0
	s.pausedTime = 0
}

// ProcessKey processes a keystroke and returns the match status
func (s *Session) ProcessKey(key string) MatchStatus {
	if s.engine == nil || s.isPaused {
		return MatchNone
	}

	s.engine.ProcessKey(key)
	s.keystrokes++
	s.keysUsed += key

	return s.CheckMatch()
}

// CheckMatch checks if the current buffer matches the desired state
func (s *Session) CheckMatch() MatchStatus {
	task := s.CurrentTask()
	if task == nil || s.engine == nil {
		return MatchNone
	}

	bufferText := s.engine.Text()
	cursorIdx := s.engine.CursorIndex()

	if task.IsMotionTask() {
		// For motion tasks, check cursor position
		if cursorIdx == task.CursorEnd {
			return MatchComplete
		}
		return MatchInProgress
	}

	// For editing tasks, check text match
	if bufferText == task.Desired {
		return MatchComplete
	}

	if bufferText == task.Initial {
		return MatchNone
	}

	return MatchInProgress
}

// CompleteTask marks the current task as complete and advances
func (s *Session) CompleteTask() *TaskResult {
	task := s.CurrentTask()
	if task == nil {
		return nil
	}

	elapsed := time.Since(s.taskStart) - s.pausedTime
	efficiency := 0.0
	if s.keystrokes > 0 {
		efficiency = (float64(task.OptimalCount) / float64(s.keystrokes)) * 100
		if efficiency > 100 {
			efficiency = 100
		}
	}

	result := TaskResult{
		TaskID:            task.ID,
		Category:          task.Category,
		Difficulty:        task.Difficulty,
		TimeMs:            elapsed.Milliseconds(),
		Keystrokes:        s.keystrokes,
		OptimalKeystrokes: task.OptimalCount,
		Efficiency:        efficiency,
		Success:           true,
		KeysUsed:          s.keysUsed,
		Resets:            s.resets,
		HintsUsed:         s.hintsUsed,
		CompletedAt:       time.Now(),
	}

	s.TaskResults = append(s.TaskResults, result)
	s.CurrentIndex++

	if s.CurrentIndex >= len(s.Tasks) {
		s.State = SessionStateCompleted
		now := time.Now()
		s.CompletedAt = &now
	} else {
		s.StartTask()
	}

	return &result
}

// SkipTask skips the current task
func (s *Session) SkipTask() bool {
	if s.SkipsRemaining <= 0 {
		return false
	}

	task := s.CurrentTask()
	if task == nil {
		return false
	}

	s.SkipsRemaining--

	result := TaskResult{
		TaskID:            task.ID,
		Category:          task.Category,
		Difficulty:        task.Difficulty,
		TimeMs:            60000, // Max time
		Keystrokes:        s.keystrokes,
		OptimalKeystrokes: task.OptimalCount,
		Efficiency:        0,
		Success:           false,
		KeysUsed:          s.keysUsed,
		Resets:            s.resets,
		HintsUsed:         s.hintsUsed,
		CompletedAt:       time.Now(),
	}

	s.TaskResults = append(s.TaskResults, result)
	s.CurrentIndex++

	if s.CurrentIndex >= len(s.Tasks) {
		s.State = SessionStateCompleted
		now := time.Now()
		s.CompletedAt = &now
	} else {
		s.StartTask()
	}

	return true
}

// ResetTask resets the current task
func (s *Session) ResetTask() {
	task := s.CurrentTask()
	if task == nil {
		return
	}

	s.resets++
	s.engine.Reset(task.Initial, task.CursorStart)
	// Timer and keystroke count continue
}

// UseHint records hint usage
func (s *Session) UseHint() {
	s.hintsUsed++
}

// Pause pauses the session timer
func (s *Session) Pause() {
	if !s.isPaused {
		s.isPaused = true
		s.pauseStart = time.Now()
		s.State = SessionStatePaused
	}
}

// Resume resumes the session timer
func (s *Session) Resume() {
	if s.isPaused {
		s.pausedTime += time.Since(s.pauseStart)
		s.isPaused = false
		s.State = SessionStateActive
	}
}

// IsPaused returns whether the session is paused
func (s *Session) IsPaused() bool {
	return s.isPaused
}

// IsComplete returns whether the session is complete
func (s *Session) IsComplete() bool {
	return s.State == SessionStateCompleted
}

// BufferText returns the current buffer text
func (s *Session) BufferText() string {
	if s.engine == nil {
		return ""
	}
	return s.engine.Text()
}

// CursorPosition returns the current cursor position
func (s *Session) CursorPosition() (x, y int) {
	if s.engine == nil {
		return 0, 0
	}
	return s.engine.CursorPosition()
}

// CursorIndex returns the current cursor index
func (s *Session) CursorIndex() int {
	if s.engine == nil {
		return 0
	}
	return s.engine.CursorIndex()
}

// Mode returns the current vim mode
func (s *Session) Mode() vim.Mode {
	if s.engine == nil {
		return vim.ModeNormal
	}
	return s.engine.Mode()
}

// ElapsedTime returns the elapsed time for the current task
func (s *Session) ElapsedTime() time.Duration {
	if s.taskStart.IsZero() {
		return 0
	}
	elapsed := time.Since(s.taskStart) - s.pausedTime
	if s.isPaused {
		elapsed = s.pauseStart.Sub(s.taskStart) - s.pausedTime
	}
	return elapsed
}

// TotalElapsedTime returns total session time
func (s *Session) TotalElapsedTime() time.Duration {
	if s.CompletedAt != nil {
		return s.CompletedAt.Sub(s.StartedAt)
	}
	return time.Since(s.StartedAt)
}

// Progress returns the current progress as a fraction
func (s *Session) Progress() float64 {
	if s.TotalTasks == 0 {
		return 0
	}
	return float64(s.CurrentIndex) / float64(s.TotalTasks)
}

// Keystrokes returns the current keystroke count
func (s *Session) Keystrokes() int {
	return s.keystrokes
}

// MatchStatus represents the match state between buffer and desired
type MatchStatus int

const (
	MatchNone       MatchStatus = iota // Buffer unchanged
	MatchInProgress                    // Buffer modified but doesn't match
	MatchComplete                      // Buffer matches desired
)

func (m MatchStatus) String() string {
	switch m {
	case MatchNone:
		return "none"
	case MatchInProgress:
		return "in_progress"
	case MatchComplete:
		return "complete"
	default:
		return "unknown"
	}
}
