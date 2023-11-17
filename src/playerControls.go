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

	PLAYER_SPEED = 4.0
)

func updatePlayerPosition(floatElapsedTime float64) {
	speed := PLAYER_SPEED * floatElapsedTime
	cosA := math.Cos(playerA)
	sinA := math.Sin(playerA)

	if keyPressedState[sdl.K_w] {
		playerX += cosA * speed
		playerY += sinA * speed

		if string(mapRoom[int(playerY)*int(mapWidth)+int(playerX)]) == "#" {
			playerX -= cosA * speed
			playerY -= sinA * speed
		}
	}

	if keyPressedState[sdl.K_s] {
		playerX -= cosA * speed
		playerY -= sinA * speed

		if string(mapRoom[int(playerY)*int(mapWidth)+int(playerX)]) == "#" {
			playerX += cosA * speed
			playerY += sinA * speed
		}
	}

	if keyPressedState[sdl.K_a] {
		playerX += sinA * speed
		playerY -= cosA * speed

		if string(mapRoom[int(playerY)*int(mapWidth)+int(playerX)]) == "#" {
			playerX -= sinA * speed
			playerY += cosA * speed
		}
	}

	if keyPressedState[sdl.K_d] {
		playerX -= sinA * speed
		playerY += cosA * speed

		if string(mapRoom[int(playerY)*int(mapWidth)+int(playerX)]) == "#" {
			playerX += sinA * speed
			playerY -= cosA * speed
		}
	}
}
