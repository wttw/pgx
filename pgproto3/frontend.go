package pgproto3

import "net"

type Frontend struct {
	r io.Reader
	w io.Writer
}

func NewFrontend(r io.Reader, w io.Writer) (*Frontend, error) {
	return &Frontend{r: r, w: w}, nil
}

func (b *Frontend) Send(msg FrontendMessage) error {

}

func (b *Frontend) Receive() (BackendMessage, error) {

}
