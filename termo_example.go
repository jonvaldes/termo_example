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
	keyChan := make(chan termo.ScanCode)
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

	// Main loop
	for {

		// Read keyboard
		select {
		case s := <-keyChan:
			if s.IsEscapeCode() {
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
					termo.Stop()
					fmt.Println(s.EscapeCode())
					return
				}
			}
			if !s.IsEscapeCode() {
				r := s.Rune()
				// Exit if Ctrl+C or Esc are pressed
				if r == 3 || r == 27 {
					os.Exit(0)
				}
			}
		default:
		}

		// Clear framebuffer
		f.Clear()

		// Draw the rectangle
		f.PutRect(posx, posy, 20, 20, termo.CellState{termo.AttrBold, termo.ColorYellow, termo.ColorRed}, '2')

		// Draw the sine wave
		t := float64(time.Now().UnixNano()-startT) / 500000.0

		chars := []rune{'.', 'o', '*', 'o', '.'}
		s := termo.CellState{termo.AttrNone, termo.ColorGreen, termo.ColorDefault}

		for i := 0; i < w; i++ {
			sh := 6 + int(5*math.Sin(0.001*t+float64(i)/float64(w)*math.Pi*2))
			for j := 0; j < 5; j++ {
				f.Put(i, sh-2+j, s, chars[j])
			}
		}

		// Draw text
		f.PutText(4, h-4, termo.CellState{termo.AttrNone, termo.ColorRed, termo.ColorGray}, "Press Up/Down/Left/Right to move")
		f.PutText(4, h-3, termo.CellState{termo.AttrBold, termo.ColorRed, termo.ColorDefault}, "Ctrl+C or Esc to exit")

		// Draw outer frame
		f.PutRect(0, 0, w, 1, termo.CellState{termo.AttrNone, termo.ColorDefault, termo.ColorDefault}, '-')
		f.PutRect(0, 0, 1, h, termo.CellState{termo.AttrNone, termo.ColorDefault, termo.ColorDefault}, '|')
		f.PutRect(w-1, 0, 1, h, termo.CellState{termo.AttrNone, termo.ColorDefault, termo.ColorDefault}, '|')
		f.PutRect(0, h-1, w, 1, termo.CellState{termo.AttrNone, termo.ColorDefault, termo.ColorDefault}, '-')

		// Push framebuffer to screen
		f.Flush()

		// Wait for next frame
		time.Sleep(64 * time.Millisecond)
	}
}
