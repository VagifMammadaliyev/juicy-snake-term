package main

import (
	"bufio"
	"os"

	"github.com/VagifMammadaliyev/juicy-snake-term/internal/game"
	"github.com/VagifMammadaliyev/juicy-snake-term/internal/terminal"
	"golang.org/x/term"
)

func main() {
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

	done := make(chan struct{})
	defer close(done)
	game := game.NewGame(stdout)
	go terminal.ListenControl(os.Stdin, game.Controls, done)
	game.Run()
}
