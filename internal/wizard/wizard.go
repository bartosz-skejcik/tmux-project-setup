package wizard

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bartosz-skejcik/tmux-setup/internal/config"
	"gopkg.in/yaml.v2"
)

func StartTemplateCreation(templateName string) {
	fmt.Printf("Creating template: %s\n", templateName)
	cfg := createConfig()
	saveTemplate(cfg, templateName)
}

func Start() {
	cfg := createConfig()
	saveConfig(cfg)
}

func createConfig() config.Config {
	cfg := config.Config{}

	// Session name
	cfg.SessionName = prompt("Enter session name", "dev")

	// Focus window
	focusStr := prompt("Enter focus window number", "1")
	cfg.FocusWindow, _ = strconv.Atoi(focusStr)

	// Global defaults
	fmt.Println("\nGlobal Defaults:")
	cfg.Defaults = config.GlobalDefaults{
		Directory:      prompt("Default directory", ""),
		InitialCommand: prompt("Default initial command", ""),
		PreCommand:     prompt("Default pre-command (optional)", ""),
		PostCommand:    prompt("Default post-command (optional)", ""),
	}

	// Windows
	cfg.Windows = []config.WindowConfig{}
	for {
		fmt.Println("\nWindow Configuration:")
		windowName := prompt("Window name (or empty to finish)", "")
		if windowName == "" {
			break
		}

		window := config.WindowConfig{
			Name:           windowName,
			Directory:      prompt("Window directory", cfg.Defaults.Directory),
			InitialCommand: prompt("Initial command", cfg.Defaults.InitialCommand),
			GitBranch:      prompt("Git branch (optional)", ""),
			PreCommand:     prompt("Pre-command (optional)", ""),
			PostCommand:    prompt("Post-command (optional)", ""),
		}

		// Panes
		window.Panes = []config.PaneConfig{}
		for {
			fmt.Println("\nPane Configuration (for window: " + window.Name + "):")
			paneDir := prompt("Pane directory (or empty to finish panes)", "")
			if paneDir == "" {
				break
			}

			pane := config.PaneConfig{
				Directory:      paneDir,
				InitialCommand: prompt("Initial command", ""),
				PreCommand:     prompt("Pre-command (optional)", ""),
				PostCommand:    prompt("Post-command (optional)", ""),
			}

			refreshStr := prompt("Refresh interval in seconds (0 for no refresh)", "0")
			pane.RefreshInterval, _ = strconv.Atoi(refreshStr)

			window.Panes = append(window.Panes, pane)
		}

		// Layout
		if len(window.Panes) > 1 {
			layoutType := prompt("Layout type (simple/advanced)", "simple")
			if layoutType == "simple" {
				window.Layout = prompt("Layout (even-horizontal/even-vertical/main-horizontal/main-vertical)", "even-horizontal")
			} else {
				layout := config.LayoutConfig{
					Direction: prompt("Layout direction (horizontal/vertical)", "horizontal"),
					Panes:     make([]config.PaneLayout, len(window.Panes)),
				}

				for i := range window.Panes {
					fmt.Printf("\nPane %d size:\n", i+1)
					layout.Panes[i] = config.PaneLayout{
						Width:  prompt("Width percentage (e.g., 30%)", ""),
						Height: prompt("Height percentage (e.g., 50%)", ""),
					}
				}
				window.Layout = layout
			}
		}

		cfg.Windows = append(cfg.Windows, window)
	}

	return cfg
}

func saveTemplate(cfg config.Config, templateName string) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		fmt.Printf("Error getting config directory: %v\n", err)
		return
	}

	templatesDir := filepath.Join(configDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		fmt.Printf("Error creating templates directory: %v\n", err)
		return
	}

	filename := filepath.Join(templatesDir, templateName+".yml")
	data, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Printf("Error marshaling configuration: %v\n", err)
		return
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("Error saving template: %v\n", err)
		return
	}

	fmt.Printf("Template saved to %s\n", filename)
}

func prompt(message, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", message, defaultValue)
	} else {
		fmt.Printf("%s: ", message)
	}

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}
	return input
}

func saveConfig(cfg config.Config) {
	filename := prompt("Save configuration as", "tmux.conf.yml")
	if !strings.HasSuffix(filename, ".yml") {
		filename += ".yml"
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Printf("Error marshaling configuration: %v\n", err)
		return
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("Error saving configuration: %v\n", err)
		return
	}

	fmt.Printf("Configuration saved to %s\n", filename)
}
