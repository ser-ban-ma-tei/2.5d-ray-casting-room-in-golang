package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func handleKeyboardEvent(t *sdl.KeyboardEvent, isGameRunning *bool) {
	if t.State == sdl.PRESSED && t.Keysym.Sym == sdl.K_ESCAPE {
		*isGameRunning = false
	}

	if t.Keysym.Sym == sdl.K_w {
		keyPressedState[sdl.K_w] = t.State == sdl.PRESSED
	}

	if t.Keysym.Sym == sdl.K_s {
		keyPressedState[sdl.K_s] = t.State == sdl.PRESSED
	}

	if t.Keysym.Sym == sdl.K_a {
		keyPressedState[sdl.K_a] = t.State == sdl.PRESSED
	}

	if t.Keysym.Sym == sdl.K_d {
		keyPressedState[sdl.K_d] = t.State == sdl.PRESSED
	}
}
