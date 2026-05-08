package browser

import (
	"os/exec"
	"syscall"
)

type Opener struct{}

func NewOpener() *Opener {
	return &Opener{}
}

func (o *Opener) Open(rawURL string) error {
	return openDetached(rawURL)
}

func (o *Opener) OpenAndWait(rawURL string) error {
	return openDetached(rawURL)
}

func openDetached(rawURL string) error {
	cmd := exec.Command("xdg-open", rawURL)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Process.Release()
}
