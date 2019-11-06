package gateway

import (
	"encoding/json"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	portscanner "github.com/anvie/port-scanner"
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

const port = 4353

// connectWebsocket connect casa gateway to casa server and retry on fails
func connectWebsocket() {
	ips := discoverServer()

	for {
		var err error
		for _, ip := range ips {
			u := url.URL{Scheme: "ws", Host: ip + ":" + strconv.Itoa(port), Path: "/v1/ws"}
			WS, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Printf("dial err:" + err.Error())
				log.Printf("wait 5 seconds to redail...")
				time.Sleep(time.Second * 5)
				continue
			}
		}
		break
	}

	message := WebsocketMessage{
		Action: "newConnection",
		Body:   []byte(""),
	}
	byteMessage, _ := json.Marshal(message)
	err := WS.WriteMessage(websocket.TextMessage, byteMessage)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func containIPAddress(arr []string, search string) bool {
	for _, addr := range arr {
		if addr == search {
			return true
		}
	}
	return false
}

var waitg sync.WaitGroup

func discoverServer() []string {
	var ips []string
	var ipAddresses []string
	ifaces, err := net.Interfaces()
	if err == nil {

		for _, iface := range ifaces {
			addrs, err := iface.Addrs()
			if err == nil {
				for _, addr := range addrs {
					cleanAddr := addr.String()[:strings.Index(addr.String(), "/")]
					if cleanAddr != "127.0.0.1" && !strings.Contains(cleanAddr, ":") && !net.ParseIP(cleanAddr).IsLoopback() {
						cleanAddr = addr.String()[:strings.LastIndex(addr.String(), ".")+1]
						if !containIPAddress(ipAddresses, cleanAddr) {
							ipAddresses = append(ipAddresses, cleanAddr)
						}
					}
				}
			}
		}

		waitg.Add(len(ipAddresses) * 255)

		for _, ipAddr := range ipAddresses {
			for i := 0; i < 255; i++ {
				go func(i int, ipAddr string) {
					ip := ipAddr + strconv.Itoa(i)
					ps := portscanner.NewPortScanner(ip, 3*time.Second, 4)
					opened := ps.IsOpen(port)
					if opened {
						ips = append(ips, ip)
					}
					waitg.Done()
				}(i, ipAddr)

			}
		}
		waitg.Wait()
	}

	return ips
}

// StartWebsocketClient start a websocket client to send and receive data between casa server and casa gateway
func StartWebsocketClient() {
	connectWebsocket()

	// Handle response from server
	go func() {
		for {
			_, readMessage, err := WS.ReadMessage()
			if err != nil {
				log.Println("read:", err)

				// When read error is 1006, retry connection
				if strings.Contains(err.Error(), "close 1006") {
					connectWebsocket()
					continue
				}
				return
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
