package socket_server

import (
	"bufio"
	"fmt"
	"net"
	"time"

	. "github.com/flutterjanus/januscaler/socket_client"
)

// "log"

// "os"

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
// 		// ss.broadcast("gfg")

// 	})
// 	ss.AcceptConnections()

// }

func handleClients(server *SocketServer, client *SocketServerClient) {
	for {
		buffer, err := bufio.NewReader(client.GetConn()).ReadBytes('\n')
		if err != nil {
			// fmt.Println("Client left.")
			client.GetConn().Close()
			return
		}
		// log.Println("Client message:", string(buffer[:len(buffer)-1]))
		if server.callback != nil {
			server.callback(string(buffer[:len(buffer)-1]), client)
		}
		client.GetConn().Write(buffer)
	}

}

type SocketServer struct {
	listener              net.Listener
	callback              func(msg string, client *SocketServerClient)
	authorizationCallback func(msg string, client *SocketServerClient) bool
	Clients               map[string]*SocketServerClient
	guests                map[string]*SocketServerClient
}

func (s *SocketServer) Broadcast(msg string) {
	for _, c := range s.Clients {
		if c.GetConn() != nil {
			c.Send(msg)
		}
	}
}
func (server *SocketServer) OnMessage(callback func(msg string, client *SocketServerClient)) {
	server.callback = callback
}
func (server *SocketServer) OnAuthorization(callback func(msg string, client *SocketServerClient) bool) {
	server.authorizationCallback = callback
}
func MakeSocketServer() *SocketServer {
	newServer := new(SocketServer)
	newServer.Clients = make(map[string]*SocketServerClient)
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
		client := &SocketServerClient{}
		client.SetConn(&c)
		server.guests[c.RemoteAddr().String()] = client
		go server.manageGuests()
		// fmt.Println("Client connected.")
		// fmt.Println("Client " + c.RemoteAddr().String() + " connected.")
		go handleClients(server, client)

	}
}
func (server *SocketServer) manageGuests() {
	for _, client := range server.guests {
		client.GetConn().SetDeadline(time.Now().Add(time.Second * 60))
		go server.handleGuestAuthorization(client)
	}
}

func (server *SocketServer) AcceptConnections() {
	server.acceptConnections()
}
func (server *SocketServer) handleGuestAuthorization(client *SocketServerClient) {
	for {
		buffer, err := bufio.NewReader(client.GetConn()).ReadBytes('\n')
		msg := string(buffer[:len(buffer)-1])
		if err != nil {
			client.GetConn().Close()
			return
		}
		if server.authorizationCallback != nil {
			if server.authorizationCallback(msg, client) {
				client.GetConn().SetDeadline(time.Time{})
				server.Clients[msg] = client
			}
		}
	}

}
