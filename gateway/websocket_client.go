package gateway

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/ItsJimi/casa-gateway/logger"
	"github.com/ItsJimi/casa-gateway/utils"
	"github.com/gorilla/websocket"
)

// WebsocketMessage define incoming message from casa server, including ActionMessage, etc...
type WebsocketMessage struct {
	Action string
	Body   []byte
}

// ActionMessage define incoming actions from casa server
type ActionMessage struct {
	PhysicalID string
	Plugin     string
	Call       string
	Config     string
	Params     string
}

// WS is the connector to write message across app
var WS *websocket.Conn

// connectWebsocket connect casa gateway to casa server and retry on fails
func connectWebsocket(port string) {
	for {
		var err error
		ServerIP = utils.DiscoverServer()
		if ServerIP == "" {
			logger.WithFields(logger.Fields{}).Debugf("Wait 5 seconds to redail...")
			time.Sleep(time.Second * 5)
			continue
		}
		u := url.URL{Scheme: "ws", Host: ServerIP + ":" + utils.ServerPort, Path: "/v1/ws"}
		WS, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			logger.WithFields(logger.Fields{"code": "CGGWCCW001"}).Errorf("%s", err.Error())
			logger.WithFields(logger.Fields{}).Debugf("Wait 5 seconds to redail...")
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}

	addr := strings.Split(WS.LocalAddr().String(), ":")[0]
	message := WebsocketMessage{
		Action: "newConnection",
		Body:   []byte(addr + ":" + port),
	}
	byteMessage, _ := json.Marshal(message)
	err := WS.WriteMessage(websocket.TextMessage, byteMessage)
	if err != nil {
		logger.WithFields(logger.Fields{"code": "CGGWCCW002"}).Errorf("%s", err.Error())
		return
	}
	logger.WithFields(logger.Fields{}).Debugf("Websocket connected!")
}

// StartWebsocketClient start a websocket client to send and receive data between casa server and casa gateway
func StartWebsocketClient(port string) {
	connectWebsocket(port)

	// Handle response from server
	go func() {
		for {
			_, readMessage, err := WS.ReadMessage()
			if err != nil {
				logger.WithFields(logger.Fields{"code": "CGGWCSWC001"}).Errorf("%s", err.Error())

				// When read error is 1006, retry connection
				if strings.Contains(err.Error(), "close 1006") {
					connectWebsocket(port)
					continue
				}
				return
			}

			var parsedMessage WebsocketMessage
			err = json.Unmarshal(readMessage, &parsedMessage)
			if err != nil {
				logger.WithFields(logger.Fields{"code": "CGGWCSWC002"}).Errorf("%s", err.Error())
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
					logger.WithFields(logger.Fields{"code": "CGGWCSWC003"}).Errorf("%s", err.Error())
					return
				}
				break
			case "callAction":
				var action ActionMessage
				err = json.Unmarshal(parsedMessage.Body, &action)
				if err != nil {
					logger.WithFields(logger.Fields{"code": "CGGWCSWC004"}).Errorf("%s", err.Error())
					continue
				}

				PluginFromName(action.Plugin).CallAction(action.PhysicalID, action.Call, []byte(action.Params), []byte(action.Config))
				break
			default:
			}
		}
	}()
}
