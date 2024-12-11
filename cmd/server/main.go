package main

import (
	"fmt"
	"game/common"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Error while reading message:", err)
		return
	}

	cmd := common.ParseMessage(msg)
	fmt.Println("Command:", cmd.Command, "Args:", cmd.Args)

	switch cmd.Command {
	case common.CommandJoin:
		if !g.Player1HasJoined {
			common.SendInitMessage(conn, 1)
			g.Player1HasJoined = true
			g.Player1.ID = cmd.Args[0]
			g.Player1.X = 0 + 20
			g.Player1.Y = 10
			g.Player1Connection = conn
		} else {
			common.SendInitMessage(conn, 2)
			g.Player2HasJoined = true
			g.Player2.ID = cmd.Args[0]
			g.Player2.X = 800 - (30)
			g.Player2.Y = 20
			g.Player2Connection = conn
			g.BallSpawned = true
			g.BallX = 800 / 2
			g.BallY = 600 / 2
		}

		go handleInputFromPlayer(conn)

		fmt.Println("Player Joined")
	default:
		fmt.Println("Invalid command:", cmd)
	}
}

// When input comes in from a websocket connection, it is sent to input channel
func handleInputFromPlayer(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			return
		}

		cmd := common.ParseMessage(msg)

		switch cmd.Command {
		case common.CommandQuit:
			fmt.Println("Player Quit")
			return
		case common.CommandInput:
			fmt.Println("Input:", cmd.Args[0], cmd.Args[1], cmd.Args[2])
			inchan <- cmd
		default:
			fmt.Println("Invalid command:", cmd)
		}
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/ws", handleWebsocket)
	go gameLoop(g)

	fmt.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
