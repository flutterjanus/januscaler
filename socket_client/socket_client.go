package socket_client

import (
	"net"
)

type SocketServerClient struct {
	conn *net.Conn
}

func (client *SocketServerClient) Send(msg string) {
	(*client.conn).Write([]byte(msg))
}
func (client *SocketServerClient) GetConn() net.Conn {
	return *client.conn
}
func (client *SocketServerClient) SetConn(conn *net.Conn) {
	client.conn = conn
}
