package browser

import "os/exec"

type Opener struct{}

func NewOpener() *Opener {
	return &Opener{}
}

func (o *Opener) Open(rawURL string) error {

	cmd := exec.Command("xdg-open", rawURL)
	return cmd.Start()
}
