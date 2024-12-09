package tmux

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/bartosz-skejcik/tmux-setup/internal/config"
	"github.com/bartosz-skejcik/tmux-setup/internal/hooks"
)

// CheckDependencies verifies all required dependencies are available
func CheckDependencies(dependencies []string) {
	for _, dep := range dependencies {
		if _, err := exec.LookPath(dep); err != nil {
			log.Fatalf("Dependency missing: %s", dep)
		}
	}
}

// CreateSession creates a new tmux session with the given configuration
func CreateSession(sessionName string, cfg config.Config) error {
	exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-n", "placeholder").Run()
	exec.Command("tmux", "set-option", "-g", "base-index", "1").Run()
	exec.Command("tmux", "set-window-option", "-g", "pane-base-index", "1").Run()

	for i, window := range cfg.Windows {
		if err := createWindow(sessionName, i, window, cfg.Defaults); err != nil {
			return fmt.Errorf("failed to create window %d: %v", i+1, err)
		}
	}

	return nil
}

// AttachSession attaches to an existing tmux session
func AttachSession(sessionName string, focusWindow int) error {
	if focusWindow == 0 {
		focusWindow = 1
	}

	tmuxPath, err := exec.LookPath("tmux")
	if err != nil {
		return err
	}

	exec.Command("tmux", "select-window", "-t", fmt.Sprintf("%s:%d", sessionName, focusWindow)).Run()

	return syscall.Exec(tmuxPath, []string{"tmux", "attach-session", "-t", sessionName}, os.Environ())
}

func createWindow(sessionName string, index int, window config.WindowConfig, defaults config.GlobalDefaults) error {
	windowName := window.Name
	if windowName == "" {
		windowName = fmt.Sprintf("window-%d", index+1)
	}

	if err := hooks.RunPreWindowHooks(window); err != nil {
		return err
	}

	if index == 0 {
		exec.Command("tmux", "rename-window", "-t", fmt.Sprintf("%s:1", sessionName), windowName).Run()
	} else {
		exec.Command("tmux", "new-window", "-t", fmt.Sprintf("%s:%d", sessionName, index+1), "-n", windowName).Run()
	}

	// Set working directory
	dir := resolveDirectory(defaults.Directory, window.Directory)
	if dir != "" {
		sendKeys(sessionName, index+1, 1, fmt.Sprintf("cd %s", dir))
	}

	// Handle Git integration
	if window.GitBranch != "" {
		sendKeys(sessionName, index+1, 1, fmt.Sprintf("git checkout %s", window.GitBranch))
	}

	// Create panes and set up layouts
	if err := createPanes(sessionName, index+1, window, dir); err != nil {
		return err
	}

	if err := hooks.RunPostWindowHooks(window); err != nil {
		return err
	}

	return nil
}

func createPanes(sessionName string, windowIndex int, window config.WindowConfig, defaultDir string) error {
	for i, pane := range window.Panes {
		if i > 0 {
			splitType := "-h"
			if layout, ok := window.Layout.(string); ok && strings.Contains(layout, "vertical") {
				splitType = "-v"
			}
			exec.Command("tmux", "split-window", splitType, "-t", fmt.Sprintf("%s:%d", sessionName, windowIndex)).Run()
		}

		paneDir := resolveDirectory(defaultDir, pane.Directory)
		if paneDir != "" {
			sendKeys(sessionName, windowIndex, i+1, fmt.Sprintf("cd %s", paneDir))
		}

		if pane.InitialCommand != "" {
			sendKeys(sessionName, windowIndex, i+1, pane.InitialCommand)
		}

		if pane.RefreshInterval > 0 {
			go refreshPane(sessionName, windowIndex, i+1, pane)
		}
	}

	// Apply layout
	if window.Layout != nil {
		applyLayout(sessionName, windowIndex, window.Layout)
	}

	return nil
}

func refreshPane(sessionName string, windowIndex, paneIndex int, pane config.PaneConfig) {
	ticker := time.NewTicker(time.Duration(pane.RefreshInterval) * time.Second)
	for range ticker.C {
		sendKeys(sessionName, windowIndex, paneIndex, pane.InitialCommand)
	}
}

func applyLayout(sessionName string, windowIndex int, layout interface{}) {
	switch l := layout.(type) {
	case string:
		exec.Command("tmux", "select-layout", "-t", fmt.Sprintf("%s:%d", sessionName, windowIndex), l).Run()
	case config.LayoutConfig:
		// Apply custom layout using resize-pane commands
		for i, pane := range l.Panes {
			if pane.Width != "" {
				exec.Command("tmux", "resize-pane", "-t", fmt.Sprintf("%s:%d.%d", sessionName, windowIndex, i+1),
					"-x", strings.TrimSuffix(pane.Width, "%")).Run()
			}
			if pane.Height != "" {
				exec.Command("tmux", "resize-pane", "-t", fmt.Sprintf("%s:%d.%d", sessionName, windowIndex, i+1),
					"-y", strings.TrimSuffix(pane.Height, "%")).Run()
			}
		}
	}
}

func sendKeys(sessionName string, windowIndex, paneIndex int, command string) {
	exec.Command("tmux", "send-keys", "-t", fmt.Sprintf("%s:%d.%d", sessionName, windowIndex, paneIndex),
		command, "C-m").Run()
}

func resolveDirectory(parent, child string) string {
	if child != "" {
		if filepath.IsAbs(child) {
			return child
		}
		return filepath.Join(parent, child)
	}
	return parent
}
