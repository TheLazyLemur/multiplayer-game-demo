# Multiplayer Game Example

This is a multiplayer game example using Raylib and Go.
It uses a server authoritative game state.

## Architecture

The server is responsible for managing the game state and broadcasting it to all connected clients.
The client is responsible for handling user input and updating the game state.
Communication between the server and client is done using WebSockets. The client effectively acts as a controller for the game.

## How to run

### Run Server

```bash
go run ./cmd/server/...
```

### Run Client

```bash
go run ./cmd/client/...
```

Once 2 clients are connected, the ball will start to move, and players can move the paddle.
