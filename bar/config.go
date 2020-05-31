package main

var (
	x      int16  = 0
	y      int16  = 0
	width  uint16 = 1920
	height uint16 = 25

	foreground uint32 = 0xF8F8F8
	background uint32 = 0x181818

	fonts = []string{
		"-*-gohufont-medium-r-*-*-14-*-*-*-*-*-*-*",
	}

	modules = []Module{
		Time{
			x:          (1920 - 110),
			y:          y,
			width:      100,
			height:     height,
			font:       fonts[0],
			foreground: foreground,
			background: background,
			format:     "Mon 2 03:04",
		},
		CPU{
			x:          (1920 - 160),
			y:          y,
			width:      50,
			height:     height,
			font:       fonts[0],
			foreground: foreground,
			background: background,
		},
		Workspace{
			x:                  10,
			y:                  y,
			width:              25,
			height:             height,
			font:               fonts[0],
			foreground:         foreground,
			backgroundActive:   0x7CAFC2,
			backgroundInactive: background,
		},
	}
)
