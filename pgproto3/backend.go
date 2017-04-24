package pgproto3

import "net"

type Backend struct {
	r io.Reader
	w io.Writer
}

func NewBackend(r io.Reader, w io.Writer) (*Backend, error) {
	return &Backend{r: r, w: w}, nil
}

func (b *Backend) Send(msg BackendMessage) error {

}

func (b *Backend) Receive() (FrontendMessage, error) {

}
