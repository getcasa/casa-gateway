package gateway

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// WebsocketMessage define incoming message from casa server, including ActionMessage, etc...
type WebsocketMessage struct {
	Action string
	Body   []byte
}

// ActionMessage define incoming actions from casa server
type ActionMessage struct {
	Plugin string
	Call   string
	Params string
}

// WS is the connector to write message across app
var WS *websocket.Conn

// StartWebsocketClient start a websocket client to send and receive data between casa server and casa gateway
func StartWebsocketClient() {
	u := url.URL{Scheme: "ws", Host: "192.168.1.21:3000", Path: "/v1/ws"}
	log.Printf("connecting to %s", u.String())
	var err error
	WS, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	message := WebsocketMessage{
		Action: "newConnection",
		Body:   []byte(""),
	}
	byteMessage, _ := json.Marshal(message)
	err = WS.WriteMessage(websocket.TextMessage, byteMessage)
	if err != nil {
		log.Println("write:", err)
		return
	}

	// Handle response from server
	go func() {
		for {
			_, readMessage, err := WS.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				continue
			}

			var parsedMessage WebsocketMessage
			err = json.Unmarshal(readMessage, &parsedMessage)
			if err != nil {
				log.Println("read:", err)
				continue
			}

			switch parsedMessage.Action {
			case "hello":
				message := WebsocketMessage{
					Action: "newData",
					Body:   []byte("Hello from casa gateway!"),
				}
				byteMessage, _ := json.Marshal(message)
				err := WS.WriteMessage(websocket.TextMessage, byteMessage)
				if err != nil {
					log.Println("write:", err)
					return
				}
				break
			case "callAction":
				var action ActionMessage
				err = json.Unmarshal(parsedMessage.Body, &action)
				if err != nil {
					log.Println("read:", err)
					continue
				}

				PluginFromName(action.Plugin).CallAction(action.Call, []byte(action.Params))
				break
			default:
			}
		}
	}()
}
