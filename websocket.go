package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

// https://github.com/otiai10/colorchat/blob/master/main.go

// The list of participants.
// FIXME: It's not necessary to be list.List, but it's easy for now :p
var participants = list.List{}

// Event represents an event from server side published.
type Event struct {
	Type string `json:"type"`
	Text string `json:"text"`
	User string `json:"user"`
}

func socket(conn *websocket.Conn) {

	// Add this request user to participants list.
	participation := participants.PushBack(conn)
	// FIXME: It's better to embed more information to participants list.
	//        For example, hmm..., yes, like "ID"

	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("#%02x%02x%02x", rand.Intn(255), rand.Intn(255), rand.Intn(255))
	logger := log.New(os.Stdout, fmt.Sprintf("[%s]\t", id), 0)

	defer func() {
		conn.Close()
		participants.Remove(participation) // clean up
		logger.Println("Exited loop")
	}()

	// Sturct for decoding message from client side
	msg := struct {
		Text string
		Type string
	}{}

	// {{{ FIXME: Tell who this request user is.
	ev := &Event{Type: "CONNECT", Text: "yourself", User: id}
	b, _ := json.Marshal(ev)
	conn.Write(b)
	// }}}

	// This loop keeps alive unless any error raises.
	for {

		if err := websocket.JSON.Receive(conn, &msg); err != nil {
			if err == io.EOF {
				logger.Println("Connection closed:", err)
			} else {
				logger.Println("Unexpected error:", err)
			}
			return // Exit from this loop
		}

		switch msg.Type {
		case "KEEPALIVE":
			// do nothing
		default:
			// event := &Event{
			// 	Type: "MESSAGE",
			// 	Text: msg.Text,
			// 	User: id,
			// }
			// b, _ := json.Marshal(event)
			// // Publish to all participants.
			// for e := participants.Front(); e != nil; e = e.Next() {
			// 	// FIXME: Type assertion validation :p
			// 	// FIXME: Error handling on Write :p
			// 	e.Value.(*websocket.Conn).Write(b)
			// }
			// FIXME: It's better to make some func to separate this common process :p
		}

		// continue: Keep waiting for message from this connection.
	}

}
