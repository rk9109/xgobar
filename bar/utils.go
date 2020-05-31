package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

// Intern an atom
// Return unique identifier (atom) associated to property name.
func getAtom(conn *xgb.Conn, name string) (xproto.Atom, error) {
	reply, err := xproto.InternAtom(conn, true, uint16(len(name)), name).Reply()
	if err != nil {
		return 0, err
	}
	return reply.Atom, nil
}

// Update properties
// Converts 32-bit properties to correct format.
// Based on https://github.com/BurntSushi/xgbutil/xprop/xprop.go
func updateProp32(conn *xgb.Conn, window xproto.Window, mode byte,
	propName string, typeName string, properties ...uint) error {

	data := make([]byte, len(properties)*4)
	for i, prop := range properties {
		xgb.Put32(data[(i*4):], uint32(prop))
	}

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
		mode,
		window,
		propAtom, typeAtom,
		32,
		uint32(len(data)/4),
		data,
	).Check()
}
