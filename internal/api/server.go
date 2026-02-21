package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/timlinux/macaco/internal/config"
	"github.com/timlinux/macaco/internal/game"
)

// Server represents the REST API server
type Server struct {
	cfg      *config.Config
	engine   *game.Engine
	server   *http.Server
	startTime time.Time
}

// NewServer creates a new API server
func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:       cfg,
		engine:    game.NewEngine(cfg),
		startTime: time.Now(),
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/v1/health", s.handleHealth)
	mux.HandleFunc("/api/v1/sessions", s.handleSessions)
	mux.HandleFunc("/api/v1/sessions/", s.handleSessionByID)
	mux.HandleFunc("/api/v1/tasks", s.handleTasks)
	mux.HandleFunc("/api/v1/tasks/", s.handleTaskByID)
	mux.HandleFunc("/api/v1/rounds", s.handleRounds)
	mux.HandleFunc("/api/v1/rounds/", s.handleRoundByType)
	mux.HandleFunc("/api/v1/stats/lifetime", s.handleLifetimeStats)
	mux.HandleFunc("/api/v1/stats/export", s.handleStatsExport)

	// Create PID file
	s.writePIDFile()

	s.server = &http.Server{
		Addr:         s.cfg.ServerAddr,
		Handler:      corsMiddleware(loggingMiddleware(mux)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting MoCaCo API server on %s", s.cfg.ServerAddr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.server.Shutdown(ctx)
	}
	s.removePIDFile()
}

// writePIDFile writes the server PID to a file
func (s *Server) writePIDFile() {
	pidDir := s.cfg.DataDir
	if pidDir == "" {
		homeDir, _ := os.UserHomeDir()
		pidDir = filepath.Join(homeDir, ".config", "macaco")
	}
	os.MkdirAll(pidDir, 0755)
	pidFile := filepath.Join(pidDir, "server.pid")
	os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
}

// removePIDFile removes the PID file
func (s *Server) removePIDFile() {
	pidDir := s.cfg.DataDir
	if pidDir == "" {
		homeDir, _ := os.UserHomeDir()
		pidDir = filepath.Join(homeDir, ".config", "macaco")
	}
	pidFile := filepath.Join(pidDir, "server.pid")
	os.Remove(pidFile)
}

