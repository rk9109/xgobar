package main

import (
	"time"
	//"github.com/i3/go-i3"
)

// Rectangle representation
//
type Rectangle struct {
	x      int16
	y      int16
	width  uint16
	height uint16
	color  uint32
}

// Text representation
//
type Text struct {
	text  string
	font  string
	color uint32
}

// Block representation
//
type Block struct {
	text      Text
	rectangle Rectangle
}

// Module interface
//
type Module interface {
	run(ch chan []Block)
}

// Clock module
//
// Module outputs the current time. See https://golang.org/pkg/time/#Time.Format
// to customize output format
type Clock struct {
	x          int16
	y          int16
	width      uint16
	height     uint16
	font       string
	foreground uint32
	background uint32
	format     string
}

func (c Clock) run(ch chan []Block) {
	go func() {
		for {
			currentTime := time.Now().Format(c.format)

			block := []Block{
				Block{
					rectangle: Rectangle{
						x:      c.x,
						y:      c.y,
						width:  c.width,
						height: c.height,
						color:  c.background,
					},
					text: Text{
						text:  currentTime,
						font:  c.font,
						color: c.foreground,
					},
				},
			}
			ch <- block

			time.Sleep(time.Second)
		}
	}()
}

// Workspace module
//
type Workspace struct {
}

func (w Workspace) run(ch chan []Block) {
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
		}
	}()
}
