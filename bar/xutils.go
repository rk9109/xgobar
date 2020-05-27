package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

// Intern an atom
// Return unique identifier (atom) associated to property name
func getAtom(conn *xgb.Conn, name string) (xproto.Atom, error) {
	reply, err := xproto.InternAtom(conn, true, uint16(len(name)), name).Reply()
	if err != nil {
		return 0, err
	}
	return reply.Atom, nil
}

// Update ...
// Based on https://github.com/BurntSushi/xgbutil/
func updateProp(conn *xgb.Conn, window xproto.Window, format byte,
	propName string, typeName string, data []byte) error {

	propAtom, err := getAtom(conn, propName)
	if err != nil {
		return err
	}

	typeAtom, err := getAtom(conn, typeName)
	if err != nil {
		return err
	}

	return xproto.ChangePropertyChecked(
		conn,
		xproto.PropModeReplace,
		window,
		propAtom,
		typeAtom,
		format,
		uint32(len(data)/(int(format)/8)),
		data,
	).Check()
}
