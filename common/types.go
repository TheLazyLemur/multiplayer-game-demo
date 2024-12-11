package common

import (
	"encoding/json"
	"strconv"

	"github.com/gorilla/websocket"
)

// msg = Command:PlayerID

type Command string

const (
	CommandJoin  Command = "join"
	CommandQuit  Command = "quit"
	CommandInit  Command = "init"
	CommandState Command = "state"
	CommandInput Command = "input"
)

type Message struct {
	Command Command
	Args    []string
}

func (m Message) String() string {
	json, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return string(json)
}

func ParseMessage(msg []byte) Message {
	var m Message

	err := json.Unmarshal(msg, &m)
	if err != nil {
		panic(err)
	}

	return m
}

func SendJoinMessage(conn *websocket.Conn, playerID string) error {
	msg := Message{
		Command: CommandJoin,
		Args:    []string{playerID},
	}
	return conn.WriteMessage(websocket.TextMessage, []byte(msg.String()))
}

func SendInitMessage(conn *websocket.Conn, p int) error {
	msg := Message{
		Command: CommandInit,
		Args:    []string{strconv.Itoa(p)},
	}
	return conn.WriteMessage(websocket.TextMessage, []byte(msg.String()))
}

func SendStateMessage(conn *websocket.Conn, g GameState) error {
	gsJson, err := json.Marshal(g)
	if err != nil {
		panic(err)
	}
	msg := Message{
		Command: CommandState,
		Args:    []string{string(gsJson)},
	}
	return conn.WriteMessage(websocket.TextMessage, []byte(msg.String()))
}

func SendQuitMessage(conn *websocket.Conn, playerID string) error {
	msg := Message{
		Command: CommandQuit,
		Args:    []string{playerID},
	}
	return conn.WriteMessage(websocket.TextMessage, []byte(msg.String()))
}

func SendInputMessage(conn *websocket.Conn, key string, action string, playerID string) error {
	msg := Message{
		Command: CommandInput,
		Args:    []string{key, action, playerID},
	}
	return conn.WriteMessage(websocket.TextMessage, []byte(msg.String()))
}

type Player struct {
	ID string
	X  int
	Y  int
}

type GameState struct {
	Player1HasJoined  bool
	Player2HasJoined  bool
	Player1           Player
	Player2           Player
	Player1Connection *websocket.Conn
	Player2Connection *websocket.Conn
	BallSpawned       bool
	BallLaunched      bool
	BallX             int
	BallY             int
}
