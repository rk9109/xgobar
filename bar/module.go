package main

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/i3/go-i3"
)

type Text struct {
	text  string
	font  string
	color uint32
}

type Underline struct {
	height uint16
	color  uint32
}

// Rectangle representation
//
type Rectangle struct {
	x      int16
	y      int16
	width  uint16
	height uint16
	color  uint32
}

// Block representation
//
type Block struct {
	text      Text
	underline Underline
	rectangle Rectangle
	name      string
}

// Module interface
//
type Module interface {
	run(ch chan []Block)
}

// Plaintext module
//
// Module outputs a single constant string. Used to construct static elements in the
// bar (e.g. icons).
type Plaintext struct {
	x          int16
	y          int16
	width      uint16
	height     uint16
	font       string
	foreground uint32
	background uint32
	text       string
	name       string
}

func (p Plaintext) run(ch chan []Block) {
	go func() {
		for {
			ch <- []Block{
				Block{
					rectangle: Rectangle{
						x:      p.x,
						y:      p.y,
						width:  p.width,
						height: p.height,
						color:  p.background,
					},
					text: Text{
						text:  p.text,
						font:  p.font,
						color: p.foreground,
					},
					name: p.name,
				},
			}
			time.Sleep(1000 * time.Second)
		}
	}()
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

func (t *Time) update() []Block {
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
	// public
	x          int16
	y          int16
	width      uint16
	height     uint16
	font       string
	foreground uint32
	background uint32

	// private
	inactive uint64
	active   uint64
}

func (c CPU) run(ch chan []Block) {
	go func() {
		for {
			ch <- c.update()
			time.Sleep(2 * time.Second)
		}
	}()
}

func (c *CPU) update() []Block {
	content, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return []Block{}
	}

	lines := strings.Split(string(content), "\n")
	counts := strings.Fields(lines[0])

	var inactive, active uint64 = 0, 0
	for i := 1; i < len(counts); i++ {
		count, err := strconv.ParseUint(counts[i], 10, 64)
		if err != nil {
			return []Block{}
		}
		if i != 4 {
			active += count
		} else {
			inactive = count
		}
	}
	// calculate CPU usage
	usage := int(100 * (active - c.active) /
		((active - c.active) + (inactive - c.inactive)))
	c.active = active
	c.inactive = inactive

	return []Block{
		Block{
			rectangle: Rectangle{
				x:      c.x,
				y:      c.y,
				width:  c.width,
				height: c.height,
				color:  c.background,
			},
			text: Text{
				text:  strconv.Itoa(usage),
				font:  c.font,
				color: c.foreground,
			},
			name: "cpu",
		},
	}
}

// Memory module
//
// Module outputs current memory usage. Memory usage is calculated by polling the
// contents of /proc/meminfo. Uses MemAvailable in /proc/meminfo, requiring Linux
// kernel 3.14 or higher.
type Memory struct {
	x          int16
	y          int16
	width      uint16
	height     uint16
	font       string
	foreground uint32
	background uint32
}

func (m Memory) run(ch chan []Block) {
	go func() {
		for {
			ch <- m.update()
			time.Sleep(2 * time.Second)
		}
	}()
}

func (m *Memory) update() []Block {
	content, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return []Block{}
	}

	lines := strings.Split(string(content), "\n")
	memoryTotal, err := strconv.ParseInt(strings.Fields(lines[0])[1], 10, 0)
	memoryAvailable, err := strconv.ParseInt(strings.Fields(lines[2])[1], 10, 0)
	if err != nil {
		return []Block{}
	}
	// calculate memory usage
	usage := int(100 * (memoryTotal - memoryAvailable) / memoryTotal)

	return []Block{
		Block{
			rectangle: Rectangle{
				x:      m.x,
				y:      m.y,
				width:  m.width,
				height: m.height,
				color:  m.background,
			},
			text: Text{
				text:  strconv.Itoa(usage),
				font:  m.font,
				color: m.foreground,
			},
			name: "memory",
		},
	}
}

// Battery module
//
// Module output current battery percentage. Battery percentage is calculated
// by polling /sys/class/power_supply/BAT1/capacity, requiring Linux kernel 3.19
// or higher.
type Battery struct {
	x          int16
	y          int16
	width      uint16
	height     uint16
	font       string
	foreground uint32
	background uint32
}

func (b Battery) run(ch chan []Block) {
	go func() {
		for {
			ch <- b.update()
			time.Sleep(15 * time.Second)
		}
	}()
}

func (b *Battery) update() []Block {
	content, err := ioutil.ReadFile("/sys/class/power_supply/BAT1/capacity")
	if err != nil {
		return []Block{}
	}

	percentage, err := strconv.ParseInt(strings.Split(string(content), "\n")[0], 10, 0)
	if err != nil {
		return []Block{}
	}

	return []Block{
		Block{
			rectangle: Rectangle{
				x:      b.x,
				y:      b.y,
				width:  b.width,
				height: b.height,
				color:  b.background,
			},
			text: Text{
				text:  strconv.Itoa(int(percentage)),
				font:  b.font,
				color: b.foreground,
			},
			name: "battery",
		},
	}
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

func (w *Workspace) update() []Block {
	workspaces, err := i3.GetWorkspaces()
	if err != nil {
		return []Block{}
	}

	blocks := make([]Block, len(workspaces))
	x := w.x
	for i, workspace := range workspaces {
		blocks[i] = Block{
			rectangle: Rectangle{
				x:      x,
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
		x += int16(w.width) + 5
	}
	return blocks
}
