package stats

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// Tracker manages statistics persistence and aggregation
type Tracker struct {
	filePath string
	data     *StatsData
}

// StatsData represents the complete statistics file
type StatsData struct {
	Version     string         `json:"version"`
	UserID      string         `json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	LastUpdated time.Time      `json:"last_updated"`
	Lifetime    *LifetimeStats `json:"lifetime"`
	Sessions    []*SessionStats `json:"sessions"`
	Achievements []Achievement  `json:"achievements"`
	Preferences  *Preferences   `json:"preferences"`
}

// LifetimeStats aggregates statistics over all sessions
type LifetimeStats struct {
	TotalRounds        int                       `json:"total_rounds"`
	TotalTasks         int                       `json:"total_tasks"`
	TotalTimeMs        int64                     `json:"total_time_ms"`
	TotalKeystrokes    int                       `json:"total_keystrokes"`
	TotalPracticeTimeMs int64                    `json:"total_practice_time_ms"`
	ByCategory         map[string]*CategoryStats `json:"by_category"`
	PersonalBests      *PersonalBests            `json:"personal_bests"`
}

// CategoryStats holds statistics for a task category
type CategoryStats struct {
	TasksAttempted  int     `json:"tasks_attempted"`
	TasksCompleted  int     `json:"tasks_completed"`
	TotalTimeMs     int64   `json:"total_time_ms"`
	TotalKeystrokes int     `json:"total_keystrokes"`
	TotalEfficiency float64 `json:"total_efficiency"`
	BestTimeMs      int64   `json:"best_time_ms"`
	AvgTimeMs       int64   `json:"avg_time_ms"`
	AvgEfficiency   float64 `json:"avg_efficiency"`
	SuccessRate     float64 `json:"success_rate"`
}

// PersonalBests tracks personal records
type PersonalBests struct {
	FastestTask    *BestRecord `json:"fastest_task,omitempty"`
	BestEfficiency *BestRecord `json:"best_efficiency,omitempty"`
	FastestRound   *BestRecord `json:"fastest_round,omitempty"`
}

// BestRecord represents a personal best
type BestRecord struct {
	ID         string    `json:"id"`
	Value      float64   `json:"value"`
	Date       time.Time `json:"date"`
	RoundType  string    `json:"round_type,omitempty"`
}

// SessionStats holds statistics for a single session
type SessionStats struct {
	SessionID      string                    `json:"session_id"`
	RoundType      string                    `json:"round_type"`
	StartedAt      time.Time                 `json:"started_at"`
	CompletedAt    time.Time                 `json:"completed_at,omitempty"`
	TotalTimeMs    int64                     `json:"total_time_ms"`
	TasksCompleted int                       `json:"tasks_completed"`
	TasksAttempted int                       `json:"tasks_attempted"`
	Grade          string                    `json:"grade"`
	AvgEfficiency  float64                   `json:"avg_efficiency"`
	AvgTimeMs      int64                     `json:"avg_time_ms"`
	CategoryStats  map[string]*CategoryStats `json:"category_stats,omitempty"`
	Tasks          []*TaskStats              `json:"tasks,omitempty"`
}

// TaskStats holds statistics for a single task attempt
type TaskStats struct {
	TaskID            string    `json:"task_id"`
	Category          string    `json:"category"`
	Difficulty        int       `json:"difficulty"`
	TimeMs            int64     `json:"time_ms"`
	Keystrokes        int       `json:"keystrokes"`
	OptimalKeystrokes int       `json:"optimal_keystrokes"`
	Efficiency        float64   `json:"efficiency"`
	Success           bool      `json:"success"`
	KeysUsed          string    `json:"keys_used"`
	Resets            int       `json:"resets"`
	HintsUsed         int       `json:"hints_used"`
	CompletedAt       time.Time `json:"completed_at"`
}

// Achievement represents an unlocked achievement
type Achievement struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UnlockedAt  time.Time `json:"unlocked_at"`
}

// Preferences stores user preferences
type Preferences struct {
	Theme              string  `json:"theme"`
	AutoAdvanceDelayMs int     `json:"auto_advance_delay_ms"`
	ShowHints          bool    `json:"show_hints"`
	EnableSounds       bool    `json:"enable_sounds"`
	AnimationSpeed     float64 `json:"animation_speed"`
}

// NewTracker creates a new statistics tracker
func NewTracker(filePath string) (*Tracker, error) {
	t := &Tracker{
		filePath: filePath,
	}

	if err := t.load(); err != nil {
		// Initialize new stats data
		t.data = t.newStatsData()
	}

	return t, nil
}

// newStatsData creates a new stats data structure
func (t *Tracker) newStatsData() *StatsData {
	return &StatsData{
		Version:     "1.0.0",
		UserID:      uuid.New().String(),
		CreatedAt:   time.Now(),
		LastUpdated: time.Now(),
		Lifetime: &LifetimeStats{
			ByCategory:    make(map[string]*CategoryStats),
			PersonalBests: &PersonalBests{},
		},
		Sessions:     make([]*SessionStats, 0),
		Achievements: make([]Achievement, 0),
		Preferences: &Preferences{
			Theme:              "dark",
			AutoAdvanceDelayMs: 500,
			ShowHints:          true,
			EnableSounds:       false,
			AnimationSpeed:     1.0,
		},
	}
}

// load loads statistics from file
func (t *Tracker) load() error {
	data, err := os.ReadFile(t.filePath)
	if err != nil {
		return err
	}

	t.data = &StatsData{}
	if err := json.Unmarshal(data, t.data); err != nil {
		return err
	}

	// Ensure maps are initialized
	if t.data.Lifetime == nil {
		t.data.Lifetime = &LifetimeStats{
			ByCategory:    make(map[string]*CategoryStats),
			PersonalBests: &PersonalBests{},
		}
	}
	if t.data.Lifetime.ByCategory == nil {
		t.data.Lifetime.ByCategory = make(map[string]*CategoryStats)
	}

	return nil
}

// Save saves statistics to file
func (t *Tracker) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(t.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	t.data.LastUpdated = time.Now()

	data, err := json.MarshalIndent(t.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(t.filePath, data, 0644)
}

// RecordSession records a completed session
func (t *Tracker) RecordSession(session *SessionStats) {
	// Add to sessions list
	t.data.Sessions = append(t.data.Sessions, session)

	// Limit stored sessions
	if len(t.data.Sessions) > 100 {
		t.data.Sessions = t.data.Sessions[len(t.data.Sessions)-100:]
	}

	// Update lifetime stats
	t.updateLifetimeStats(session)

	// Check for new achievements
	t.checkAchievements(session)
}

// updateLifetimeStats updates aggregate statistics
func (t *Tracker) updateLifetimeStats(session *SessionStats) {
	lifetime := t.data.Lifetime

	lifetime.TotalRounds++
	lifetime.TotalTasks += session.TasksAttempted
	lifetime.TotalTimeMs += session.TotalTimeMs
	lifetime.TotalPracticeTimeMs += session.TotalTimeMs

	// Update category stats
	for cat, catStats := range session.CategoryStats {
		if lifetime.ByCategory[cat] == nil {
			lifetime.ByCategory[cat] = &CategoryStats{}
		}
		lc := lifetime.ByCategory[cat]
		lc.TasksAttempted += catStats.TasksAttempted
		lc.TasksCompleted += catStats.TasksCompleted
		lc.TotalTimeMs += catStats.TotalTimeMs
		lc.TotalEfficiency += catStats.TotalEfficiency

		// Update averages
		if lc.TasksAttempted > 0 {
			lc.AvgTimeMs = lc.TotalTimeMs / int64(lc.TasksAttempted)
			lc.AvgEfficiency = lc.TotalEfficiency / float64(lc.TasksAttempted)
			lc.SuccessRate = float64(lc.TasksCompleted) / float64(lc.TasksAttempted) * 100
		}

		// Update best time
		if catStats.BestTimeMs > 0 && (lc.BestTimeMs == 0 || catStats.BestTimeMs < lc.BestTimeMs) {
			lc.BestTimeMs = catStats.BestTimeMs
		}
	}

	// Update personal bests
	if session.TotalTimeMs > 0 {
		if lifetime.PersonalBests.FastestRound == nil ||
			session.TotalTimeMs < int64(lifetime.PersonalBests.FastestRound.Value) {
			lifetime.PersonalBests.FastestRound = &BestRecord{
				ID:        session.SessionID,
				Value:     float64(session.TotalTimeMs),
				Date:      session.CompletedAt,
				RoundType: session.RoundType,
			}
		}
	}
}

// checkAchievements checks for newly unlocked achievements
func (t *Tracker) checkAchievements(session *SessionStats) {
	// First Steps - complete first round
	if t.data.Lifetime.TotalRounds == 1 {
		t.unlockAchievement("first-steps", "First Steps", "Complete your first round")
	}

	// Dedicated - complete 10 rounds
	if t.data.Lifetime.TotalRounds == 10 {
		t.unlockAchievement("dedicated", "Dedicated", "Complete 10 rounds")
	}

	// Expert - complete 100 rounds
	if t.data.Lifetime.TotalRounds == 100 {
		t.unlockAchievement("expert", "Expert", "Complete 100 rounds")
	}

	// Perfect efficiency - 100% efficiency on a task
	if session.AvgEfficiency >= 100 {
		t.unlockAchievement("optimal-path", "Optimal Path", "Achieve 100% efficiency on a task")
	}

	// Flawless Victory - complete round without mistakes
	if session.TasksCompleted == session.TasksAttempted && session.AvgEfficiency >= 95 {
		t.unlockAchievement("flawless-victory", "Flawless Victory", "Complete a round with 95%+ efficiency")
	}
}

// unlockAchievement unlocks an achievement if not already unlocked
func (t *Tracker) unlockAchievement(id, name, description string) {
	// Check if already unlocked
	for _, a := range t.data.Achievements {
		if a.ID == id {
			return
		}
	}

	t.data.Achievements = append(t.data.Achievements, Achievement{
		ID:          id,
		Name:        name,
		Description: description,
		UnlockedAt:  time.Now(),
	})
}

// GetLifetime returns lifetime statistics
func (t *Tracker) GetLifetime() *LifetimeStats {
	return t.data.Lifetime
}

// GetRecentSessions returns recent sessions
func (t *Tracker) GetRecentSessions(count int) []*SessionStats {
	sessions := t.data.Sessions
	if len(sessions) <= count {
		return sessions
	}
	return sessions[len(sessions)-count:]
}

// GetAchievements returns all achievements
func (t *Tracker) GetAchievements() []Achievement {
	return t.data.Achievements
}

// GetPreferences returns user preferences
func (t *Tracker) GetPreferences() *Preferences {
	return t.data.Preferences
}

// SetPreferences updates user preferences
func (t *Tracker) SetPreferences(prefs *Preferences) {
	t.data.Preferences = prefs
}

// ExportJSON exports statistics as JSON
func (t *Tracker) ExportJSON() ([]byte, error) {
	return json.MarshalIndent(t.data, "", "  ")
}
