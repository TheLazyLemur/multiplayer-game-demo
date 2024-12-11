package main

import (
	"fmt"
	"game/common"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var g = &common.GameState{}
var p1MoveAction = ""
var p2MoveAction = ""
var hasJoined = false
var ballVelX = 1
var ballVelY = 1

var inchan = make(chan common.Message, 10)

func handleInput() {
	if g.Player1HasJoined {
		if g.Player1Connection == nil {
			fmt.Println("Connection is nil")
		} else {
			common.SendStateMessage(g.Player1Connection, *g)
		}
	}

	if g.Player2HasJoined {
		if g.Player2Connection == nil {
			fmt.Println("Connection is nil")
		} else {
			common.SendStateMessage(g.Player2Connection, *g)
		}
	}

	select {
	case cmd := <-inchan:
		if cmd.Args[1] == "press" {
			if cmd.Args[2] == g.Player1.ID {
				p1MoveAction = cmd.Args[0]
			}
			if cmd.Args[2] == g.Player2.ID {
				p2MoveAction = cmd.Args[0]
			}
		}

		if cmd.Args[1] == "release" {
			if cmd.Args[2] == g.Player1.ID {
				if cmd.Args[0] == p1MoveAction {
					p1MoveAction = ""
				}
			}
			if cmd.Args[2] == g.Player2.ID {
				if cmd.Args[0] == p2MoveAction {
					p2MoveAction = ""
				}
			}
		}
	default:
	}
}

func gameLoop(g *common.GameState) {
	for {
		handleInput()

		if g.BallSpawned {
			g.BallX += ballVelX
			g.BallY += ballVelY
			if g.BallX > 800 || g.BallX < 0 {
				ballVelX = -ballVelX
			}

			if g.BallY > 600 || g.BallY < 0 {
				ballVelY = -ballVelY
			}

			p1Rect := rl.Rectangle{X: float32(g.Player1.X), Y: float32(g.Player1.Y), Width: 10, Height: 100}
			p2Rect := rl.Rectangle{X: float32(g.Player2.X), Y: float32(g.Player2.Y), Width: 10, Height: 100}
			ballRect := rl.Rectangle{X: float32(g.BallX), Y: float32(g.BallY), Width: 10, Height: 100}

			if rl.CheckCollisionRecs(p1Rect, ballRect) {
				ballVelX = -ballVelX

				if p1MoveAction == "up" {
					ballVelY = -1

				}
				if p1MoveAction == "down" {
					ballVelY = 1
				}

				if p1MoveAction != "up" && p1MoveAction != "down" {
					ballVelY = -ballVelY
				}
			}

			if rl.CheckCollisionRecs(p2Rect, ballRect) {
				if p2MoveAction == "up" {
					ballVelY = -1
				}

				if p2MoveAction == "down" {
					ballVelY = 1
				}

				if p2MoveAction != "up" && p2MoveAction != "down" {
					ballVelY = -ballVelY
				}
			}
		}

		if p1MoveAction == "up" {
			g.Player1.Y -= 1
		}

		if p1MoveAction == "down" {
			g.Player1.Y += 1
		}

		if p2MoveAction == "up" {
			g.Player2.Y -= 1
		}

		if p2MoveAction == "down" {
			g.Player2.Y += 1
		}

		time.Sleep(time.Millisecond * 16)
	}
}
