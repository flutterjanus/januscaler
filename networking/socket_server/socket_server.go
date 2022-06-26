package socket_server

import (
	"bufio"
	"fmt"
	""
	// "log"
	"net"
	// "os"
)

const (
	connHost = "0.0.0.0"
	connPort = "8081"
	connType = "tcp"
)

// func main() {
// 	ss := MakeSocketServer()
// 	ss.Listen(connType, "0.0.0.0:8081")
// 	ss.OnMessage(func(msg string, client *SocketServerClient) {
// 		fmt.Println(msg)
// 		client.Send("djfn")
// 	})
// 	ss.AcceptConnections()

// }

func handleClients(server *SocketServer, client *SocketServerClient) {
	for {
		buffer, err := bufio.NewReader(*(client).conn).ReadBytes('\n')
		if err != nil {
			// fmt.Println("Client left.")
			(*client.conn).Close()
			return
		}

		// log.Println("Client message:", string(buffer[:len(buffer)-1]))
		if server.callback != nil {
			server.callback(string(buffer[:len(buffer)-1]), client)
		}

		(*(client).conn).Write(buffer)
	}

}

type SocketServer struct {
	listener net.Listener
	callback func(msg string, client *SocketServerClient)
	clients  map[string]*SocketServerClient
}

func (server *SocketServer) OnMessage(callback func(msg string, client *SocketServerClient)) {

	server.callback = callback
}
func MakeSocketServer() *SocketServer {
	newServer := new(SocketServer)
	newServer.clients = make(map[string]*SocketServerClient)
	return newServer
}

func (server *SocketServer) Listen(connType string, connHost string) error {
	var err error
	server.listener, err = net.Listen(connType, connHost)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return err
	}
	fmt.Println("Starting " + connType + " server on " + connHost)
	return nil
}
func (server *SocketServer) acceptConnections() error {
	defer server.listener.Close()
	for {
		c, err := server.listener.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return err
		}

		client := &SocketServerClient{conn: &c}
		server.clients[c.RemoteAddr().String()] = client

		fmt.Println("Client connected.")
		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")
		go handleClients(server, client)

	}
}
func (server *SocketServer) AcceptConnections() {
	server.acceptConnections()
}
