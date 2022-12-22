package ubuntu

import (
	"fmt"
	"os/exec"
)

func InitGoModule(moduleName string) error {
	cmd := exec.Command("go", "mod", "init", moduleName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("terminal: can't init golang module - %s", err)
	}
	return nil
}

func DownloadGoModuleDependencies() error {
	cmd := exec.Command("go", "mod", "download")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("terminal: can't download module dependencies - %s", err)
	}
	return nil
}

func TidyModuleDependencies() error {
	cmd := exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("terminal: can't tidy modeule dependencies - %s", err)
	}
	return nil
}
