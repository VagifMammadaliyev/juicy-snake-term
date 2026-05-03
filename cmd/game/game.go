package main

import (
	"bufio"
	"log"
	"net"
	"os"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/game"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
	"golang.org/x/term"
)

func main() {
	// prepare terminal
	defer os.Stdin.Close()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	stdout := bufio.NewWriter(os.Stdout)
	defer stdout.Flush()

	terminal.Alternate.Set(stdout, terminal.Enable)
	defer terminal.Alternate.Set(stdout, terminal.Disable)

	terminal.HideCursor(stdout)
	defer terminal.RestoreCursor(stdout)
	// end prepare terminal

	// connect to server
	hostAddress := "localhost:8080"
	if len(os.Args) > 1 {
		hostAddress = os.Args[1]
	}

	addr, err := net.ResolveUDPAddr("udp", hostAddress)
	if err != nil {
		log.Fatalf("can't resolve server addr: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("can't connect to server: %v", err)
	}
	defer conn.Close()
	// end connect to server

	// TODO: looks like a good place for context
	done := make(chan struct{})
	game := game.NewGameClient(conn, stdout)
	go terminal.ListenControl(os.Stdin, game.Controls, done)
	game.Run(done)
}
