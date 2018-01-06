package main

import (
	"net/http"

	"log"

	"github.com/gorilla/websocket"
	"github.com/montera82/go-chat-app/pkg/trace"
	"github.com/stretchr/objx"
)

type room struct {
	//messages for this room
	forward chan *message

	//client wanting to join this room
	join chan *client

	//client wanting to leave this room
	leave chan *client

	//clients in this room
	clients map[*client]bool

	//tracer attribute to enable logging of room
	trace trace.Tracer

	avatar Avatar
}

//activates the room to listenForJoinsDepartureAndMessages for activities
//and aproprietly handle them
func (r *room) listenForJoinsDepartureAndMessages() {

	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.trace.Trace("client joined the room")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.trace.Trace("client left the room")
		case msg := <-r.forward:
			//forward message to all clients
			for client := range r.clients {
				client.send <- msg
				r.trace.Trace("Msg sent ", msg.Message)
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
	WriteBufferSize: messageBufferSize}

//handler which sets up a client pushes them to a room
func (r *room) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	socket, err := upgrader.Upgrade(w, rq, nil)
	if nil != err {
		log.Fatal("serveHTTP: ", err)
	}

	authCookie, err := rq.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
	}

	client := &client{
		socket:   socket,
		room:     r,
		send:     make(chan *message, messageBufferSize),
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.join <- client

	defer func() {
		r.leave <- client
	}()

	go client.write()                  //reads message from clients send channel, and writes it to the socket
	client.readMessagesFromWebSocket() //blocks the main routine
}

func newRoom(useAvatar Avatar) *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		trace:   trace.Off(),
		avatar:  useAvatar,
	}
}
