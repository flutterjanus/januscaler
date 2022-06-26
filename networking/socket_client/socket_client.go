package networking

import (
	"net"
)

type SocketServerClient struct {
	conn *net.Conn
}

func (client *SocketServerClient) Send(msg string) {
	(*client.conn).Write([]byte(msg + "\n"))

}
