package main

import (
	"fmt"
	"log"
	"net"
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
		fmt.Println(string(buf[0:n]))
	}
}
