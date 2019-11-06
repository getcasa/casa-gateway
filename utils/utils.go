package utils

import (
	cryptorand "crypto/rand"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	portscanner "github.com/anvie/port-scanner"
	"github.com/getcasa/sdk"
	"github.com/oklog/ulid/v2"
)

// NewULID create an ulid
func NewULID() ulid.ULID {
	id, _ := ulid.New(ulid.Timestamp(time.Now()), cryptorand.Reader)
	return id
}

// GetIDFile get ID from config file
func GetIDFile() string {
	var id string

	file, err := os.OpenFile(".casa", os.O_APPEND, 0644)
	if err != nil {
		id = string(NewULID().String())
		err = ioutil.WriteFile(".casa", []byte(id), 0644)
		Check(err, "error")
	} else {
		data := make([]byte, 100)
		count, err := file.Read(data)
		Check(err, "error")
		id = string(data[:count])
	}
	return id
}

// Check check error
func Check(e error, typ string) {
	if e != nil {
		panic(e)
	}
}

// FindTriggerFromName find trigger with name trigger
func FindTriggerFromName(triggers []sdk.Trigger, name string) sdk.Trigger {
	for _, trigger := range triggers {
		if trigger.Name == name {
			return trigger
		}
	}
	return sdk.Trigger{}
}

// FindFieldFromName find field with name field
func FindFieldFromName(fields []sdk.Field, name string) sdk.Field {
	for _, field := range fields {
		if field.Name == name {
			return field
		}
	}
	return sdk.Field{}
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

// ServerPort is the default port use by casa server
const ServerPort = 4353

// DiscoverServer get ips of casa servers
func DiscoverServer() []string {
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
					opened := ps.IsOpen(ServerPort)
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
