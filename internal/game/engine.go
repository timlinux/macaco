package game

import (
	"sync"

	"github.com/timlinux/macaco/internal/config"
	"github.com/timlinux/macaco/internal/stats"
)

// Engine manages game sessions and state
type Engine struct {
	cfg          *config.Config
	taskDB       *TaskDatabase
	generator    *TaskGenerator
	sessions     map[string]*Session
	statsTracker *stats.Tracker
	mu           sync.RWMutex
}

// NewEngine creates a new game engine
func NewEngine(cfg *config.Config) *Engine {
	var taskDB *TaskDatabase
	if cfg.TasksFile != "" {
		var err error
		taskDB, err = LoadTaskDatabase(cfg.TasksFile)
		if err != nil {
			// Fall back to generated tasks
			taskDB = NewGeneratedTaskDatabase()
		}
	} else {
		// Use procedurally generated tasks by default
		taskDB = NewGeneratedTaskDatabase()
	}

	tracker, _ := stats.NewTracker(cfg.StatsFile)
	generator := NewTaskGenerator()

	return &Engine{
		cfg:          cfg,
		taskDB:       taskDB,
		generator:    generator,
		sessions:     make(map[string]*Session),
		statsTracker: tracker,
	}
}

// CreateSession creates a new game session with procedurally generated tasks
func (e *Engine) CreateSession(roundType string) *Session {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Generate fresh tasks for this session using public domain texts
	tasks := e.generator.GenerateTasksForRound(roundType)

	// Convert to pointers
	taskPtrs := make([]*Task, len(tasks))
	for i := range tasks {
		taskPtrs[i] = &tasks[i]
	}

	session := NewSession(roundType, taskPtrs)
	session.StartTask()

	e.sessions[session.ID] = session
	return session
}

// GetTextAttribution returns attribution for the public domain texts used
func (e *Engine) GetTextAttribution() string {
	return e.generator.GetAttribution()
}

// GetSession retrieves a session by ID
func (e *Engine) GetSession(sessionID string) *Session {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.sessions[sessionID]
}

// DeleteSession removes a session
func (e *Engine) DeleteSession(sessionID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.sessions, sessionID)
}

// ProcessKey processes a keystroke for a session
func (e *Engine) ProcessKey(sessionID string, key string) (MatchStatus, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	session, ok := e.sessions[sessionID]
	if !ok {
		return MatchNone, ErrSessionNotFound
	}

	status := session.ProcessKey(key)
	return status, nil
}

// CompleteTask marks the current task as complete
func (e *Engine) CompleteTask(sessionID string) (*TaskResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	session, ok := e.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}

	result := session.CompleteTask()

	// If session is complete, save stats
	if session.IsComplete() && e.statsTracker != nil {
		e.saveSessionStats(session)
	}

	return result, nil
}

// SkipTask skips the current task
func (e *Engine) SkipTask(sessionID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	session, ok := e.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	if !session.SkipTask() {
		return ErrNoSkipsRemaining
	}

	return nil
}

// ResetTask resets the current task
func (e *Engine) ResetTask(sessionID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	session, ok := e.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	session.ResetTask()
	return nil
}

// UseHint records a hint usage
func (e *Engine) UseHint(sessionID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	session, ok := e.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	session.UseHint()
	return nil
}

// PauseSession pauses a session
func (e *Engine) PauseSession(sessionID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	session, ok := e.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	session.Pause()
	return nil
}

// ResumeSession resumes a session
func (e *Engine) ResumeSession(sessionID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	session, ok := e.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	session.Resume()
	return nil
}

// GetSessionStats returns statistics for a completed session
func (e *Engine) GetSessionStats(sessionID string) (*stats.SessionStats, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	session, ok := e.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}

	return e.calculateSessionStats(session), nil
}

// GetLifetimeStats returns lifetime statistics
func (e *Engine) GetLifetimeStats() *stats.LifetimeStats {
	if e.statsTracker != nil {
		return e.statsTracker.GetLifetime()
	}
	return nil
}

// GetTask returns a task by ID
func (e *Engine) GetTask(taskID string) *Task {
	return e.taskDB.GetTask(taskID)
}

// GetRoundTypes returns available round types
func (e *Engine) GetRoundTypes() []string {
	return []string{"beginner", "intermediate", "advanced", "expert", "mixed"}
}

// calculateSessionStats calculates statistics for a session
func (e *Engine) calculateSessionStats(session *Session) *stats.SessionStats {
	if len(session.TaskResults) == 0 {
		return nil
	}

	sessionStats := &stats.SessionStats{
		SessionID:      session.ID,
		RoundType:      session.RoundType,
		StartedAt:      session.StartedAt,
		TotalTimeMs:    session.TotalElapsedTime().Milliseconds(),
		TasksCompleted: 0,
		TasksAttempted: len(session.TaskResults),
		CategoryStats:  make(map[string]*stats.CategoryStats),
	}

	if session.CompletedAt != nil {
		sessionStats.CompletedAt = *session.CompletedAt
	}

	var totalEfficiency float64
	var totalTimeMs int64

	for _, result := range session.TaskResults {
		if result.Success {
			sessionStats.TasksCompleted++
		}
		totalEfficiency += result.Efficiency
		totalTimeMs += result.TimeMs

		// Update category stats
		cat := string(result.Category)
		if sessionStats.CategoryStats[cat] == nil {
			sessionStats.CategoryStats[cat] = &stats.CategoryStats{}
		}
		cs := sessionStats.CategoryStats[cat]
		cs.TasksAttempted++
		if result.Success {
			cs.TasksCompleted++
		}
		cs.TotalTimeMs += result.TimeMs
		cs.TotalEfficiency += result.Efficiency
	}

	if sessionStats.TasksAttempted > 0 {
		sessionStats.AvgEfficiency = totalEfficiency / float64(sessionStats.TasksAttempted)
		sessionStats.AvgTimeMs = totalTimeMs / int64(sessionStats.TasksAttempted)
	}

	sessionStats.Grade = calculateGrade(sessionStats)

	return sessionStats
}

// saveSessionStats saves session statistics
func (e *Engine) saveSessionStats(session *Session) {
	sessionStats := e.calculateSessionStats(session)
	if sessionStats != nil {
		e.statsTracker.RecordSession(sessionStats)
		e.statsTracker.Save()
	}
}

// calculateGrade calculates the grade for a session
func calculateGrade(s *stats.SessionStats) string {
	if s.TasksAttempted == 0 {
		return "F"
	}

	completionRate := float64(s.TasksCompleted) / float64(s.TasksAttempted)
	efficiency := s.AvgEfficiency

	// Target times by difficulty (in ms)
	// Simplified: using overall average
	avgTargetTime := int64(8000) // 8 seconds average

	if completionRate == 1.0 && efficiency >= 95 && s.AvgTimeMs <= avgTargetTime {
		return "S"
	}
	if completionRate == 1.0 && efficiency >= 85 && s.AvgTimeMs <= int64(float64(avgTargetTime)*1.2) {
		return "A"
	}
	if completionRate >= 0.9 && efficiency >= 75 {
		return "B"
	}
	if completionRate >= 0.75 && efficiency >= 60 {
		return "C"
	}
	if completionRate >= 0.5 {
		return "D"
	}
	return "F"
}

// Errors
type GameError string

func (e GameError) Error() string {
	return string(e)
}

const (
	ErrSessionNotFound   GameError = "session not found"
	ErrTaskNotFound      GameError = "task not found"
	ErrNoSkipsRemaining  GameError = "no skips remaining"
	ErrSessionCompleted  GameError = "session already completed"
)
