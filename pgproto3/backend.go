package pgproto3

import "net"

type Backend struct {
	conn net.Conn
}

func NewBackend() (*Backend, error) {

}

func (b *Backend) Send(msg BackendMessage) error {

}

func (b *Backend) Receive() (FrontendMessage, error) {

}
