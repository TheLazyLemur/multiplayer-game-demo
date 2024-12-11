package main

import (
	"encoding/json"
	"fmt"
	"game/common"
	"log"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var gs common.GameState

func main() {

	id := uuid.New().String()

	url := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()

	fmt.Println("Connected to WebSocket server")

	err = common.SendJoinMessage(conn, id)
	if err != nil {
		log.Println("Error writing message:", err)
		return
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		return
	}

	cmd := common.ParseMessage(msg)
	if cmd.Command != common.CommandInit {
		log.Println("Unexpected message:", msg)
		return
	}
	fmt.Println(cmd)
	p, _ := strconv.Atoi(cmd.Args[0])
	fmt.Println("Player:", p)

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}

			cmd := common.ParseMessage(msg)
			if cmd.Command == common.CommandState {
				var ns common.GameState

				err = json.Unmarshal([]byte(cmd.Args[0]), &ns)
				if err != nil {
					log.Println("Error unmarshalling state:", err)
					return
				}

				gs = ns
			}
		}
	}()

	rl.InitWindow(int32(800), int32(600), "raylib [core] example - basic window")
	rl.SetTargetFPS(60)
	defer func() {
		rl.CloseWindow()

		err = common.SendQuitMessage(conn, id)
		if err != nil {
			log.Println("Error writing message:", err)
			return
		}

		conn.Close()
	}()

	for !rl.WindowShouldClose() {

		if rl.IsKeyPressed(rl.KeyUp) {
			common.SendInputMessage(conn, "up", "press", id)
		}
		if rl.IsKeyReleased(rl.KeyUp) {
			common.SendInputMessage(conn, "up", "release", id)
		}

		if rl.IsKeyPressed(rl.KeyDown) {
			common.SendInputMessage(conn, "down", "press", id)
		}
		if rl.IsKeyReleased(rl.KeyDown) {
			common.SendInputMessage(conn, "down", "release", id)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

		if gs.Player1HasJoined {
			rl.DrawRectangle(int32(gs.Player1.X), int32(gs.Player1.Y), 10, 100, rl.Red)
		}

		if gs.Player2HasJoined {
			rl.DrawRectangle(int32(gs.Player2.X), int32(gs.Player2.Y), 10, 100, rl.Blue)
		}

		if gs.BallSpawned {
			rl.DrawCircle(int32(gs.BallX), int32(gs.BallY), 10, rl.Green)
		}

		rl.EndDrawing()
	}
}
