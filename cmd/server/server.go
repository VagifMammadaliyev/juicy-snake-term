package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/config"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/logic"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/server"
)

func main() {
	conf := config.NewGameConfig()
	hostAddress := fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)
	if len(os.Args) > 1 {
		hostAddress = os.Args[1]
	}

	log.Printf("starting server at %s", hostAddress)

	addr, err := net.ResolveUDPAddr("udp", hostAddress)
	if err != nil {
		log.Fatalf("can't resolve addr: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("can't listen: %v", err)
	}
	defer conn.Close()

	logic := logic.NewLogicWithMaxArea()
	gameServer := server.NewGameServer(conn, logic)
	go gameServer.Run()

	buffer := make([]byte, 1024) // might be an overkill
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("can't read: %v", err)
			continue
		}

		dataBuffer := bytes.NewBuffer(bytes.Clone(buffer[:n]))
		go gameServer.HandlePlayerConnection(conn, dataBuffer, clientAddr)
	}

}
