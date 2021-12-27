// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = false
	g.Mouse = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}

var initialMouseX, initialMouseY, xOffset int
var msgMouseDown, movingMsg bool

func layout(g *gocui.Gui) error {
	if _, err := g.View("msg"); msgMouseDown && err == nil {
		moveMsg(g)
	}
	if v, err := g.SetView("but1", 2, 2, 22, 7, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "Button 1 - line 1")
		fmt.Fprintln(v, "Button 1 - line 2")
		fmt.Fprintln(v, "Button 1 - line 3")
		fmt.Fprintln(v, "Button 1 - line 4")
		if _, err := g.SetCurrentView("but1"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("but2", 24, 2, 44, 4, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "Button 2 - line 1")
	}
	updateHighlightedView(g)
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	for _, n := range []string{"but1", "but2"} {
		if err := g.SetKeybinding(n, gocui.MouseLeft, gocui.ModNone, showMsg); err != nil {
			return err
		}
	}
	if err := g.SetKeybinding("", gocui.MouseRelease, gocui.ModNone, mouseUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.MouseLeft, gocui.ModNone, msgDown); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func showMsg(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-10, maxY/2, maxX/2+10, maxY/2+2, 0); err == nil || errors.Is(err, gocui.ErrUnknownView) {
		v.Clear()
		v.SelBgColor = gocui.ColorCyan
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, l)
	}
	return nil
}

func updateHighlightedView(g *gocui.Gui) {
	mx, my := g.MousePosition()
	for _, view := range g.Views() {
		view.Highlight = false
	}
	if v, err := g.ViewByPosition(mx, my); err == nil {
		v.Highlight = true
	}
}

func moveMsg(g *gocui.Gui) {
	mx, my := g.MousePosition()
	if !movingMsg && (mx != initialMouseX || my != initialMouseY) {
		movingMsg = true
	}
	g.SetView("msg", mx-xOffset, my-1, mx-xOffset+20, my+1, 0)
}

func msgDown(g *gocui.Gui, v *gocui.View) error {
	initialMouseX, initialMouseY = g.MousePosition()
	xOffset, _ = v.Cursor()
	xOffset += 1
	msgMouseDown = true
	return nil
}

func mouseUp(g *gocui.Gui, v *gocui.View) error {
	if msgMouseDown {
		msgMouseDown = false
		if movingMsg {
			movingMsg = false
			return nil
		} else {
			g.DeleteView("msg")
		}
	}
	return nil
}
