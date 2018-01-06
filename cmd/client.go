package main

import (
	"time"

	"github.com/gorilla/websocket"
)

//client represents a single user on our platform
type client struct {
	socket   *websocket.Conn //the websocket being use to communicate by this client
	room     *room           //the room where client is chatting in
	send     chan *message   //client will use this to send a message
	userData map[string]interface{}
}

//reads messages off the socket associated to this client
func (c *client) readMessagesFromWebSocket() {
	defer c.socket.Close()
	//scan through the messages from websocket
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)

		if err != nil {
			return
		}

		msg.Time = time.Now()
		msg.Name = c.userData["name"].(string)
		//if avatarUrl, ok := c.userData["avatar_url"]; ok {
		//	msg.AvatarURL = avatarUrl.(string)
		//}

		msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
		//send the message to the rooms forward channel
		c.room.forward <- msg
	}
}

//check through the message available in send channel
//and write it to the socket
func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.send {

		err := c.socket.WriteJSON(msg)

		if nil != err {
			return
		}
	}
}
