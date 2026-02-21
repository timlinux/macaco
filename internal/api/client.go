package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/timlinux/macaco/internal/game"
	"github.com/timlinux/macaco/internal/stats"
)

// Client is an HTTP client for the MoCaCo API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient(serverAddr string) *Client {
	return &Client{
		baseURL: "http://" + serverAddr + "/api/v1",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// HealthCheck checks if the server is healthy
func (c *Client) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	return nil
}

// SessionResponse represents a session response from the API
type SessionResponse struct {
	SessionID        string                 `json:"session_id"`
	RoundType        string                 `json:"round_type"`
	TotalTasks       int                    `json:"total_tasks"`
	CurrentTaskIndex int                    `json:"current_task_index"`
	StartedAt        time.Time              `json:"started_at"`
	CurrentTask      map[string]interface{} `json:"current_task"`
	BufferState      string                 `json:"buffer_state"`
	CursorPosition   int                    `json:"cursor_position"`
	ElapsedTimeMs    int64                  `json:"elapsed_time_ms"`
}

// CreateSession creates a new game session
func (c *Client) CreateSession(roundType string) (*SessionResponse, error) {
	body := map[string]string{"round_type": roundType}
	jsonBody, _ := json.Marshal(body)

	resp, err := c.httpClient.Post(c.baseURL+"/sessions", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, c.parseError(resp)
	}

	var session SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, err
	}

	return &session, nil
}

// GetSession retrieves a session by ID
func (c *Client) GetSession(sessionID string) (*SessionResponse, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/sessions/" + sessionID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var session SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, err
	}

	return &session, nil
}

// DeleteSession deletes a session
func (c *Client) DeleteSession(sessionID string) error {
	req, err := http.NewRequest(http.MethodDelete, c.baseURL+"/sessions/"+sessionID, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// KeystrokeResponse represents a keystroke response from the API
type KeystrokeResponse struct {
	BufferState    string `json:"buffer_state"`
	CursorPosition int    `json:"cursor_position"`
	CurrentMode    string `json:"current_mode"`
	MatchStatus    string `json:"match_status"`
	TaskCompleted  bool   `json:"task_completed"`
	ElapsedTimeMs  int64  `json:"elapsed_time_ms"`
}

// SendKeystroke sends a single keystroke to the session
func (c *Client) SendKeystroke(sessionID, key string) (*KeystrokeResponse, error) {
	body := map[string]interface{}{
		"key":       key,
		"modifiers": []string{},
	}
	jsonBody, _ := json.Marshal(body)

	resp, err := c.httpClient.Post(
		c.baseURL+"/sessions/"+sessionID+"/keystroke",
		"application/json",
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var result KeystrokeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SendKeystrokes sends multiple keystrokes to the session
func (c *Client) SendKeystrokes(sessionID string, keys []string) (*KeystrokeResponse, error) {
	body := map[string]interface{}{
		"keys": keys,
	}
	jsonBody, _ := json.Marshal(body)

	resp, err := c.httpClient.Post(
		c.baseURL+"/sessions/"+sessionID+"/keystrokes",
		"application/json",
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var result KeystrokeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CompleteTaskResponse represents a task completion response
type CompleteTaskResponse struct {
	TaskCompleted  bool                    `json:"task_completed"`
	RoundComplete  bool                    `json:"round_complete"`
	TasksRemaining int                     `json:"tasks_remaining"`
	NextTask       map[string]interface{}  `json:"next_task,omitempty"`
	Result         *game.TaskResult        `json:"result,omitempty"`
}

// CompleteTask marks the current task as complete
func (c *Client) CompleteTask(sessionID string) (*CompleteTaskResponse, error) {
	resp, err := c.httpClient.Post(
		c.baseURL+"/sessions/"+sessionID+"/complete",
		"application/json",
		nil,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var result CompleteTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SkipTask skips the current task
func (c *Client) SkipTask(sessionID string) error {
	resp, err := c.httpClient.Post(
		c.baseURL+"/sessions/"+sessionID+"/skip",
		"application/json",
		nil,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.parseError(resp)
	}

	return nil
}

// ResetTask resets the current task
func (c *Client) ResetTask(sessionID string) (*KeystrokeResponse, error) {
	resp, err := c.httpClient.Post(
		c.baseURL+"/sessions/"+sessionID+"/reset",
		"application/json",
		nil,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var result KeystrokeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSessionStats retrieves statistics for a session
func (c *Client) GetSessionStats(sessionID string) (*stats.SessionStats, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/sessions/" + sessionID + "/stats")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var result stats.SessionStats
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetLifetimeStats retrieves lifetime statistics
func (c *Client) GetLifetimeStats() (*stats.LifetimeStats, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/stats/lifetime")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var result stats.LifetimeStats
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetRoundTypes retrieves available round types
func (c *Client) GetRoundTypes() ([]string, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/rounds")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var result struct {
		RoundTypes []string `json:"round_types"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.RoundTypes, nil
}

// parseError parses an error response
func (c *Client) parseError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	var errResp struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
		return fmt.Errorf("%s: %s", errResp.Error.Code, errResp.Error.Message)
	}

	return fmt.Errorf("API error: status %d", resp.StatusCode)
}
