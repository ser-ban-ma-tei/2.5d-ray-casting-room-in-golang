package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	PLAYER_MOVE_FORWARD   uint8 = 0
	PLAYER_MOVE_BACKWARDS uint8 = 1
	PLAYER_MOVE_LEFT      uint8 = 2
	PLAYER_MOVE_RIGHT     uint8 = 3

	PLAYER_SPEED = 5.0
)

func updatePlayerPosition(floatElapsedTime float64) {
	speed := PLAYER_SPEED * floatElapsedTime

	if keyPressedState[sdl.K_w] {
		playerX += math.Sin(playerA) * speed
		playerY += math.Cos(playerA) * speed

		if string(mapRoom[int(playerY)*int(mapWidth)+int(playerX)]) == "#" {
			playerX -= math.Sin(playerA) * speed
			playerY -= math.Cos(playerA) * speed
		}
	}

	if keyPressedState[sdl.K_s] {
		playerX -= math.Sin(playerA) * speed
		playerY -= math.Cos(playerA) * speed

		if string(mapRoom[int(playerY)*int(mapWidth)+int(playerX)]) == "#" {
			playerX += math.Sin(playerA) * speed
			playerY += math.Cos(playerA) * speed
		}
	}

	if keyPressedState[sdl.K_a] {
		playerX -= math.Cos(playerA) * speed
		playerY += math.Sin(playerA) * speed

		if string(mapRoom[int(playerY)*int(mapWidth)+int(playerX)]) == "#" {
			playerX += math.Cos(playerA) * speed
			playerY -= math.Sin(playerA) * speed
		}
	}

	if keyPressedState[sdl.K_d] {
		playerX += math.Cos(playerA) * speed
		playerY -= math.Sin(playerA) * speed

		if string(mapRoom[int(playerY)*int(mapWidth)+int(playerX)]) == "#" {
			playerX -= math.Cos(playerA) * speed
			playerY += math.Sin(playerA) * speed
		}
	}
}
