package main

import "time"

var (
	name string = ""

	x      int16  = 0
	y      int16  = 0
	width  uint16 = 1920
	height uint16 = 25

	foreground uint32 = 0xF8F8F8
	background uint32 = 0x181818

	fontName = "-*-gohufont-medium-r-*-*-14-*-*-*-*-*-*-*"

	modules = []Module{
		Test{
			X:          25,
			Y:          0,
			W:          100,
			H:          height,
			Foreground: 0xF8F8F8,
			Background: 0x7CAFC2,
			Interval:   500 * time.Millisecond,
		},
		Test{
			X:          150,
			Y:          0,
			W:          100,
			H:          height,
			Foreground: 0xF8F8F8,
			Background: 0xAB4642,
			Interval:   1000 * time.Millisecond,
		},
	}
)
