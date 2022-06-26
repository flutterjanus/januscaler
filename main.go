package main

import (
	"flag"
	"fmt"
	. "github.com/flutterjanus/januscaler/socket_client"
	. "github.com/flutterjanus/januscaler/socket_server"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	// "os"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")
var mode = flag.String("mode", "server", "starting as server")
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type onMessageCallBack = func(msg string)
type Connection struct {
	response  http.ResponseWriter
	request   *http.Request
	conn      *websocket.Conn
	onMessage onMessageCallBack
	close     chan bool
}

func (conn *Connection) bindOnMessage(onMessage onMessageCallBack) {
	conn.onMessage = onMessage
}
func handleClient(connection *Connection) {
	var err error
	connection.conn, err = upgrader.Upgrade(connection.response, connection.request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer connection.conn.Close()
	for {
		_, message, err := connection.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if connection.onMessage != nil {
			connection.onMessage(string(message))
		}
		closeIt := <-connection.close
		if closeIt {
			break
		}
		// err = c.WriteMessage(1, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}

}

var connections map[string]*Connection = make(map[string]*Connection)

func home(response http.ResponseWriter, request *http.Request) {
	connections[request.RemoteAddr] = &Connection{request: request,
		close:    make(chan bool),
		response: response}
	connections[request.RemoteAddr].bindOnMessage(func(msg string) {
		fmt.Println(msg)
		connections[request.RemoteAddr].conn.WriteMessage(1, []byte(msg))
	})
	go handleClient(connections[request.RemoteAddr])
	fmt.Println(connections)
	// fmt.Println(networking.GetOutboundIP().String())
	<-connections[request.RemoteAddr].close
}
func main() {
	ss := MakeSocketServer()
	ss.Listen("tcp", "0.0.0.0:8081")
	ss.OnAuthorization(func(msg string, client *SocketServerClient) bool {
		return true
	})
	ss.OnMessage(func(msg string, client *SocketServerClient) {
		ss.Broadcast(fmt.Sprintf("client:%s sent:- %s", client.GetConn().RemoteAddr(), msg))
	})
	ss.AcceptConnections()

}
