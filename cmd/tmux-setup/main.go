package main

import (
	"flag"
	"log"
	"os"

	"github.com/bartosz-skejcik/tmux-setup/internal/config"
	"github.com/bartosz-skejcik/tmux-setup/internal/hooks"
	"github.com/bartosz-skejcik/tmux-setup/internal/tmux"
	"github.com/bartosz-skejcik/tmux-setup/internal/wizard"
)

func main() {
	args := os.Args

	flags := hooks.GetFlagsFromArgs(args)

	templateName := flags.Lookup("template").Value.(flag.Getter).Get().(string)
	createTemplate := flags.Lookup("create-template").Value.(flag.Getter).Get().(string)

	// Handle wizard with template creation
	if len(args) > 1 && args[1] == "wizard" {
		if createTemplate != "" {
			wizard.StartTemplateCreation(createTemplate)
			// wizard.Run(*createTemplate)
		} else {
			wizard.StartTemplateCreation("")
			// wizard.Run("")
		}
		return
	}

	var cfg config.Config
	var err error

	// If template is specified, load it directly
	if templateName != "" {
		cfg, err = config.LoadTemplate(templateName)
		if err != nil {
			log.Fatalf("Failed to load template: %v", err)
		}
	} else {
		// Otherwise look for local config file
		configPath := config.FindConfigFile()
		if configPath == "" {
			log.Fatal("No tmux.conf.yml found in current or parent directories.")
		}
		cfg, err = config.Load(configPath)
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}
	}

	// Perform dependency check
	if len(cfg.Dependencies) > 0 {
		tmux.CheckDependencies(cfg.Dependencies)
	}

	// Create the tmux session
	sessionName := cfg.SessionName
	if sessionName == "" {
		sessionName = "dev"
	}

	// Run pre-session hooks
	if err := hooks.RunPreSessionHooks(cfg); err != nil {
		log.Printf("Warning: pre-session hooks failed: %v", err)
	}

	err = tmux.CreateSession(sessionName, cfg)
	if err != nil {
		log.Fatalf("Failed to create tmux session: %v", err)
	}

	// Run post-session hooks
	if err := hooks.RunPostSessionHooks(cfg); err != nil {
		log.Printf("Warning: post-session hooks failed: %v", err)
	}

	// Attach to the tmux session if no template argument is provided
	if len(os.Args) == 1 || templateName != "" {
		err := tmux.AttachSession(sessionName, cfg.FocusWindow)
		if err != nil {
			log.Fatalf("Failed to attach to tmux session: %v", err)
		}
	}
}
