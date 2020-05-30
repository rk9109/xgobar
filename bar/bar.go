package main

import (
	"errors"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

// Bar representation
type Bar struct {
	//
	conn   *xgb.Conn
	screen *xproto.ScreenInfo
	pixmap xproto.Pixmap
	window xproto.Window
	gc     xproto.Gcontext

	//
	font       map[string]xproto.Font
	fontWidth  map[string]int16
	fontHeight map[string]int16

	//
	rectangle Rectangle

	//
	modules []Module

	//
	ch chan []Block
}

// Return a bar initialized to configuration
// NewBar initializes connection to X server, and retrieves setup information.
func NewBar() (*Bar, error) {
	conn, err := xgb.NewConn()
	if err != nil {
		return &Bar{}, err
	}

	setup := xproto.Setup(conn)
	screen := setup.DefaultScreen(conn)

	rectangle := Rectangle{
		x:      x,
		y:      y,
		width:  width,
		height: height,
		color:  background,
	}

	return &Bar{
		conn:       conn,
		screen:     screen,
		font:       make(map[string]xproto.Font),
		fontWidth:  make(map[string]int16),
		fontHeight: make(map[string]int16),
		rectangle:  rectangle,
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

	err = b.getFonts()
	if err != nil {
		return err
	}

	err = xproto.CreateWindowChecked(
		b.conn,
		b.screen.RootDepth,
		b.window,
		b.screen.Root,
		b.rectangle.x, b.rectangle.y,
		b.rectangle.width, b.rectangle.height,
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

	// initialize bar background
	b.drawBlock(Block{
		rectangle: b.rectangle,
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
func (b *Bar) getFonts() error {
	for _, fontName := range fonts {
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

		reply, err := xproto.QueryFont(b.conn, xproto.Fontable(font)).Reply()
		if err != nil {
			return err
		}

		b.font[fontName] = font
		b.fontWidth[fontName] = reply.MaxBounds.CharacterWidth
		b.fontHeight[fontName] = reply.FontAscent - reply.FontDescent
	}

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
		b.rectangle.width,
		b.rectangle.height,
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
		block.rectangle.x, block.rectangle.y,
		block.rectangle.x, block.rectangle.y,
		block.rectangle.width, block.rectangle.height,
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
		[]uint32{block.rectangle.color},
	).Check()
	if err != nil {
		return err
	}

	rectangle := xproto.Rectangle{
		X:      block.rectangle.x,
		Y:      block.rectangle.y,
		Width:  block.rectangle.width,
		Height: block.rectangle.height,
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
	if block.text.text == "" {
		return nil
	}

	font, ok := b.font[block.text.font]
	fontWidth, ok := b.fontWidth[block.text.font]
	fontHeight, ok := b.fontHeight[block.text.font]
	if !ok {
        return errors.New("invalid font")
	}

	// calculate coordinates to center text inside rectangle
	fontX := block.rectangle.x + (int16(block.rectangle.width)-int16(len(block.text.text))*fontWidth)/2
	fontY := block.rectangle.y + int16(block.rectangle.height)/2 + fontHeight/2

	err := xproto.ChangeGCChecked(
		b.conn,
		b.gc,
		xproto.GcForeground|xproto.GcBackground|xproto.GcFont,
		[]uint32{block.text.color, block.rectangle.color, uint32(font)},
	).Check()
	if err != nil {
		return err
	}

	err = xproto.ImageText8Checked(
		b.conn,
		byte(len(block.text.text)),
		xproto.Drawable(b.pixmap),
		b.gc,
		fontX, fontY,
		block.text.text,
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

	return nil
}
