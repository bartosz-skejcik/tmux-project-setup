package hooks

import (
	"flag"
	"fmt"
	"log"
	"os/exec"

	"github.com/bartosz-skejcik/tmux-setup/internal/config"
)

func GetFlagsFromArgs(args []string) flag.FlagSet {
	flags := flag.NewFlagSet("tmux-setup", flag.ExitOnError)
	flags.String("template", "", "Use a template from ~/.config/tmux-setup/templates/")
	flags.String("create-template", "", "Create a new template using the wizard")

	for i, arg := range args {
		if arg == "--template" {
			if args[i+1] == "" {
				log.Fatalln("Please provide a template name")
			}
			flags.Set("template", args[i+1])
		}
		if arg == "--create-template" {
			if args[i+1] == "" {
				log.Fatalln("Please provide a name for the template")
			}
			flags.Set("create-template", args[i+1])
		}
	}

	return *flags
}

func RunPreSessionHooks(cfg config.Config) error {
	if cfg.Defaults.PreCommand != "" {
		if err := runCommand(cfg.Defaults.PreCommand); err != nil {
			return err
		}
	}
	return nil
}

func RunPostSessionHooks(cfg config.Config) error {
	if cfg.Defaults.PostCommand != "" {
		if err := runCommand(cfg.Defaults.PostCommand); err != nil {
			return err
		}
	}
	return nil
}

func RunPreWindowHooks(window config.WindowConfig) error {
	if window.PreCommand != "" {
		if err := runCommand(window.PreCommand); err != nil {
			return err
		}
	}
	return nil
}

func RunPostWindowHooks(window config.WindowConfig) error {
	if window.PostCommand != "" {
		if err := runCommand(window.PostCommand); err != nil {
			return err
		}
	}
	return nil
}

func runCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("hook command failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}
