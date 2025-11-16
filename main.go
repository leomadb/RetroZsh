package main

import (
	"fmt"
	"math/rand/v2"
	"os"

	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
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

				// Draw the frame
				for x := 0; x < w; x++ {
					screen.SetContent(x, 0, '█', nil, tcell.StyleDefault)
					screen.SetContent(x, h-1, '█', nil, tcell.StyleDefault)
				}
				for y := 0; y < h; y++ {
					screen.SetContent(0, y, '█', nil, tcell.StyleDefault)
					screen.SetContent(w-1, y, '█', nil, tcell.StyleDefault)
				}
				showApple(screen)

				screen.Show()

				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	for {
		event := screen.PollEvent()
		switch event := event.(type) {
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				close(quit)
				return
			}
		}
	}
}

func showApple(screen tcell.Screen) {
	w, h := screen.Size()

	x := rand.IntN(w)
	y := rand.IntN(h)
	screen.SetContent(x, y, '🍎', nil, tcell.StyleDefault)

}
