package main

import (
	"bytes"
	"log"
	"net"
	"os"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/game"
)

func main() {
	hostAddress := "localhost:8080"
	if len(os.Args) > 1 {
		hostAddress = os.Args[1]
	}

	addr, err := net.ResolveUDPAddr("udp", hostAddress)
	if err != nil {
		log.Fatalf("can't resolve addr: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("can't listen: %v", err)
	}
	defer conn.Close()

	gameServer := game.NewGameServer(conn)
	go func() {
		gameServer.Run()
	}()

	buffer := make([]byte, 1024)
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
