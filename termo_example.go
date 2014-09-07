package main

import (
	"fmt"
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
	errChan := make(chan error)

	termo.StartKeyReadLoop(keyChan, errChan)

	ticker := time.Tick(500 * time.Millisecond)

	//var altPressed bool

	var lastScanCode termo.ScanCode
	// Main loop
	for {
		// Check for terminal resize
		if _w, _h, _ := termo.Size(); w != _w || h != _h {
			w = _w
			h = _h
			f = termo.NewFramebuffer(w, h)
		}

		// Clear framebuffer
		f.Clear()

		f.SetText(2, 0, " File  Edit  Search  View  Options  Help")
		f.AttribRect(0, 0, w, 1, termo.CellState{termo.AttrBold, termo.ColorBlack, termo.ColorGray})
		f.AttribRect(0, 1, w, h, termo.CellState{termo.AttrNone, termo.ColorGray, termo.ColorBlue})

		// Draw outer frame
		f.ASCIIRect(0, 1, w, h+4, false, false)
		filename := " C:\\CONFIG.SYS "
		f.SetText(w/2-len(filename)/2, 1, filename)
		f.AttribRect(w/2-len(filename)/2, 1, len(filename), 1, termo.CellState{termo.AttrBold, termo.ColorBlue, termo.ColorGray})

		// Draw text
		f.SetText(1, 2, "FILES=40\nBUFFERS=25\nDEVICE=C:\\HIMEM.SYS\nDEVICE=C:\\EMM386.EXE NOEMS\nDEVICEHIGH=C:\\NET\\ifshlp.sys")

		// Draw popup
		f.AttribRect(w/2-20, h/2-5, 40, 1, termo.BoldBlackOnWhite)
		f.AttribRect(w/2-19, h/2-4, 40, 10, termo.BoldWhiteOnBlack)
		f.AttribRect(w/2-20, h/2-4, 40, 9, termo.CellState{termo.AttrBold, termo.ColorBlack, termo.ColorGray})

		f.CenterText(w/2, h/2-5, "About termo_example.\n")
		f.CenterText(w/2, h/2-3, "Created by Jon Valdes\n"+
			"as an example for the termo\n"+
			"library, which you can get from:\n\n"+
			"https://github.com/jonvaldes/termo")

		f.AttribText(w/2-3, h/2+3, termo.BoldBlackOnWhite, "< OK >")
		f.SetRect(w/2-2, h/2+4, 6, 1, termo.CellState{termo.AttrNone, termo.ColorBlack, termo.ColorGray}, '▀')
		f.Set(w/2+3, h/2+3, termo.CellState{termo.AttrNone, termo.ColorBlack, termo.ColorGray}, '▄')

		f.SetText(1, 8, fmt.Sprint(lastScanCode))
		// Read keyboard
		select {
		case <-ticker:
			// Periodically flush framebuffer to screen
			f.Flush()
		case s := <-keyChan:
			lastScanCode = s
			if s.IsMouseMoveEvent() {
				x, y := s.MouseCoords()
				termo.SetCursor(x, y)
			} else if s.IsMouseDownEvent() {
				x, y := s.MouseCoords()
				f.SetRect(x-2, y-2, 5, 5, termo.CellState{termo.AttrBold, termo.ColorYellow, termo.ColorYellow}, '#')
				f.Flush()
			} else if s.IsMouseUpEvent() {
				x, y := s.MouseCoords()
				f.SetRect(x-2, y-2, 5, 5, termo.CellState{termo.AttrBold, termo.ColorGreen, termo.ColorGreen}, '#')
				f.Flush()
			} else if s.IsEscapeCode() {
				switch s.EscapeCode() {
				case 65: // Up
				case 66: // Down
				case 67: // Right
				case 68: // Left
				}
			} else {
				r := s.Rune()
				// Exit if Ctrl+C or Esc are pressed
				if r == 3 || r == 27 {
					termo.Stop()
					os.Exit(0)
				}
			}
		}
	}
}
