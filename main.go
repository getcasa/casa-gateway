package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ItsJimi/casa/devices"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	ip   = "224.0.0.50"
	port = "9898"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ip+":"+port)
	if err != nil {
		log.Panic(err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	fmt.Printf("Listening gateway events\n")

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Panic("Can't read udp", err)
		}

		var res Event
		err = json.Unmarshal(buf[0:n], &res)
		if err != nil {
			log.Println(err)
		}

		fmt.Println(res)

		// call ifttt webhook with switch
		if res.Model == "switch" {
			data := []byte(res.Data.(string))
			var button Switch
			err = json.Unmarshal(data, &button)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(button)
			req := []byte(`{"value1": "test", "value2": "bla", "value3": "toto"}`)
			if button.Status == "click" {
				resp, err := http.Post("https://maker.ifttt.com/trigger/button_pressed/with/key/"+os.Getenv("IFTTT_KEY"), "application/json", bytes.NewBuffer(req))
				if err != nil {
					log.Println(err)
				}
				defer resp.Body.Close()
			}
		}
	}
}
