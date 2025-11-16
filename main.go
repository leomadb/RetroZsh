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
		w, h := screen.Size()
		applePosition := []int{w / 2, h / 2}
		chars := []rune{'#'}
		snakeStyle := tcell.Style.Foreground(tcell.StyleDefault, tcell.ColorGreen).Background(tcell.ColorBlack)
		snakeBodyStyle := tcell.Style.Foreground(tcell.StyleDefault, tcell.ColorDarkGreen).Background(tcell.ColorBlack)
		appleStyle := tcell.Style.Foreground(tcell.StyleDefault, tcell.ColorRed).Background(tcell.ColorBlack)
		/*
			backdropStyle1 := tcell.Style.Background(tcell.StyleDefault, tcell.ColorGray)
			backdropStyle2 := tcell.Style.Background(tcell.StyleDefault, tcell.ColorBlack)
		*/

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
					if snakeHead[0] == 1 {
						snakeHead[0] = w - 2
					} else {
						snakeHead[0]--
					}
				case 2:
					if snakeHead[1] == h-2 {
						snakeHead[1] = 1
					} else {
						snakeHead[1]++
					}
				case 3:
					if snakeHead[0] == w-2 {
						snakeHead[0] = 2
					} else {
						snakeHead[0]++
					}
				case 4:
					if snakeHead[1] == 1 {
						snakeHead[1] = h - 2
					} else {
						snakeHead[1]--
					}
				}

				// Check for apple collision
				ateApple := snakeHead[0] == applePosition[0] && snakeHead[1] == applePosition[1]
				if ateApple {
					score++
					applePosition = showApple(screen, nil)
				}

				// Save tail position BEFORE moving (needed for expansion)
				var oldTailPos []int
				if len(snakeBody) > 0 {
					oldTailPos = []int{snakeBody[len(snakeBody)-1][0], snakeBody[len(snakeBody)-1][1]}
				}

				// Move body segments forward
				if len(snakeBody) > 0 {
					// Move all segments forward (from tail to head to avoid overwriting)
					for i := len(snakeBody) - 1; i > 0; i-- {
						snakeBody[i][0] = snakeBody[i-1][0]
						snakeBody[i][1] = snakeBody[i-1][1]
					}
					// Move first segment to where head was
					snakeBody[0][0] = lastHeadPosition[0]
					snakeBody[0][1] = lastHeadPosition[1]
				}

				// Add new segment if apple was eaten (at the OLD tail position)
				if ateApple {
					if len(snakeBody) == 0 {
						// If no body yet, add first segment at last head position
						snakeBody = append(snakeBody, []int{lastHeadPosition[0], lastHeadPosition[1]})
					} else {
						// Add new segment at the OLD tail position (before it moved)
						// This makes the tail stay in place, effectively growing the snake
						snakeBody = append(snakeBody, []int{oldTailPos[0], oldTailPos[1]})
					}
				}

				/*
					--------[ Handle Screen ]--------
				*/

				// Clear The Screen
				screen.Clear()

				// Draw the frame
				for x := 0; x < w; x++ {
					screen.SetContent(x, 0, '█', nil, tcell.StyleDefault)
					screen.SetContent(x, h-1, '█', nil, tcell.StyleDefault)
				}
				for y := 0; y < h; y++ {
					screen.SetContent(0, y, '█', nil, tcell.StyleDefault)
					screen.SetContent(w-1, y, '█', nil, tcell.StyleDefault)
				}

				// Score
				for i, r := range "score: " + strconv.Itoa(score) {
					screen.SetContent(1+i, 0, r, nil, tcell.StyleDefault)
				}

				// Backdrop
				/*
					for x := 0; x < w; x++ {
						for y := 0; y < h; y++ {
							if x == 0 || x == w-1 || y == 0 || y == h-1 {
								continue
							} else {
								if (x+y)%2 == 0 {
									screen.SetContent(x, y, ' ', nil, backdropStyle1)
								} else {
									screen.SetContent(x, y, ' ', nil, backdropStyle2)
								}
							}
						}
					}
				*/
				screen.SetContent(snakeHead[0], snakeHead[1], chars[0], nil, snakeStyle)
				for _, pos := range snakeBody {
					screen.SetContent(pos[0], pos[1], chars[0], nil, snakeBodyStyle)
				}

				// Draw the apple
				screen.SetContent(applePosition[0], applePosition[1], 'O', nil, appleStyle)

				screen.Show()

				time.Sleep(75 * time.Millisecond)
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
				direction = 4
			} else if event.Key() == tcell.KeyRight {
				direction = 3
			} else if event.Key() == tcell.KeyDown {
				direction = 2
			}
		}
	}
}

func showApple(screen tcell.Screen, pos []int) []int {
	w, h := screen.Size()

	var x, y int
	if pos != nil {
		x = pos[0]
		y = pos[1]
	} else {
		// Generate random position, avoiding the borders (frame)
		x = rand.IntN(w-4) + 2
		y = rand.IntN(h-4) + 2
	}
	return []int{x, y}
}
