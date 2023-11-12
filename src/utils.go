package main

import (
	"fmt"
	"os"
)

func getCurrentWorkingDirectory() string {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get current working directory: %s", err)
		os.Exit(1)
	}

	return cwd
}
