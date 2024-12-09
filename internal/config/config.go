package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Configuration struct for the YAML config file
type Config struct {
	SessionName  string         `yaml:"session_name"`
	FocusWindow  int            `yaml:"focus_window"`
	Defaults     GlobalDefaults `yaml:"defaults"`
	Dependencies []string       `yaml:"dependencies"`
	Windows      []WindowConfig `yaml:"windows"`
	Template     string         `yaml:"template,omitempty"`
}

type GlobalDefaults struct {
	Directory      string `yaml:"directory"`
	InitialCommand string `yaml:"initial_command"`
	PreCommand     string `yaml:"pre_command,omitempty"`
	PostCommand    string `yaml:"post_command,omitempty"`
}

type WindowConfig struct {
	Name           string       `yaml:"name"`
	Directory      string       `yaml:"directory"`
	InitialCommand string       `yaml:"initial_command"`
	Layout         interface{}  `yaml:"layout"` // Can be string or LayoutConfig
	GitBranch      string       `yaml:"git_branch"`
	Panes          []PaneConfig `yaml:"panes"`
	PreCommand     string       `yaml:"pre_command,omitempty"`
	PostCommand    string       `yaml:"post_command,omitempty"`
}

type PaneConfig struct {
	Directory       string `yaml:"directory"`
	InitialCommand  string `yaml:"initial_command"`
	RefreshInterval int    `yaml:"refresh_interval,omitempty"`
	PreCommand      string `yaml:"pre_command,omitempty"`
	PostCommand     string `yaml:"post_command,omitempty"`
}

type LayoutConfig struct {
	Direction string       `yaml:"direction"`
	Panes     []PaneLayout `yaml:"panes"`
}

type PaneLayout struct {
	Width  string `yaml:"width,omitempty"`
	Height string `yaml:"height,omitempty"`
}

// FindConfigFile locates the config file in current or parent directories
func FindConfigFile() string {
	currentDir, _ := os.Getwd()
	for currentDir != "/" {
		configPath := filepath.Join(currentDir, "tmux.conf.yml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
		currentDir = filepath.Dir(currentDir)
	}
	return ""
}

// Load YAML configuration file
func Load(path string) (Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	// If template is specified, merge with template configuration
	if config.Template != "" {
		templateConfig, err := LoadTemplate(config.Template)
		if err != nil {
			return config, err
		}
		config = MergeConfigs(templateConfig, config)
	}

	return config, nil
}

// GetConfigDir returns the path to the configuration directory
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "tmux-setup")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}
