package browser

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func (o *Opener) openNewWindow(rawURL string) error {
	command, err := defaultBrowserCommand()
	if err != nil {
		return openDetached("xdg-open", rawURL)
	}

	name := filepath.Base(command)

	switch {
	case isFirefox(name):
		return openDetached(command, "--new-window", rawURL)

	case isChromium(name):
		return openDetached(
			command,
			"--new-window",
			rawURL,
		)

	default:
		return openDetached("xdg-open", rawURL)
	}
}

func openDetached(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Process.Release()
}

func defaultBrowserCommand() (string, error) {
	if browser := strings.TrimSpace(os.Getenv("BROWSER")); browser != "" {
		return strings.Fields(browser)[0], nil
	}

	desktopID, err := output("xdg-settings", "get", "default-web-browser")
	if err != nil || desktopID == "" {
		desktopID, err = output("xdg-mime", "query", "default", "x-scheme-handler/https")
		if err != nil || desktopID == "" {
			return "", errors.New("default browser not found")
		}
	}

	return commandFromDesktopFile(desktopID)
}

func commandFromDesktopFile(desktopID string) (string, error) {
	home, _ := os.UserHomeDir()

	paths := []string{
		filepath.Join(home, ".local/share/applications", desktopID),
		filepath.Join("/usr/local/share/applications", desktopID),
		filepath.Join("/usr/share/applications", desktopID),
	}

	for _, path := range paths {
		command, err := execCommandFromDesktopFile(path)
		if err == nil && command != "" {
			return command, nil
		}
	}

	return "", errors.New("desktop file exec command not found")
}

func execCommandFromDesktopFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if !strings.HasPrefix(line, "Exec=") {
			continue
		}

		execline := strings.TrimPrefix(line, "Exec=")
		fields := strings.Fields(execline)
		if len(fields) == 0 {
			return "", errors.New("empty desktop exec command")
		}

		return fields[0], nil
	}

	return "", scanner.Err()
}

func output(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func isFirefox(name string) bool {
	switch name {
	case "firefox", "firefox-bin", "librewolf", "librewolf-bin", "zen", "zen-browser":
		return true
	default:
		return false
	}
}

func isChromium(name string) bool {
	switch name {
	case "brave",
		"brave-browser",
		"chromium",
		"chromium-browser",
		"google-chrome",
		"google-chrome-stable",
		"chrome",
		"vivaldi",
		"vivaldi-stable",
		"microsoft-edge",
		"microsoft-edge-stable":
		return true
	default:
		return false
	}
}
