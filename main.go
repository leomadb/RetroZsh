package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer screen.Fini()

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				// Clear The Screen
				screen.Clear()

				// Dimensions
				w, h := screen.Size()

				// Draw the maze

			}
		}
	}()
}
