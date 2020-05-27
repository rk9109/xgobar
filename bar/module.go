package main

import (
	"strconv"
	"time"
)

type Block struct {
	//
	X, Y int16

	//
	W, H uint16

	//
	Text string

	//
	Foreground uint32
	Background uint32
}

type Module interface {
	run(ch chan []Block)
}

type Test struct {
	//
	X, Y int16

	//
	W, H uint16

	//
	Foreground uint32
	Background uint32

	//
	Interval time.Duration
}

func (t Test) run(ch chan []Block) {
	go func() {
		counter := 0
		for {
			counter++
			blocks := []Block{
				Block{
					X:          t.X,
					Y:          t.Y,
					W:          t.W,
					H:          t.H,
					Text:       "counter: " + strconv.Itoa(counter),
					Foreground: t.Foreground,
					Background: t.Background,
				},
			}
			ch <- blocks
			time.Sleep(t.Interval)
		}
	}()
}
