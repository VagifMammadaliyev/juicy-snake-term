# Juicy Snake

Multiplayer snake game to play directly in your terminal while you wait for agent output.

## Requirements

- Go 1.26 or later (only required to build from source)
- A terminal that supports ANSI escape sequences

## Download

Pre-built binaries for Linux, macOS, and Windows are attached to each tagged release on the [Releases page](https://github.com/VagifMammadaliyev/juicy-snake-term/releases/latest). Linux and macOS builds ship as `.tar.gz`, Windows builds as `.zip`.

Pick the archive matching your OS and CPU architecture:

- `juicy-snake-server-<os>-<arch>.tar.gz` / `.zip` — the server
- `juicy-snake-game-<os>-<arch>.tar.gz` / `.zip` — the client

For example, on a 64-bit Linux machine you want `juicy-snake-server-linux-amd64.tar.gz` and `juicy-snake-game-linux-amd64.tar.gz`. On Apple Silicon, use `darwin-arm64`.

Extract the archives:

```sh
tar -xzf juicy-snake-server-linux-amd64.tar.gz
tar -xzf juicy-snake-game-linux-amd64.tar.gz
```

On Windows, right-click the `.zip` and choose "Extract All", or run `Expand-Archive` in PowerShell.

You'll get two binaries: `juicy-snake-server` and `juicy-snake-game` (with `.exe` suffix on Windows). The examples below use `./server` and `./game` for brevity — substitute the actual filename, or rename them.

## Build from source

Alternatively, build the server and the game client yourself:

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
