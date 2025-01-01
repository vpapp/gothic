package main

import (
	"time"

	"github.com/vpapp/gothic"
)

func main() {
	ir := gothic.NewInterpreter(`
		pack [ttk::progressbar .bar1] -padx 20 -pady 20
		pack [ttk::progressbar .bar2] -padx 20 -pady 20
	`)

	go func() {
		i := 0
		inc := -1
		for {
			if i > 99 || i < 1 {
				inc = -inc
			}
			i += inc
			time.Sleep(50 * time.Millisecond)
			ir.Eval(`.bar1 configure -value %{}`, i)
		}
	}()

	go func() {
		i := 0
		inc := -1

		for {
			if i > 99 || i < 1 {
				inc = -inc
			}
			i += inc
			time.Sleep(100 * time.Millisecond)
			ir.Eval(`.bar2 configure -value %{}`, i)
		}
	}()

	<-ir.Done
}
