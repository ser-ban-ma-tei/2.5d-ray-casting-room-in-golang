package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/veandco/go-sdl2/img"
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

const (
	WALL_IMAGE_WIDTH  int = 160
	WALL_IMAGE_HEIGHT int = 160
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
	mapRoom += "#....##........#"
	mapRoom += "#....##........#"
	mapRoom += "#....##........#"
	mapRoom += "#..............#"
	mapRoom += "#..............#"
	mapRoom += "#..............#"
	mapRoom += "#.........######"
	mapRoom += "#..............#"
	mapRoom += "#..............#"
	mapRoom += "################"
}

func run() int {
	initRoom()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %s\n", err)
		return 1
	}
	defer sdl.Quit()

	windowWidth := int32(960)
	windowHeight := int32(540)

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

	wallImage, err := img.Load(filepath.Join(getCurrentWorkingDirectory(), "assets", "brick_wall.png"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load wall image: %s", err)
		return 1
	}
	defer wallImage.Free()

	wallPixels := wallImage.Pixels()

	isGameRunning := true

	tp1 := time.Now()
	var tp2 time.Time

	mouseSensitivity := 0.04
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

			sampleX := 0.0

			for !isWallHit && distanceToWall <= mapDepth {
				distanceToWall += 0.02

				testX := int(playerX + eyeX*distanceToWall)
				testY := int(playerY + eyeY*distanceToWall)

				if testX < 0 || testX >= int(mapWidth) || testY < 0 || testY >= int(mapHeight) {
					isWallHit = true
					distanceToWall = mapDepth
				} else {
					if string(mapRoom[testY*int(mapWidth)+testX]) == "#" {
						isWallHit = true

						blockMidX := float64(testX) + 0.5
						blockMidY := float64(testY) + 0.5

						testPointX := playerX + eyeX*distanceToWall
						testPointY := playerY + eyeY*distanceToWall

						testAngle := math.Atan2(testPointY-blockMidY, testPointX-blockMidX)

						if testAngle >= -1*math.Pi*0.25 && testAngle < math.Pi*0.25 {
							sampleX = testPointY - float64(testY)
						}
						if testAngle >= math.Pi*0.25 && testAngle < math.Pi*0.75 {
							sampleX = testPointX - float64(testX)
						}
						if testAngle < -1*math.Pi*0.25 && testAngle >= -1*math.Pi*0.75 {
							sampleX = testPointX - float64(testX)
						}
						if testAngle >= math.Pi*0.75 || testAngle < -1*math.Pi*0.75 {
							sampleX = testPointY - float64(testY)
						}

					}
				}
			}

			ceiling := int32(float64(windowHeight)/2.0 - float64(windowHeight)/distanceToWall)
			floor := windowHeight - ceiling

			ceilingColor := sdl.Color{R: 0, G: 0, B: 0, A: 255}

			for y := int32(0); y < windowHeight; y++ {
				if y <= ceiling {
					setTexturePixel(pixels, x, y, windowWidth, ceilingColor)
				} else if y > ceiling && y <= floor {
					if distanceToWall < mapDepth {
						sampleY := (float64(y) - float64(ceiling)) / (float64(floor) - float64(ceiling))
						pixelColor := sampleImageColor(wallPixels, sampleX, sampleY, WALL_IMAGE_WIDTH)
						setTexturePixel(pixels, x, y, windowWidth, pixelColor)
					} else {
						setTexturePixel(pixels, x, y, windowWidth, sdl.Color{R: 0, G: 0, B: 0, A: 255})
					}

				} else {
					b := 1.0 - (float64(y)-float64(windowHeight)/2.0)/(float64(windowHeight)/2.0)
					var floorColor sdl.Color
					if b < 0.25 {
						floorColor = sdl.Color{R: 0, G: 80, B: 0, A: 255}
					} else if b < 0.5 {
						floorColor = sdl.Color{R: 0, G: 64, B: 0, A: 255}
					} else if b < 0.75 {
						floorColor = sdl.Color{R: 0, G: 48, B: 0, A: 255}
					} else if b < 0.9 {
						floorColor = sdl.Color{R: 0, G: 32, B: 0, A: 255}
					} else {
						floorColor = sdl.Color{R: 0, G: 16, B: 0, A: 255}
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

func getPixelColor(pixels []byte, x, y, width int) sdl.Color {
	index := (y*width + x) * 4

	if index >= cap(pixels) {
		return sdl.Color{R: 0, G: 0, B: 0, A: 255}
	}

	return sdl.Color{
		R: pixels[index],
		G: pixels[index+1],
		B: pixels[index+2],
		A: pixels[index+3],
	}
}

func sampleImageColor(pixels []byte, x, y float64, width int) sdl.Color {
	sx := int(x * float64(width))
	sy := int(y*float64(width) - 1.0)

	if sx < 0 || sx >= width || sy < 0 || sy >= width {
		return sdl.Color{R: 0, G: 0, B: 0, A: 255}
	}

	return getPixelColor(pixels, sx, sy, width)
}

func main() {
	os.Exit(run())
}
