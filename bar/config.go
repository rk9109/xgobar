package main

var (
	x      int16  = 0
	y      int16  = 0
	width  uint16 = 1920
	height uint16 = 25

	foreground     uint32 = 0xF8F8F8
	foregroundDark uint32 = 0xB8B8B8
	background     uint32 = 0x181818

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
		Battery{
			x:          (1920 - 150),
			y:          y,
			width:      25,
			height:     height,
			font:       fonts[0],
			foreground: foreground,
			background: background,
		},
		Plaintext{
			x:          (1920 - 190),
			y:          y,
			width:      40,
			height:     height,
			font:       fonts[0],
			foreground: foregroundDark,
			background: background,
			text:       "BAT:",
		},
		CPU{
			x:          (1920 - 230),
			y:          y,
			width:      25,
			height:     height,
			font:       fonts[0],
			foreground: foreground,
			background: background,
		},
		Plaintext{
			x:          (1920 - 270),
			y:          y,
			width:      40,
			height:     height,
			font:       fonts[0],
			foreground: foregroundDark,
			background: background,
			text:       "CPU:",
		},
		Memory{
			x:          (1920 - 310),
			y:          y,
			width:      25,
			height:     height,
			font:       fonts[0],
			foreground: foreground,
			background: background,
		},
		Plaintext{
			x:          (1920 - 350),
			y:          y,
			width:      40,
			height:     height,
			font:       fonts[0],
			foreground: foregroundDark,
			background: background,
			text:       "RAM:",
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
