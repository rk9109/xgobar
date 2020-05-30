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
		Clock{
			x:          (1920 - 100) / 2,
			y:          y,
			width:      100,
			height:     height,
			font:       fonts[0],
			foreground: foreground,
			background: background,
			format:     "3:04:05 PM",
		},
	}
)