// Middleware

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Handlers

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"status":         "ok",
		"version":        "1.0.0",
		"uptime_seconds": int(time.Since(s.startTime).Seconds()),
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createSession(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) createSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoundType string `json:"round_type"`
		UserID    string `json:"user_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.RoundType == "" {
		req.RoundType = "beginner"
	}

	session := s.engine.CreateSession(req.RoundType)
	task := session.CurrentTask()

	response := map[string]interface{}{
		"session_id":         session.ID,
		"round_type":         session.RoundType,
		"total_tasks":        session.TotalTasks,
		"current_task_index": session.CurrentIndex,
		"started_at":         session.StartedAt,
		"current_task":       taskToMap(task),
	}

	writeJSON(w, http.StatusCreated, response)
}

func (s *Server) handleSessionByID(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/sessions/")
	parts := strings.Split(path, "/")
	sessionID := parts[0]

	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "MISSING_SESSION_ID", "Session ID is required")
		return
	}

	// Check for sub-routes
	if len(parts) > 1 {
		switch parts[1] {
		case "keystroke":
			s.handleKeystroke(w, r, sessionID)
			return
		case "keystrokes":
			s.handleKeystrokes(w, r, sessionID)
			return
		case "complete":
			s.handleCompleteTask(w, r, sessionID)
			return
		case "skip":
			s.handleSkipTask(w, r, sessionID)
			return
		case "reset":
			s.handleResetTask(w, r, sessionID)
			return
		case "stats":
			s.handleSessionStats(w, r, sessionID)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		s.getSession(w, r, sessionID)
	case http.MethodDelete:
		s.deleteSession(w, r, sessionID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) getSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	session := s.engine.GetSession(sessionID)
	if session == nil {
		writeError(w, http.StatusNotFound, "SESSION_NOT_FOUND", "Session not found")
		return
	}

	response := map[string]interface{}{
		"session_id":         session.ID,
		"round_type":         session.RoundType,
		"total_tasks":        session.TotalTasks,
		"current_task_index": session.CurrentIndex,
		"started_at":         session.StartedAt,
		"current_task":       taskToMap(session.CurrentTask()),
		"buffer_state":       session.BufferText(),
		"cursor_position":    session.CursorIndex(),
		"elapsed_time_ms":    session.ElapsedTime().Milliseconds(),
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) deleteSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	s.engine.DeleteSession(sessionID)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleKeystroke(w http.ResponseWriter, r *http.Request, sessionID string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Key       string   `json:"key"`
		Modifiers []string `json:"modifiers,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	session := s.engine.GetSession(sessionID)
	if session == nil {
		writeError(w, http.StatusNotFound, "SESSION_NOT_FOUND", "Session not found")
		return
	}

	status := session.ProcessKey(req.Key)

	response := map[string]interface{}{
		"buffer_state":    session.BufferText(),
		"cursor_position": session.CursorIndex(),
		"current_mode":    session.Mode().String(),
		"match_status":    status.String(),
		"task_completed":  status == game.MatchComplete,
		"elapsed_time_ms": session.ElapsedTime().Milliseconds(),
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleKeystrokes(w http.ResponseWriter, r *http.Request, sessionID string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Keys []string `json:"keys"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	session := s.engine.GetSession(sessionID)
	if session == nil {
		writeError(w, http.StatusNotFound, "SESSION_NOT_FOUND", "Session not found")
		return
	}

	var status game.MatchStatus
	for _, key := range req.Keys {
		status = session.ProcessKey(key)
		if status == game.MatchComplete {
			break
		}
	}

	response := map[string]interface{}{
		"buffer_state":    session.BufferText(),
		"cursor_position": session.CursorIndex(),
		"current_mode":    session.Mode().String(),
		"match_status":    status.String(),
		"task_completed":  status == game.MatchComplete,
		"elapsed_time_ms": session.ElapsedTime().Milliseconds(),
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleCompleteTask(w http.ResponseWriter, r *http.Request, sessionID string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result, err := s.engine.CompleteTask(sessionID)
	if err != nil {
		writeError(w, http.StatusNotFound, "SESSION_NOT_FOUND", err.Error())
		return
	}

	session := s.engine.GetSession(sessionID)
	response := map[string]interface{}{
		"task_completed":  true,
		"round_complete":  session == nil || session.IsComplete(),
		"tasks_remaining": 0,
		"result":          result,
	}

	if session != nil && !session.IsComplete() {
		response["tasks_remaining"] = session.TotalTasks - session.CurrentIndex
		response["next_task"] = taskToMap(session.CurrentTask())
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleSkipTask(w http.ResponseWriter, r *http.Request, sessionID string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := s.engine.SkipTask(sessionID); err != nil {
		if err == game.ErrNoSkipsRemaining {
			writeError(w, http.StatusBadRequest, "NO_SKIPS_REMAINING", "No skips remaining")
		} else {
			writeError(w, http.StatusNotFound, "SESSION_NOT_FOUND", err.Error())
		}
		return
	}

	session := s.engine.GetSession(sessionID)
	response := map[string]interface{}{
		"task_skipped":    true,
		"tasks_remaining": session.TotalTasks - session.CurrentIndex,
	}

	if !session.IsComplete() {
		response["next_task"] = taskToMap(session.CurrentTask())
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleResetTask(w http.ResponseWriter, r *http.Request, sessionID string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := s.engine.ResetTask(sessionID); err != nil {
		writeError(w, http.StatusNotFound, "SESSION_NOT_FOUND", err.Error())
		return
	}

	session := s.engine.GetSession(sessionID)
	response := map[string]interface{}{
		"task_reset":      true,
		"buffer_state":    session.BufferText(),
		"cursor_position": session.CursorIndex(),
		"elapsed_time_ms": session.ElapsedTime().Milliseconds(),
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleSessionStats(w http.ResponseWriter, r *http.Request, sessionID string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := s.engine.GetSessionStats(sessionID)
	if err != nil {
		writeError(w, http.StatusNotFound, "SESSION_NOT_FOUND", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Return task list summary
	response := map[string]interface{}{
		"message": "Use /api/v1/tasks/:task_id to get a specific task",
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	taskID := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
	task := s.engine.GetTask(taskID)
	if task == nil {
		writeError(w, http.StatusNotFound, "TASK_NOT_FOUND", "Task not found")
		return
	}

	writeJSON(w, http.StatusOK, taskToMap(task))
}

func (s *Server) handleRounds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rounds := s.engine.GetRoundTypes()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"round_types": rounds,
	})
}

func (s *Server) handleRoundByType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roundType := strings.TrimPrefix(r.URL.Path, "/api/v1/rounds/")

	// Just return round info
	response := map[string]interface{}{
		"round_type": roundType,
		"task_count": 30,
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleLifetimeStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := s.engine.GetLifetimeStats()
	if stats == nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"total_rounds": 0,
			"total_tasks":  0,
		})
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

func (s *Server) handleStatsExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	stats := s.engine.GetLifetimeStats()

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=macaco-stats.json")
		json.NewEncoder(w).Encode(stats)
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=macaco-stats.csv")
		// Simple CSV export
		fmt.Fprintf(w, "Category,Tasks Completed,Avg Time (ms),Avg Efficiency\n")
		if stats != nil {
			for cat, cs := range stats.ByCategory {
				fmt.Fprintf(w, "%s,%d,%d,%.1f\n", cat, cs.TasksCompleted, cs.AvgTimeMs, cs.AvgEfficiency)
			}
		}
	default:
		writeError(w, http.StatusBadRequest, "INVALID_FORMAT", "Supported formats: json, csv")
	}
}

// Helper functions

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	})
}

func taskToMap(task *game.Task) map[string]interface{} {
	if task == nil {
		return nil
	}
	return map[string]interface{}{
		"task_id":      task.ID,
		"category":     task.Category,
		"difficulty":   task.Difficulty,
		"initial":      task.Initial,
		"desired":      task.Desired,
		"cursor_start": task.CursorStart,
		"cursor_end":   task.CursorEnd,
		"description":  task.Description,
		"hint":         task.Hint,
	}
}
