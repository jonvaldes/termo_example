package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/jonvaldes/termo"
)

func main() {

	// Initialize termo
	if err := termo.Init(); err != nil {
		panic(err)
	}
	termo.EnableMouseEvents()
	termo.ShowCursor()

	// Reset terminal if we panic
	defer func() {
		termo.Stop()
		if err := recover(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	// Create termo framebuffer
	w, h, err := termo.Size()
	if err != nil {
		panic(err)
	}
	f := termo.NewFramebuffer(w, h)

	// Start reading input from stdin
	keyChan := make(chan termo.ScanCode, 100)
	go func() {
		for {
			s, err := termo.ReadScanCode()
			if err != nil {
				panic(err)
			}
			keyChan <- s
		}
	}()

	posx := 10
	posy := 10
	startT := time.Now().UnixNano()
	var lastScanCode termo.ScanCode

	// Main loop
	for {
		// Clear framebuffer
		f.Clear()

		// Draw outer frame
		f.ASCIIRect(0, 0, w, h, true, false)

		// Read keyboard
	readAgain:
		select {
		case s := <-keyChan:
			lastScanCode = s
			if s.IsMouseMoveEvent() {
				x, y := s.MouseCoords()
				termo.SetCursor(x, y)
			} else if s.IsMouseDownEvent() {
				x, y := s.MouseCoords()
				f.SetRect(x-2, y-2, 5, 5, termo.CellState{termo.AttrBold, termo.ColorYellow, termo.ColorYellow}, '#')
			} else if s.IsMouseUpEvent() {
				x, y := s.MouseCoords()
				f.SetRect(x-2, y-2, 5, 5, termo.CellState{termo.AttrBold, termo.ColorGreen, termo.ColorGreen}, '#')
			} else if s.IsEscapeCode() {
				switch s.EscapeCode() {
				case 65: // Up
					posy--
				case 66: // Down
					posy++
				case 67: // Right
					posx++
				case 68: // Left
					posx--
				default:
					//termo.Stop()
					//termo.PutText(fmt.Sprint(s.EscapeCode()))
					//return
				}
			} else {
				r := s.Rune()
				// Exit if Ctrl+C or Esc are pressed
				if r == 3 || r == 27 {
					termo.Stop()
					os.Exit(0)
				}
			}
			goto readAgain
		default:
		}

		// Draw the rectangle
		f.SetRect(posx, posy, 20, 20, termo.CellState{termo.AttrBold, termo.ColorYellow, termo.ColorRed}, '2')

		// Draw the sine wave
		t := float64(time.Now().UnixNano()-startT) / 500000.0

		chars := []rune{'.', 'o', '*', 'o', '.'}
		s := termo.CellState{termo.AttrNone, termo.ColorGreen, termo.ColorDefault}

		for i := 1; i < w-1; i++ {
			sh := 6 + int(5*math.Sin(0.001*t+float64(i)/float64(w)*math.Pi*2))
			for j := 1; j < 5; j++ {
				f.Set(i, sh-2+j, s, chars[j])
			}
		}

		// Draw text
		f.AttribRect(4, h-5, 10, 10, termo.StateDefault)
		f.SetText(4, h-5, fmt.Sprint(lastScanCode))
		f.SetText(4, h-4, "Press Up/Down/Left/Right to move")
		f.SetText(4, h-3, "Ctrl+C or Esc to exit")

		// Push framebuffer to screen
		f.Flush()

		// Wait for next frame
		time.Sleep(64 * time.Millisecond)
	}
}
