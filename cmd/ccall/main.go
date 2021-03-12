package main

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	e "github.com/ryhszk/ccall/internal/err"
	"github.com/ryhszk/ccall/internal/ui"
)

func main() {
	if err := ui.NewApp(); err != nil {
		e.ErrExit(err.Error())
	}
}
