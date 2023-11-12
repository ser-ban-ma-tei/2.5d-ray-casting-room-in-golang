package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	windowTitle string = "Sewer City"
)

var (
	playerX   float64 = 8.0
	playerY   float64 = 8.0
	playerA   float64 = 0.0
	playerFOV float64 = math.Pi / 4.0
)

var (
	mapWidth  byte    = 16
	mapHeight byte    = 16
	mapDepth  float64 = 16.0
)

var mapRoom string

var keyPressedState = make(map[sdl.Keycode]bool)

func initRoom() {
	mapRoom += "################"
	mapRoom += "#..............#"
	mapRoom += "#..............#"
	mapRoom += "#..............#"
	mapRoom += "#....##........#"
	mapRoom += "#....##........#"
	mapRoom += "#..............#"
	mapRoom += "#..............#"
	mapRoom += "#..............#"
	mapRoom += "#..............#"
	mapRoom += "#..............#"
	mapRoom += "#..............#"
	mapRoom += "#.........######"
	mapRoom += "#..............#"
	mapRoom += "#..............#"
	mapRoom += "################"
}

// var prevMouseX, prevMouseY int32

func run() int {
	initRoom()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %s\n", err)
		return 1
	}
	defer sdl.Quit()

	// displayBounds, err := sdl.GetDisplayBounds(0)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to get display bounds: %s\n", err)
	// 	return 1
	// }

	windowWidth := int32(800)
	windowHeight := int32(600)

	window, err := sdl.CreateWindow(
		windowTitle,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		windowWidth,
		windowHeight,
		sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a window: %s", err)
		return 1
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a renderer: %s", err)
	}
	defer renderer.Destroy()

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB888, sdl.TEXTUREACCESS_STREAMING, windowWidth, windowHeight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a texture: %s", err)
	}

	isGameRunning := true

	tp1 := time.Now()
	var tp2 time.Time

	mouseSensitivity := 0.05
	sdl.SetRelativeMouseMode(true)
	defer sdl.SetRelativeMouseMode(false)

	for isGameRunning {
		tp2 = time.Now()
		elapsedTime := tp2.Sub(tp1)
		tp1 = tp2
		var floatElapsedTime float64 = elapsedTime.Seconds()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				isGameRunning = false

			case *sdl.MouseMotionEvent:
				playerA += float64(t.XRel) * mouseSensitivity * floatElapsedTime

			case *sdl.KeyboardEvent:
				handleKeyboardEvent(t, &isGameRunning)
			}
		}

		updatePlayerPosition(floatElapsedTime)

		pixels, _, err := texture.Lock(nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to lock texture: %s", err)
			return 1
		}

		for x := int32(0); x < windowWidth; x++ {
			rayAngle := (playerA - playerFOV/2.0) + (float64(x)/float64(windowWidth))*playerFOV

			distanceToWall := 0.0
			isWallHit := false

			eyeX := math.Sin(rayAngle)
			eyeY := math.Cos(rayAngle)

			for !isWallHit && distanceToWall <= 16.0 {
				distanceToWall += 0.1

				testX := int(playerX + eyeX*distanceToWall)
				testY := int(playerY + eyeY*distanceToWall)

				if testX < 0 || testX >= int(mapWidth) || testY < 0 || testY >= int(mapHeight) {
					isWallHit = true
					distanceToWall = mapDepth
				} else {
					if string(mapRoom[testY*int(mapWidth)+testX]) == "#" {
						isWallHit = true
					}
				}
			}

			ceiling := int32(float64(windowHeight)/2.0 - float64(windowHeight)/distanceToWall)
			floor := windowHeight - ceiling

			wallShade := sdl.Color{R: 0, G: 0, B: 0, A: 255}

			if distanceToWall <= mapDepth/4.0 {
				wallShade = sdl.Color{R: 255, G: 255, B: 255, A: 255}
			} else if distanceToWall <= mapDepth/3.0 {
				wallShade = sdl.Color{R: 192, G: 192, B: 192, A: 192}
			} else if distanceToWall <= mapDepth/2.0 {
				wallShade = sdl.Color{R: 128, G: 128, B: 128, A: 128}
			} else if distanceToWall <= mapDepth/1.0 {
				wallShade = sdl.Color{R: 64, G: 64, B: 64, A: 64}
			}

			ceilingColor := sdl.Color{R: 0, G: 0, B: 0, A: 255}

			for y := int32(0); y < windowHeight; y++ {
				if y <= ceiling {
					setTexturePixel(pixels, x, y, windowWidth, ceilingColor)
				} else if y > ceiling && y <= floor {
					setTexturePixel(pixels, x, y, windowWidth, wallShade)
				} else {
					b := 1.0 - (float64(y)-float64(windowHeight)/2.0)/(float64(windowHeight)/2.0)
					var floorColor sdl.Color
					if b < 0.25 {
						floorColor = sdl.Color{R: 0, G: 0, B: 255, A: 255}
					} else if b < 0.5 {
						floorColor = sdl.Color{R: 0, G: 0, B: 224, A: 255}
					} else if b < 0.75 {
						floorColor = sdl.Color{R: 0, G: 0, B: 192, A: 255}
					} else if b < 0.9 {
						floorColor = sdl.Color{R: 0, G: 0, B: 160, A: 255}
					} else {
						floorColor = sdl.Color{R: 0, G: 0, B: 128, A: 255}
					}

					setTexturePixel(pixels, x, y, windowWidth, floorColor)
				}
			}

		}

		texture.Unlock()
		renderer.Copy(texture, nil, nil)
		renderer.Present()
		time.Sleep(time.Second / 60)
	}

	return 0
}

func setTexturePixel(pixels []byte, x, y, width int32, color sdl.Color) {
	index := (y*width + x) * 4
	pixels[index] = color.R
	pixels[index+1] = color.G
	pixels[index+2] = color.B
	pixels[index+3] = color.A
}

func main() {
	os.Exit(run())
}

// for y := int32(0); y < windowHeight; y++ {
// 	if y > floor {
// 		b := 1.0 - (float64(y)-float64(windowHeight)/2.0)/(float64(windowHeight)/2.0)
// 		if b < 0.25 {
// 			renderer.SetDrawColor(255, 0, 0, 255)
// 		} else if b < 0.5 {
// 			renderer.SetDrawColor(224, 0, 0, 255)
// 		} else if b < 0.75 {
// 			renderer.SetDrawColor(192, 0, 0, 255)
// 		} else if b < 0.9 {
// 			renderer.SetDrawColor(160, 0, 0, 255)
// 		} else {
// 			renderer.SetDrawColor(128, 0, 0, 255)
// 		}
// 		renderer.DrawPoint(x, y)
// 	}
// }
