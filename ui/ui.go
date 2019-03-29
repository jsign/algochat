package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/jsign/algochat/algochat"
	"github.com/pkg/errors"
)

const (
	msgViewID      = "msgView"
	msgViewTitle   = "Algorand Chat"
	inputViewID    = "inputView"
	inputViewTitle = "Chat here"
	logViewID      = "logView"
	logViewTitle   = "Log"
)

// StartAndLoop runs de UI
func StartAndLoop(in <-chan *algochat.ChatMessage, out chan<- string, logg <-chan string) error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return errors.Wrap(err, "error creating gui")
	}
	defer g.Close()
	g.Highlight = false
	g.Cursor = true

	g.SetManagerFunc(func(g *gocui.Gui) error {
		setOutputView(g)
		setInputView(g)
		setLogView(g)
		return nil
	})

	err = setLogView(g)
	if err != nil {
		return errors.Wrap(err, "error while creating log view")
	}

	err = setOutputView(g)
	if err != nil {
		return errors.Wrap(err, "error while creating output view")
	}

	iv, err := setInputView(g)
	if err != nil {
		return errors.Wrap(err, "error while creating output view")
	}
	err = iv.SetCursor(0, 0)
	if err != nil {
		return errors.Wrap(err, "failed to reset cursor")
	}
	_, err = g.SetCurrentView(inputViewID)
	if err != nil {
		return errors.Wrap(err, "can't set focus to input view")
	}

	err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})
	if err != nil {
		return errors.Wrap(err, "couldn't bind to ctrl+c action")
	}

	err = g.SetKeybinding(inputViewID, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, input *gocui.View) error {
		inputBuffer := input.Buffer()
		if len(inputBuffer) >= 2 {
			out <- strings.TrimRight(inputBuffer, "\n")
		}
		if err := input.SetCursor(0, 0); err != nil {
			return errors.Wrap(err, "failed to reset cursor")
		}
		input.Clear()
		input.Rewind()
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "coudn't bind to enter action")
	}

	go readInChan(g, in)
	go showLogg(g, logg)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	return nil
}

func setLogView(g *gocui.Gui) error {
	maxWidth, maxHeight := g.Size()
	lv, err := g.SetView(logViewID, maxWidth-40, 0, maxWidth-1, maxHeight-4)
	if err != nil && err != gocui.ErrUnknownView {
		return errors.Wrap(err, "failed to create log view")
	}
	lv.Autoscroll = true
	lv.Title = logViewTitle
	lv.Wrap = true

	return nil
}

func setOutputView(g *gocui.Gui) error {
	maxWidth, maxHeight := g.Size()
	mv, err := g.SetView(msgViewID, 0, 0, maxWidth-41, maxHeight-4)
	if err != nil && err != gocui.ErrUnknownView {
		return errors.Wrap(err, "failed to create messages view")
	}
	mv.Autoscroll = true
	mv.Title = msgViewTitle
	mv.Wrap = true

	return nil
}

func setInputView(g *gocui.Gui) (*gocui.View, error) {
	maxWidth, maxHeight := g.Size()
	input, err := g.SetView(inputViewID, 0, maxHeight-3, maxWidth-1, maxHeight-1)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, errors.Wrap(err, "failed to create input view")
	}
	input.Editable = true
	input.Title = inputViewTitle
	input.FgColor = gocui.ColorGreen

	return input, nil
}

func readInChan(g *gocui.Gui, in <-chan *algochat.ChatMessage) {
	mv, _ := g.View(msgViewID)
	for {
		select {
		case m := <-in:
			g.Update(func(g *gocui.Gui) error {
				_, _ = fmt.Fprintf(mv, "<%v-%v>: %v\n", m.Addr, m.Username, m.Message)
				return nil
			})
		}
	}
}

func showLogg(g *gocui.Gui, logg <-chan string) {
	lv, _ := g.View(logViewID)
	for {
		select {
		case l := <-logg:
			g.Update(func(g *gocui.Gui) error {
				_, _ = fmt.Fprintf(lv, "%v\n", l)
				return nil
			})
		}
	}
}
