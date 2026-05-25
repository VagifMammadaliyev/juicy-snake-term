# Juicy Snake

Multiplayer snake game to play directly in your terminal while you wait for agent output.

## Requirements

- Go 1.26 or later (only required to build from source)
- A terminal that supports ANSI escape sequences

## Build

Build the server and the game client:

```sh
go build ./cmd/server
go build ./cmd/game
```

This produces two binaries in the current directory: `server` and `game`.

## Run the server

Start the server on an available port:

```sh
./server "0.0.0.0:5555"
```

If no address is passed, the server uses `localhost:21000` by default. Find your local network IP address and share `IP:PORT` with the other players.

The server can also be configured via environment variables:

| Variable          | Default     | Description                |
| ----------------- | ----------- | -------------------------- |
| `JST_SERVER_HOST` | `localhost` | Host the server binds to   |
| `JST_SERVER_PORT` | `21000`     | UDP port the server uses   |

A command-line argument takes precedence over the environment variables.

## Connect

Run the game client with the address shared by the host:

```sh
./game "XXX.XXX.X.XXX:5555"
```

Press any arrow key to spawn your snake and start moving. Every player on the network runs the same command with the same server address.

## How to play

- Arrow keys to change snake direction.
- Eat food to grow.
- Colliding with bricks or other snakes kills you.
- `Ctrl-C` quits the game.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) and the [Code of Conduct](CODE_OF_CONDUCT.md).

## License

Released under the [MIT License](LICENSE).
