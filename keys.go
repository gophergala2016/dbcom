package main

import "github.com/grafov/gocui"

func initKeys(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}
