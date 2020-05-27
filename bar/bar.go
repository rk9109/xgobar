package main

import (
	"fmt" // TEMP

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type Bar struct {
	//
	conn   *xgb.Conn
	screen *xproto.ScreenInfo
	pixmap xproto.Pixmap
	window xproto.Window

	//
	font       xproto.Font
	fontWidth  int16
	fontHeight int16

	//
	gc xproto.Gcontext

	//
	name string

	//
	x, y int16

	//
	w, h uint16

	//
	fg uint32
	bg uint32

	//
	modules []Module

	//
	ch chan []Block
}

//
// TODO document
//
func NewBar() (Bar, error) {
	conn, err := xgb.NewConn()
	if err != nil {
		return Bar{}, err
	}

	setup := xproto.Setup(conn)
	screen := setup.DefaultScreen(conn)

	return Bar{
		conn:    conn,
		screen:  screen,
		name:    name,
		x:       int16(x), // avoid typecasting (?)
		y:       int16(y),
		w:       uint16(w),
		h:       uint16(h),
		fg:      uint32(fg),
		bg:      uint32(bg),
		modules: modules,
		ch:      make(chan []Block),
	}, nil
}

//
// TODO document
//
func (b *Bar) Map() error {
	window, err := xproto.NewWindowId(b.conn)
	b.window = window // TEMP
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

	err = b.getFont()
	if err != nil {
		return err
	}

	err = xproto.CreateWindowChecked(
		b.conn,
		b.screen.RootDepth,
		b.window,
		b.screen.Root,
		b.x, b.y,
		b.w, b.h,
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

	return nil
}

//
// TODO document
//
func (b *Bar) Draw() error {
	for _, module := range b.modules {
		module.run(b.ch)
	}

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

func (b *Bar) getFont() error {
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

	b.font = font
	b.fontWidth = reply.MaxBounds.CharacterWidth
	b.fontHeight = reply.FontAscent + reply.FontDescent
	fmt.Println("fontAscent: ", reply.FontAscent)
	fmt.Println("fontDescent: ", reply.FontDescent)
	fmt.Println("fontWidth: ", b.fontWidth)
	fmt.Println("fontHeight: ", b.fontHeight)

	return nil
}

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
		b.w,
		b.h,
	).Check()
	if err != nil {
		return err
	}

	b.pixmap = pixmap

	return nil
}

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

func (b *Bar) drawBlock(block Block) error {
	err := b.drawRect(block)
	if err != nil {
		return err
	}

	err = b.drawText(block)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bar) drawRect(block Block) error {
	rectangle := xproto.Rectangle{
		X:      block.X,
		Y:      block.Y,
		Width:  block.W,
		Height: block.H,
	}

	err := xproto.ChangeGCChecked(
		b.conn,
		b.gc,
		xproto.GcForeground,
		[]uint32{block.Background},
	).Check()
	if err != nil {
		return err
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

	// TODO do not copy entire pixmap (?)
	xproto.CopyArea(
		b.conn,
		xproto.Drawable(b.pixmap),
		xproto.Drawable(b.window),
		b.gc,
		0, 0, 0, 0,
		b.w, b.h,
	)

	return nil
}

func (b *Bar) drawText(block Block) error {
	// TODO
	fontX := block.X + (int16(block.W)-b.fontWidth*int16(len(block.Text)))/2
	fontY := block.Y + (int16(block.H)-b.fontHeight)/2 + b.fontHeight - 3

	err := xproto.ChangeGCChecked(
		b.conn,
		b.gc,
		xproto.GcForeground|xproto.GcBackground|xproto.GcFont,
		[]uint32{block.Foreground, block.Background, uint32(b.font)},
	).Check()
	if err != nil {
		return err
	}

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

	xproto.CopyArea(
		b.conn,
		xproto.Drawable(b.pixmap),
		xproto.Drawable(b.window),
		b.gc,
		0, 0, 0, 0,
		b.w, b.h,
	)

	return nil
}

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
