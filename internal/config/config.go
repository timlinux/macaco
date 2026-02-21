package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	// Server settings
	ServerAddr string `json:"server_addr"`

	// Game settings
	AutoAdvanceDelay int  `json:"auto_advance_delay_ms"`
	ShowHints        bool `json:"show_hints"`
	EnableSounds     bool `json:"enable_sounds"`
	AnimationSpeed   float64 `json:"animation_speed"`

	// Appearance
	Theme    string `json:"theme"` // "dark", "light", "high-contrast"
	FontSize int    `json:"font_size"`

	// Data paths
	DataDir   string `json:"data_dir"`
	StatsFile string `json:"stats_file"`
	TasksFile string `json:"tasks_file"`
}

// Default returns the default configuration
func Default() *Config {
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".config", "macaco")

	return &Config{
		ServerAddr:       "localhost:8080",
		AutoAdvanceDelay: 500,
		ShowHints:        true,
		EnableSounds:     false,
		AnimationSpeed:   1.0,
		Theme:            "dark",
		FontSize:         1,
		DataDir:          dataDir,
		StatsFile:        filepath.Join(dataDir, "stats.json"),
		TasksFile:        "", // Use embedded tasks by default
	}
}

// Load loads configuration from a file
func Load(path string) (*Config, error) {
	cfg := Default()

	// If no path specified, try default locations
	if path == "" {
		homeDir, _ := os.UserHomeDir()
		defaultPaths := []string{
			filepath.Join(homeDir, ".config", "macaco", "config.json"),
			filepath.Join(homeDir, ".macaco.json"),
		}
		for _, p := range defaultPaths {
			if _, err := os.Stat(p); err == nil {
				path = p
				break
			}
		}
	}

	if path == "" {
		return cfg, nil // No config file found, use defaults
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Save saves the configuration to a file
func (c *Config) Save(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// EnsureDataDir creates the data directory if it doesn't exist
func (c *Config) EnsureDataDir() error {
	return os.MkdirAll(c.DataDir, 0755)
}
