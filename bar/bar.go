package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Bar struct {
	//
	conn   *xgb.Conn
	screen *xproto.ScreenInfo
	pixmap xproto.Pixmap
	window xproto.Window
	gc     xproto.Gcontext

	//
	font     map[string]xproto.Font
	fontInfo map[string]*xproto.QueryFontReply

	//
	name string

	//
	x, y int16

	//
	width  uint16
	height uint16

	//
	foreground uint32
	background uint32

	//
	modules []Module

	//
	ch chan []Block
}

// Return a bar initialized to configuration
// NewBar initializes connection to X server, and retrieves setup information.
func NewBar() (Bar, error) {
	conn, err := xgb.NewConn()
	if err != nil {
		return Bar{}, err
	}

	setup := xproto.Setup(conn)
	screen := setup.DefaultScreen(conn)

	return Bar{
		conn:       conn,
		screen:     screen,
		name:       name,
		font:       make(map[string]xproto.Font),
		fontInfo:   make(map[string]*xproto.QueryFontReply),
		x:          x,
		y:          y,
		width:      width,
		height:     height,
		foreground: foreground,
		background: background,
		modules:    modules,
		ch:         make(chan []Block),
	}, nil
}

// Render bar on screen
// Map initializes resources associated to the bar, including the pixmap, graphics
// context, and font.
func (b *Bar) Map() error {
	window, err := xproto.NewWindowId(b.conn)
	b.window = window
	if err != nil {
		return err
	}

	err = b.getPixmap()
	if err != nil {
		return err
	}

	err = b.getGcontext()
	if err != nil {
		return err
	}

	for _, font := range fonts {
		err = b.getFont(font)
		if err != nil {
			return err
		}
	}

	err = xproto.CreateWindowChecked(
		b.conn,
		b.screen.RootDepth,
		b.window,
		b.screen.Root,
		b.x, b.y,
		b.width, b.height,
		0,
		xproto.WindowClassInputOutput,
		b.screen.RootVisual,
		xproto.CwBackPixmap,
		[]uint32{uint32(b.pixmap)},
	).Check()
	if err != nil {
		return err
	}

	err = b.updateEWMH()
	if err != nil {
		return err
	}

	xproto.MapWindow(b.conn, b.window)

	// TODO cleanup (?)
	b.drawBlock(Block{
		X: 0, Y: 0,
		W: b.width, H: b.height,
		Foreground: b.foreground,
		Background: b.background,
		Text:       "",
		Font:       "",
	})

	return nil
}

// Run modules and update bar
// Modules send updated blocks on a single channel; bar is redrawn upon
// receiving blocks on the channel.
func (b *Bar) Draw() error {
	// Launch background goroutines for each module
	for _, module := range b.modules {
		module.run(b.ch)
	}

	// Listen on channel for updated blocks
	for {
		blocks := <-b.ch
		for _, block := range blocks {
			err := b.drawBlock(block)
			if err != nil {
				return err
			}
		}
	}
}

// Load font and font properties
func (b *Bar) getFont(fontName string) error {
	font, err := xproto.NewFontId(b.conn)
	if err != nil {
		return err
	}

	err = xproto.OpenFontChecked(
		b.conn,
		font,
		uint16(len(fontName)),
		fontName,
	).Check()
	if err != nil {
		return err
	}

	fontInfo, err := xproto.QueryFont(b.conn, xproto.Fontable(font)).Reply()
	if err != nil {
		return err
	}

	b.font[fontName] = font
	b.fontInfo[fontName] = fontInfo

	return nil
}

// Return an uninitialized pixmap
func (b *Bar) getPixmap() error {
	pixmap, err := xproto.NewPixmapId(b.conn)
	if err != nil {
		return err
	}

	err = xproto.CreatePixmapChecked(
		b.conn,
		b.screen.RootDepth,
		pixmap,
		xproto.Drawable(b.screen.Root),
		b.width,
		b.height,
	).Check()
	if err != nil {
		return err
	}

	b.pixmap = pixmap

	return nil
}

// Return an uninitialized graphics context
func (b *Bar) getGcontext() error {
	gc, err := xproto.NewGcontextId(b.conn)
	if err != nil {
		return err
	}

	err = xproto.CreateGCChecked(
		b.conn,
		gc,
		xproto.Drawable(b.pixmap),
		0, []uint32{},
	).Check()
	if err != nil {
		return err
	}

	b.gc = gc

	return nil
}

// Update block
// drawBlock updates the bar pixmap and copies the updated area to
// the bar.
func (b *Bar) drawBlock(block Block) error {
	err := b.drawRect(block)
	if err != nil {
		return err
	}

	err = b.drawText(block)
	if err != nil {
		return err
	}

	xproto.CopyArea(
		b.conn,
		xproto.Drawable(b.pixmap),
		xproto.Drawable(b.window),
		b.gc,
		block.X, block.Y, block.X, block.Y,
		block.W, block.H,
	)

	return nil
}

// Update block rectangle
// drawRect updates the bar pixmap, but changes are not rendered until
// xproto.CopyArea() is called.
func (b *Bar) drawRect(block Block) error {
	err := xproto.ChangeGCChecked(
		b.conn,
		b.gc,
		xproto.GcForeground,
		[]uint32{block.Background},
	).Check()
	if err != nil {
		return err
	}

	rectangle := xproto.Rectangle{
		X:      block.X,
		Y:      block.Y,
		Width:  block.W,
		Height: block.H,
	}

	err = xproto.PolyFillRectangleChecked(
		b.conn,
		xproto.Drawable(b.pixmap),
		b.gc,
		[]xproto.Rectangle{rectangle},
	).Check()
	if err != nil {
		return err
	}

	return nil
}

// Update block text
// Text is centered inside the block
// drawText updates the bar pixmap, but changes are not rendered until
// xproto.CopyArea() is called.
func (b *Bar) drawText(block Block) error {
	err := xproto.ChangeGCChecked(
		b.conn,
		b.gc,
		xproto.GcForeground|xproto.GcBackground|xproto.GcFont,
		[]uint32{block.Foreground, block.Background, uint32(b.font[block.Font])},
	).Check()
	if err != nil {
		return err
	}

	// calculate (x, y) coordinates to center text
	fontX := block.X
	fontY := block.Y + int16(block.H)

	err = xproto.ImageText8Checked(
		b.conn,
		byte(len(block.Text)),
		xproto.Drawable(b.pixmap),
		b.gc,
		fontX, fontY,
		block.Text,
	).Check()
	if err != nil {
		return err
	}

	return nil
}

// Update EWMH properties
// TODO
func (b *Bar) updateEWMH() error {
	//
	dataAtom, err := getAtom(b.conn, "_NET_WM_WINDOW_TYPE_DOCK")
	if err != nil {
		return err
	}
	data := make([]byte, 4)
	xgb.Put32(data, uint32(dataAtom))

	err = updateProp(
		b.conn,
		b.window,
		32,
		"_NET_WM_WINDOW_TYPE",
		"ATOM",
		data,
	)
	if err != nil {
		return err
	}

	//
	err = updateProp(
		b.conn,
		b.window,
		8,
		"_NET_WM_NAME",
		"UTF8_STRING",
		[]byte(b.name),
	)

	return nil
}
