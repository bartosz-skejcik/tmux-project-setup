package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"gopkg.in/yaml.v3"
)

// Configuration struct for the YAML config file
type Config struct {
	SessionName  string         `yaml:"session_name"`
	FocusWindow  int            `yaml:"focus_window"`
	Defaults     GlobalDefaults `yaml:"defaults"`
	Dependencies []string       `yaml:"dependencies"`
	Windows      []WindowConfig `yaml:"windows"`
}

type GlobalDefaults struct {
	Directory      string `yaml:"directory"`
	InitialCommand string `yaml:"initial_command"`
}

type WindowConfig struct {
	Name           string       `yaml:"name"`
	Directory      string       `yaml:"directory"`
	InitialCommand string       `yaml:"initial_command"`
	Layout         string       `yaml:"layout"`
	GitBranch      string       `yaml:"git_branch"`
	Panes          []PaneConfig `yaml:"panes"`
}

type PaneConfig struct {
	Directory      string `yaml:"directory"`
	InitialCommand string `yaml:"initial_command"`
}

func main() {
	configPath := findConfigFile()
	if configPath == "" {
		log.Fatal("No tmux.conf.yml found in current or parent directories.")
	}

	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Perform dependency check
	if len(config.Dependencies) > 0 {
		checkDependencies(config.Dependencies)
	}

	// Create the tmux session
	sessionName := config.SessionName
	if sessionName == "" {
		sessionName = "dev"
	}
	err = createTmuxSession(sessionName, config)
	if err != nil {
		log.Fatalf("Failed to create tmux session: %v", err)
	}

	// Attach to the tmux session if no arguments are provided
	if len(os.Args) == 1 {
		err := attachToTmuxSession(sessionName, config.FocusWindow)
		if err != nil {
			log.Fatalf("Failed to attach to tmux session: %v", err)
		}
	}
}

// Locate the config file in current or parent directories
func findConfigFile() string {
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
func loadConfig(path string) (Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	return config, err
}

// Check for required dependencies
func checkDependencies(dependencies []string) {
	for _, dep := range dependencies {
		if _, err := exec.LookPath(dep); err != nil {
			log.Fatalf("Dependency missing: %s", dep)
		}
	}
}

// Create the tmux session
func createTmuxSession(sessionName string, config Config) error {
	exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-n", "placeholder").Run()
	exec.Command("tmux", "set-option", "-g", "base-index", "1").Run()
	exec.Command("tmux", "set-window-option", "-g", "pane-base-index", "1").Run()

	for i, window := range config.Windows {
		windowName := window.Name
		if windowName == "" {
			windowName = fmt.Sprintf("window-%d", i+1)
		}
		if i == 0 {
			exec.Command("tmux", "rename-window", "-t", fmt.Sprintf("%s:1", sessionName), windowName).Run()
		} else {
			exec.Command("tmux", "new-window", "-t", fmt.Sprintf("%s:%d", sessionName, i+1), "-n", windowName).Run()
		}

		// Apply Git integration
		if window.GitBranch != "" {
			runCommandInWindow(sessionName, i+1, fmt.Sprintf("git checkout %s", window.GitBranch))
		}

		// Set working directory
		dir := resolveDirectory(config.Defaults.Directory, window.Directory)
		if dir != "" {
			runCommandInWindow(sessionName, i+1, fmt.Sprintf("cd %s", dir))
		}

		// Run initial command
		initialCommand := window.InitialCommand
		if initialCommand == "" {
			initialCommand = config.Defaults.InitialCommand
		}
		if initialCommand != "" {
			runCommandInWindow(sessionName, i+1, initialCommand)
		}

		// Handle panes
		for j, pane := range window.Panes {
			if j > 0 {
				splitType := "-h"
				exec.Command("tmux", "split-window", splitType, "-t", fmt.Sprintf("%s:%d", sessionName, i+1)).Run()
			}
			paneDir := resolveDirectory(dir, pane.Directory)
			if paneDir != "" {
				runCommandInPane(sessionName, i+1, j+1, fmt.Sprintf("cd %s", paneDir))
			}
			if pane.InitialCommand != "" {
				runCommandInPane(sessionName, i+1, j+1, pane.InitialCommand)
			}
		}

		// Apply layout
		if window.Layout != "" {
			exec.Command("tmux", "select-layout", "-t", fmt.Sprintf("%s:%d", sessionName, i+1), window.Layout).Run()
		}
	}
	return nil
}

// Attach to the tmux session
func attachToTmuxSession(sessionName string, focusWindow int) error {
	if focusWindow == 0 {
		focusWindow = 1
	}

	tmuxPath, err := exec.LookPath("tmux")
	if err != nil {
		return err
	}

	exec.Command("tmux", "select-window", "-t", fmt.Sprintf("%s:%d", sessionName, focusWindow)).Run()

	args := []string{
		"tmux",
		"attach-session",
		"-t", sessionName,
	}

	return syscall.Exec(tmuxPath, args, os.Environ())
}

// Run a command in a specific tmux window
func runCommandInWindow(sessionName string, windowIndex int, command string) {
	exec.Command("tmux", "send-keys", "-t", fmt.Sprintf("%s:%d.1", sessionName, windowIndex), command, "C-m").Run()
}

// Run a command in a specific tmux pane
func runCommandInPane(sessionName string, windowIndex, paneIndex int, command string) {
	exec.Command("tmux", "send-keys", "-t", fmt.Sprintf("%s:%d.%d", sessionName, windowIndex, paneIndex), command, "C-m").Run()
}

// Resolve directory by prioritizing child over parent
func resolveDirectory(parent, child string) string {
	if child != "" {
		if filepath.IsAbs(child) {
			return child
		}
		return filepath.Join(parent, child)
	}
	return parent
}
