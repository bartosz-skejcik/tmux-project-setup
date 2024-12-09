package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// LoadTemplate loads a template configuration from the templates directory
func LoadTemplate(templateName string) (Config, error) {
	var config Config

	configDir, err := getConfigDir()
	if err != nil {
		return config, err
	}

	templatePath := filepath.Join(configDir, "templates", templateName+".yml")
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	return config, err
}

// MergeConfigs merges template config with user config, preferring user config values
func MergeConfigs(template, user Config) Config {
	result := template

	// Override with user values if they exist
	if user.SessionName != "" {
		result.SessionName = user.SessionName
	}
	if user.FocusWindow != 0 {
		result.FocusWindow = user.FocusWindow
	}

	// Merge windows
	if len(user.Windows) > 0 {
		result.Windows = user.Windows
	}

	// Merge defaults
	if user.Defaults != (GlobalDefaults{}) {
		result.Defaults = user.Defaults
	}

	return result
}

func getConfigDir() (string, error) {
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
