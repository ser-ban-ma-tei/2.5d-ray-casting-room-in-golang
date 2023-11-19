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

const (
	WINDOW_TITLE string = "Sewer City"
	WINDOW_WIDTH int32 = 1280
	WINDOW_HEIGHT int32 = 720

	MAP_WIDTH  int32   = 16
	MAP_HEIGHT int32   = 16
	MAP_DEPTH  float64 = 16.0

	WALL_IMAGE_WIDTH  int = 160
	WALL_IMAGE_HEIGHT int = 160
)

var (
	playerX   float64 = 8.0
	playerY   float64 = 8.0
	playerA   float64 = 0.0
	playerFOV float64 = math.Pi / 4.0

	mapRoom string

	keyPressedState = make(map[sdl.Keycode]bool)
)

type f2d struct {
	X, Y float64
}

type i2d struct {
	X, Y int32
}

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

	window, err := sdl.CreateWindow(
		WINDOW_TITLE,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		WINDOW_WIDTH,
		WINDOW_HEIGHT,
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

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB888, sdl.TEXTUREACCESS_STREAMING, WINDOW_WIDTH, WINDOW_HEIGHT)
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

	mouseSensitivity := 0.04
	sdl.SetRelativeMouseMode(true)
	defer sdl.SetRelativeMouseMode(false)

	isGameRunning := true

	lastTime := time.Now()
	lastTimeFps := lastTime
	var currentTime time.Time
	fpsCounter := 0
	currentFps := 0

	for isGameRunning {
		currentTime = time.Now()
		elapsedTime := currentTime.Sub(lastTime)
		elapsedTimeFps := currentTime.Sub(lastTimeFps)
		lastTime = currentTime
		var floatElapsedTime float64 = elapsedTime.Seconds()

		fpsCounter++
		if elapsedTimeFps.Milliseconds() > 1000 {
			currentFps = fpsCounter
			fpsCounter = 0
			lastTimeFps = currentTime

			window.SetTitle(fmt.Sprintf("%s - FPS: %d", WINDOW_TITLE, currentFps))
		}

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

		for x := int32(0); x < WINDOW_WIDTH; x++ {
			rayAngle := (playerA - playerFOV/2.0) + (float64(x)/float64(WINDOW_WIDTH))*playerFOV

			rayStart := f2d{X: playerX, Y: playerY}
			rayDirection := f2d{X: math.Cos(rayAngle), Y: math.Sin(rayAngle)}

			rayUnitStepSize := f2d{X: math.Sqrt(1 + math.Pow((rayDirection.Y/rayDirection.X), 2)), Y: math.Sqrt(1 + math.Pow((rayDirection.X/rayDirection.Y), 2))}
			mapCheck := i2d{X: int32(math.Trunc(rayStart.X)), Y: int32(math.Trunc(rayStart.Y))}
			rayLength := f2d{}
			step := i2d{}

			if rayDirection.X < 0 {
				step.X = -1
				rayLength.X = (rayStart.X - float64(mapCheck.X)) * rayUnitStepSize.X
			} else {
				step.X = 1
				rayLength.X = (float64(mapCheck.X+1) - rayStart.X) * rayUnitStepSize.X
			}

			if rayDirection.Y < 0 {
				step.Y = -1
				rayLength.Y = (rayStart.Y - float64(mapCheck.Y)) * rayUnitStepSize.Y
			} else {
				step.Y = 1
				rayLength.Y = (float64(mapCheck.Y+1) - rayStart.Y) * rayUnitStepSize.Y
			}

			isWallHit := false
			distanceToWall := 0.0
			wallTextureSampleX := 0.0

			for !isWallHit && distanceToWall < MAP_DEPTH {
				if rayLength.X < rayLength.Y {
					mapCheck.X += step.X
					distanceToWall = rayLength.X
					rayLength.X += rayUnitStepSize.X
				} else {
					mapCheck.Y += step.Y
					distanceToWall = rayLength.Y
					rayLength.Y += rayUnitStepSize.Y
				}

				if mapCheck.X < 0 || mapCheck.X >= MAP_WIDTH || mapCheck.Y < 0 || mapCheck.Y >= MAP_HEIGHT {
					isWallHit = true
					distanceToWall = MAP_DEPTH
				} else {
					if string(mapRoom[mapCheck.Y*MAP_WIDTH+mapCheck.X]) == "#" {
						isWallHit = true

						blockMidX := float64(mapCheck.X) + 0.5
						blockMidY := float64(mapCheck.Y) + 0.5

						testPointX := playerX + rayDirection.X*distanceToWall
						testPointY := playerY + rayDirection.Y*distanceToWall

						testAngle := math.Atan2(testPointY-blockMidY, testPointX-blockMidX)

						if testAngle >= -1*math.Pi*0.25 && testAngle < math.Pi*0.25 {
							wallTextureSampleX = testPointY - float64(mapCheck.Y)
						}
						if testAngle >= math.Pi*0.25 && testAngle < math.Pi*0.75 {
							wallTextureSampleX = testPointX - float64(mapCheck.X)
						}
						if testAngle < -1*math.Pi*0.25 && testAngle >= -1*math.Pi*0.75 {
							wallTextureSampleX = testPointX - float64(mapCheck.X)
						}
						if testAngle >= math.Pi*0.75 || testAngle < -1*math.Pi*0.75 {
							wallTextureSampleX = testPointY - float64(mapCheck.Y)
						}
					}
				}
			}

			ceiling := int32(float64(WINDOW_HEIGHT)/2.0 - float64(WINDOW_HEIGHT)/distanceToWall)
			floor := WINDOW_HEIGHT - ceiling

			ceilingColor := sdl.Color{R: 0, G: 0, B: 0, A: 255}

			for y := int32(0); y < WINDOW_HEIGHT; y++ {
				if y <= ceiling {
					setTexturePixel(pixels, x, y, WINDOW_WIDTH, ceilingColor)
				} else if y > ceiling && y <= floor {
					if distanceToWall < MAP_DEPTH {
						wallTextureSampleY := (float64(y) - float64(ceiling)) / (float64(floor) - float64(ceiling))
						pixelColor := sampleImageColor(wallPixels, wallTextureSampleX, wallTextureSampleY, WALL_IMAGE_WIDTH)
						pixelColor.A = 255
						setTexturePixel(pixels, x, y, WINDOW_WIDTH, pixelColor)
					} else {
						setTexturePixel(pixels, x, y, WINDOW_WIDTH, sdl.Color{R: 0, G: 0, B: 0, A: 255})
					}
				} else {
					setTexturePixel(pixels, x, y, WINDOW_WIDTH, sdl.Color{R: 0, G: 80, B: 0, A: 255})
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
