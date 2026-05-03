# Juicy snake

Multiplayer snake game to play with your coworker while you vibe-code.


## How to setup

Build the server:

```sh
go build ./cmd/server
```

Run the server on available port:

```sh
./server "0.0.0.0:5555"
```

Also find your local network IP address and send it to your coworker.


## How to connect

Build the game client:

```sh
go build ./cmd/game
```

Run the game client with server address you got from your coworker:

```sh
./game "XXX.XXX.X.XXX:5555"
```

This will connect you to the server, now if you press any of arrow keys, you will see your snake and you will be able to move it around.
Same process should be done for other players (building and launching the game client).

## How to play

If you hit wall or other snake, you lose and you can't control your snake. If you press any of arrow keys, you will spawn again.
Eat food and grow so you can easily make other snakes hit you.
Ctrl-C to exit the game.
