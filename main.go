package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"

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

	direction := 1

	go func() {
		snakeHead := []int{10, 10}
		var score int = 0
		var snakeBody [][]int
		lastHeadPosition := make([]int, 2)
		applePosition := make([]int, 2)
		chars := []rune{'#'}

		showApple(screen, applePosition)
		snakeStyle := tcell.Style.Foreground(tcell.StyleDefault, tcell.ColorGreen).Background(tcell.ColorBlack)

		for {
			select {
			case <-quit:
				return
			default:
				/*
					--------[ Game Logic ]--------
				*/

				lastHeadPosition[0] = snakeHead[0]
				lastHeadPosition[1] = snakeHead[1]

				switch direction {
				case 1:
					snakeHead[0]--
				case 2:
					snakeHead[1]++
				case 3:
					snakeHead[0]++
				case 4:
					snakeHead[1]--
				}

				if snakeHead[0] == applePosition[0] && snakeHead[1] == applePosition[1] {
					eatApple(screen, applePosition, score)
					snakeBody = append(snakeBody, lastHeadPosition)
				}

				for i := range snakeBody {
					if i == 0 {
						snakeBody[i][0] = lastHeadPosition[0]
						snakeBody[i][1] = lastHeadPosition[1]
					} else {
						snakeBody[i][0] = snakeBody[i-1][0]
						snakeBody[i][1] = snakeBody[i-1][1]
					}
				}

				/*
					--------[ Handle Screen ]--------
				*/

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
				for i, r := range strconv.Itoa(score) {
					screen.SetContent(w-2-i, 0, r, nil, tcell.StyleDefault)
				}

				screen.SetContent(snakeHead[0], snakeHead[1], chars[0], nil, snakeStyle)
				for _, pos := range snakeBody {
					screen.SetContent(pos[0], pos[1], chars[0], nil, snakeStyle)
				}

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

			// Handle Directions
			if event.Key() == tcell.KeyLeft {
				direction = 1
			} else if event.Key() == tcell.KeyUp {
				direction = 2
			} else if event.Key() == tcell.KeyRight {
				direction = 3
			} else if event.Key() == tcell.KeyDown {
				direction = 4
			}
		}
	}
}

func showApple(screen tcell.Screen, pos []int) {
	w, h := screen.Size()

	x := rand.IntN(w)
	y := rand.IntN(h)
	if pos != nil {
		x = pos[0]
		y = pos[1]
	}
	screen.SetContent(x, y, '🍎', nil, tcell.StyleDefault)

}

func eatApple(screen tcell.Screen, pos []int, score int) {
	score++
	showApple(screen, pos)
}
