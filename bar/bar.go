package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

// TODO document
type Bar struct {
	//
	conn *xgb.Conn

	//
	wid xproto.Window

	//
	font xproto.Font

	//
	name string

	//
	x, y int16

	//
	w, h uint16

	//
	fg uint32
	bg uint32
}

func New(conn *xgb.Conn) Bar {
	// TEMP
	fontid, _ := initFont(conn, font)

	//
	return Bar{
		conn: conn,
		font: fontid,
		name: name,
		x:    x,
		y:    y,
		w:    w,
		h:    h,
		fg:   fg,
		bg:   bg,
	}
}

func (b Bar) Map() error {
	var err error

	//
	setup := xproto.Setup(b.conn)
	screen := setup.DefaultScreen(b.conn)

	b.wid, err = xproto.NewWindowId(b.conn)
	if err != nil {
		return err
	}

    // TODO create pixmap here

	err = xproto.CreateWindowChecked(
		b.conn,
		screen.RootDepth,
		b.wid,
		screen.Root,
		b.x,
		b.y,
		b.w,
		b.h,
		0,
		xproto.WindowClassInputOutput,
		screen.RootVisual,
		xproto.CwBackPixel, // remove (?)
		[]uint32{b.bg},     // remove (?)
	).Check()
	if err != nil {
		return err
	}

	//
	err = b.UpdateEWMH()
	if err != nil {
		return err
	}

	//
	xproto.MapWindow(b.conn, b.wid)

	//
	pix, err := initPixmap(b.conn, screen)
	if err != nil {
		return err
	}

	//
	xproto.ChangeWindowAttributes(
		b.conn,
		b.wid,
		xproto.CwBackPixmap,
		[]uint32{uint32(pix)},
	)

	//
	err = drawText(b.conn, b.font, pix, b.wid, "TEST")
	if err != nil {
		return err
	}

	return nil
}

// TODO
//
func (b Bar) UpdateEWMH() error {
	var err error

	dataAtom, err := getAtom(b.conn, "_NET_WM_WINDOW_TYPE_DOCK")
	if err != nil {
		return err
	}

	propAtom, err := getAtom(b.conn, "_NET_WM_WINDOW_TYPE")
	if err != nil {
		return err
	}

	atom, err := getAtom(b.conn, "ATOM")
	if err != nil {
		return err
	}

	data := make([]byte, 4)
	xgb.Put32(data[0:], uint32(dataAtom))

	err = xproto.ChangePropertyChecked(
		b.conn,
		xproto.PropModeReplace,
		b.wid,
		propAtom,
		atom,
		32,
		uint32(len(data)/4),
		data,
	).Check()
	if err != nil {
		return err
	}

	// TODO name

	return nil
}
