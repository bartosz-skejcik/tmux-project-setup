// file: internal/wizard/wizard.go
package wizard

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bartosz-skejcik/tmux-setup/internal/config"
	"golang.org/x/term"
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

	cfg.SessionName = prompt("Enter session name", "dev")
	focusStr := prompt("Enter focus window number", "1")
	cfg.FocusWindow, _ = strconv.Atoi(focusStr)

	wantsToConfigureGlobalDefaults := prompt("Would you like to configure any global defaults? (yes, No)", "no")
	if wantsToConfigureGlobalDefaults == "yes" || wantsToConfigureGlobalDefaults == "y" {
		options := []string{
			"Directory",
			"Initial command",
			"Pre-command",
			"Post-command",
		}

		optionsToConfigure := MultiSelect("Which global default options would you like to configure?", options)

		if optionsToConfigure != nil {
			for _, option := range optionsToConfigure {
				switch option {
				case "Directory":
					cfg.Defaults.Directory = prompt("Default directory", "")
				case "Initial command":
					cfg.Defaults.InitialCommand = prompt("Default initial command", "")
				case "Pre-command":
					cfg.Defaults.PreCommand = prompt("Default pre-command (optional)", "")
				case "Post-command":
					cfg.Defaults.PostCommand = prompt("Default post-command (optional)", "")
				}
			}
		}
	}

	// Windows
	cfg.Windows = []config.WindowConfig{}
	for {
		fmt.Println("\nWindow Configuration:")
		wouldLikeToAddWindow := prompt("Would you like to add a window? (Yes/no)", "yes")
		if wouldLikeToAddWindow != "yes" && wouldLikeToAddWindow != "y" {
			break
		}

		// window := config.WindowConfig{
		// 	Name:           prompt("Window name (or empty to finish)", ""),
		// 	Directory:      prompt("Window directory", cfg.Defaults.Directory),
		// 	InitialCommand: prompt("Initial command", cfg.Defaults.InitialCommand),
		// 	GitBranch:      prompt("Git branch (optional)", ""),
		// 	PreCommand:     prompt("Pre-command (optional)", ""),
		// 	PostCommand:    prompt("Post-command (optional)", ""),
		// }

		options := []string{
			"Name",
			"Directory",
			"Initial command",
			"Git branch",
			"Pre-command",
			"Post-command",
		}

		window := config.WindowConfig{}

		optionsToConfigure := MultiSelect("Window configuration options available.", options)

		if optionsToConfigure != nil {
			for _, option := range optionsToConfigure {
				switch option {
				case "Name":
					window.Name = prompt("Window name", "")
				case "Directory":
					window.Directory = prompt("Window directory (optional)", cfg.Defaults.Directory)
				case "Initial command":
					window.InitialCommand = prompt("Initial command (optional)", cfg.Defaults.InitialCommand)
				case "Git branch":
					window.GitBranch = prompt("Git branch (optional)", "")
				case "Pre-command":
					window.PreCommand = prompt("Pre-command (optional)", "")
				case "Post-command":
					window.PostCommand = prompt("Post-command (optional)", "")
				}
			}

			// Panes
			window.Panes = []config.PaneConfig{}
		}

		for {
			fmt.Println("\nPane Configuration (for window: " + window.Name + "):")
			wouldLikeToAddPane := prompt("Would you like to add a pane? (Yes/no)", "yes")
			if wouldLikeToAddPane != "yes" && wouldLikeToAddPane != "y" {
				break
			}

			options := []string{
				"Directory",
				"Initial command",
				"Pre-command",
				"Post-command",
				"Refresh interval",
			}

			pane := config.PaneConfig{}

			optionsToConfigure := MultiSelect("Pane configuration options available.", options)

			if optionsToConfigure != nil {
				for _, option := range optionsToConfigure {
					switch option {
					case "Directory":
						pane.Directory = prompt("Pane directory (or empty to finish panes)", "")
					case "Initial command":
						pane.InitialCommand = prompt("Initial command", "")
					case "Pre-command":
						pane.PreCommand = prompt("Pre-command (optional)", "")
					case "Post-command":
						pane.PostCommand = prompt("Post-command (optional)", "")
					case "Refresh interval":
						refreshStr := prompt("Refresh interval in seconds (0 for no refresh)", "0")
						pane.RefreshInterval, _ = strconv.Atoi(refreshStr)
					}
				}
			}

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

// MultiSelect allows user to select multiple items from a list
func MultiSelect(title string, items []string) []string {
	// Disable input buffering and echo
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error setting terminal mode:", err)
		return nil
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Clear the screen
	fmt.Print("\033[2J")
	fmt.Print("\033[H")

	selected := make(map[int]bool)
	currentIndex := 0

	for {
		fmt.Print("\033[H")
		// Redraw the list
		fmt.Printf("%s\r\n\n", title)

		for i, item := range items {
			prefix := "  "
			if i == currentIndex {
				prefix = "> "
			}

			if selected[i] {
				prefix += "[x] "
			} else {
				prefix += "[ ] "
			}

			fmt.Printf("%s%s\r\n", prefix, item)
		}

		fmt.Printf("\r\n\nUse arrow keys/j/k to navigate, SPACE to select, ENTER to confirm, q to quit\r")

		// Read a single character
		b := make([]byte, 1)
		os.Stdin.Read(b)

		// Handle special keys and vim-style navigation
		switch {
		case b[0] == 'q' || b[0] == 3: // q or Ctrl+C
			// clear the screen
			fmt.Print("\033[2J")
			return nil
		case b[0] == '\r' || b[0] == '\n': // Enter
			// Collect selected items
			var result []string
			for i, item := range items {
				if selected[i] {
					result = append(result, item)
				}
			}
			// clear the screen
			fmt.Print("\033[2J")
			return result
		case b[0] == ' ': // Space to toggle selection
			selected[currentIndex] = !selected[currentIndex]
		case b[0] == 'j' || b[0] == 66: // Down arrow or 'j'
			if currentIndex < len(items)-1 {
				currentIndex++
			}
		case b[0] == 'k' || b[0] == 65: // Up arrow or 'k'
			if currentIndex > 0 {
				currentIndex--
			}
		}
	}
}
