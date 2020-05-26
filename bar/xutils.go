package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func getAtom(conn *xgb.Conn, name string) (xproto.Atom, error) {
	//
	reply, err := xproto.InternAtom(conn, true, uint16(len(name)), name).Reply()
	if err != nil {
		return 0, err
	}
	return reply.Atom, nil
}

func initFont(conn *xgb.Conn, name string) (xproto.Font, error) {
	//
	font, err := xproto.NewFontId(conn)
	if err != nil {
		return 0, err
	}

	//
	err = xproto.OpenFontChecked(
		conn,
		font,
		uint16(len(name)),
		name,
	).Check()
	if err != nil {
		return 0, err
	}

	return font, nil
}

func initPixmap(conn *xgb.Conn, screen *xproto.ScreenInfo) (xproto.Pixmap, error) {
	//
	pix, err := xproto.NewPixmapId(conn)
	if err != nil {
		return 0, err
	}

	//
	err = xproto.CreatePixmapChecked(
		conn,
		screen.RootDepth,
		pix,
		xproto.Drawable(screen.Root),
		w, // update
		h, // update
	).Check()
	if err != nil {
		return 0, err
	}

	return pix, nil
}

func drawText(conn *xgb.Conn, font xproto.Font, pix xproto.Pixmap, wid xproto.Window, text string) error {
	//
	gcid, err := xproto.NewGcontextId(conn)
	if err != nil {
		return err
	}

	valueMask := uint32(xproto.GcForeground | xproto.GcBackground | xproto.GcFont)
	valueList := []uint32{
		0xFFFFFF,
		0x0000FF,
		uint32(font),
	}

	err = xproto.CreateGCChecked(
		conn,
		gcid,
		xproto.Drawable(pix), // pixmap
		valueMask,
		valueList,
	).Check()
	if err != nil {
		return err
	}

	err = xproto.ImageText8Checked(
		conn,
		byte(len(text)),
		xproto.Drawable(pix), // pixmap
		gcid,
		5,
		15,
		text,
	).Check()
	if err != nil {
		return err
	}

	xproto.CopyArea(
		conn,
		xproto.Drawable(pix),
		xproto.Drawable(wid),
		gcid,
		0, 0, 0, 0,
		w, h,
	)

	return nil
}
