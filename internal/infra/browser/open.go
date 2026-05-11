package browser

type Opener struct{}

func NewOpener() *Opener {
	return &Opener{}
}

func (o *Opener) Open(rawURL string) error {
	return openDetached("xdg-open", rawURL)
}

func (o *Opener) OpenNewWindow(rawURL string) error {
	return o.openNewWindow(rawURL)
}

func (o *Opener) OpenAndWait(rawURL string) error {
	return openDetached("xdg-open", rawURL)
}
