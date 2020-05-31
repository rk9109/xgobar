package main

import (
	"time"

	"github.com/i3/go-i3"
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
	name      string
}

// Module interface
//
type Module interface {
	run(ch chan []Block)
}

// Time module
//
// Module outputs the current time. See https://golang.org/pkg/time/#Time.Format
// to customize output format
type Time struct {
	x          int16
	y          int16
	width      uint16
	height     uint16
	font       string
	foreground uint32
	background uint32
	format     string
}

func (t Time) run(ch chan []Block) {
	go func() {
		for {
			ch <- t.update()
			time.Sleep(time.Second)
		}
	}()
}

func (t Time) update() []Block {
	currentTime := time.Now().Format(t.format)

	return []Block{
		Block{
			rectangle: Rectangle{
				x:      t.x,
				y:      t.y,
				width:  t.width,
				height: t.height,
				color:  t.background,
			},
			text: Text{
				text:  currentTime,
				font:  t.font,
				color: t.foreground,
			},
			name: "time",
		},
	}
}

// CPU module
//
// Module outputs current CPU usage. CPU usage is calculated by polling
// the contents of /proc/stat.
type CPU struct {
	// public configuration
	x          int16
	y          int16
	width      uint16
	height     uint16
	font       string
	foreground uint32
	background uint32

	// private
	count uint64
}

func (c CPU) run(ch chan []Block) {
	go func() {
		for {
			ch <- c.update()
			time.Sleep(time.Second)
		}
	}()
}

func (c CPU) update() []Block {
	// TODO
	return []Block{}
}

// Workspace module
//
// Module outputs open workspaces and highlights the current active
// workspace.
type Workspace struct {
	x                  int16
	y                  int16
	width              uint16
	height             uint16
	font               string
	foreground         uint32
	backgroundActive   uint32
	backgroundInactive uint32
}

func (w Workspace) run(ch chan []Block) {
	go func() {
		for {
			ch <- w.update()
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (w Workspace) update() []Block {
	workspaces, err := i3.GetWorkspaces()
	if err != nil {
		return []Block{}
	}

	blocks := make([]Block, len(workspaces))
	for i, workspace := range workspaces {
		blocks[i] = Block{
			rectangle: Rectangle{
				x:      w.x,
				y:      w.y,
				width:  w.width,
				height: w.height,
				color:  w.backgroundInactive,
			},
			text: Text{
				text:  workspace.Name,
				font:  w.font,
				color: w.foreground,
			},
			name: "workspace",
		}
		if workspace.Visible {
			blocks[i].rectangle.color = w.backgroundActive
		}
		// increment position
		w.x += int16(w.width) + 5
	}
	return blocks
}
