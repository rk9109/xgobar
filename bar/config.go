package main

import "time"

var (
	name = "TEST_BAR"

	x = 0
	y = 0
	w = 1920
	h = 25

	fg = 0xF8F8F8
	bg = 0x181818

	fontName = "-*-gohufont-medium-r-*-*-14-*-*-*-*-*-*-*"

	modules = []Module{
		Test{
			X:          25,
			Y:          0,
			W:          100,
			H:          uint16(h),
			Foreground: 0xF8F8F8,
			Background: 0x7CAFC2,
			Interval:   500 * time.Millisecond,
		},
		Test{
			X:          150,
			Y:          0,
			W:          100,
			H:          uint16(h),
			Foreground: 0xF8F8F8,
			Background: 0xAB4642,
			Interval:   1000 * time.Millisecond,
		},
	}
)
