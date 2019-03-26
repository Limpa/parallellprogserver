package wshub

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// WSHub is a hub of websockets
type WSHub struct {
	sockets      []*websocket.Conn
	upgrader     websocket.Upgrader
	onConnection func() ([]byte, string)
}

// ConnectionHandler is a http handlerfunction which creates websocket connections
// It also executes provided callbackfunction when establishing a new connection
func (hub *WSHub) ConnectionHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := hub.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg, typ := hub.onConnection()
	err = conn.WriteMessage(websocket.TextMessage, buildSocketString(msg, typ))
	if err != nil {
		fmt.Println(err)
	}
	hub.sockets = append(hub.sockets, conn)
}

// Broadcast sends a message to all connected websockets
func (hub *WSHub) Broadcast(message []byte, typ string) {
	msg := buildSocketString(message, typ)
	for _, ws := range hub.sockets {
		ws.WriteMessage(websocket.TextMessage, msg)
	}
}

// New creates a new WSHub
func New(cb func() ([]byte, string)) *WSHub {
	return &WSHub{make([]*websocket.Conn, 0), websocket.Upgrader{}, cb}
}

func buildSocketString(msg []byte, typ string) []byte {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf(`{"type":"%s","content":`, typ))
	buffer.Write(msg)
	buffer.WriteString(`}`)
	return buffer.Bytes()
}
