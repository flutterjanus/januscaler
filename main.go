package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")
var mode = flag.String("mode", "server", "starting as server")
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
        return true
    },
} 
type onMessageCallBack=func(msg string)
type Connection struct{
	response http.ResponseWriter
	request *http.Request
	conn *websocket.Conn
	onMessage onMessageCallBack
	close chan bool
}

func handleClient(connection *Connection){
	var err error
	connection.conn,err=upgrader.Upgrade(connection.response,connection.request,nil)
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
		connection.onMessage(string(message))
		closeIt:=<-connection.close
		if(closeIt){
			break;
		}
		// err = c.WriteMessage(1, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}

}
var connections map[string]*Connection= make(map[string]*Connection)
func home(response http.ResponseWriter,request *http.Request){
	connections[request.RemoteAddr]=&Connection{request: request,
		close: make(chan bool),
		response: response,onMessage:func(msg string) {
		fmt.Println(msg)
		connections[request.RemoteAddr].conn.WriteMessage(1,[]byte(msg))

	}};
	go handleClient(connections[request.RemoteAddr])
	<-connections[request.RemoteAddr].close
}
func main()  {
	// http.HandleFunc("/echo", echo)
	flag.NewFlagSet("foo", flag.ExitOnError)
	fmt.Println(os.Args[1]);
	fmt.Println(*mode)
	http.HandleFunc("/websocket", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}