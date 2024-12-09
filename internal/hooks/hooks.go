package hooks

import (
	"fmt"
	"os/exec"

	"github.com/bartosz-skejcik/tmux-setup/internal/config"
)

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
