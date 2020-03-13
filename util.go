package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	up    = emoji.Sprint(":up_arrow:")
	down  = emoji.Sprint(":down_arrow:")
	left  = emoji.Sprint(":left_arrow:")
	right = emoji.Sprint(":right_arrow:")
)

func fail(err ...interface{}) {
	fmt.Println(err...)
	os.Exit(1)
}

func termArea() (area int) {
	w, h, _ := terminal.GetSize(int(os.Stdout.Fd()))
	return w * h
}

func toColor(txt string, attr color.Attribute, override color.Attribute) string {
	c := color.New(color.Reset)

	if override == color.Concealed {
		attr = override
	}
	switch attr {
	case color.Concealed:
		c.Set().Add(color.FgWhite, color.Concealed, color.Faint)
	case color.Reset:
		c.Set().Add(color.FgHiWhite)
	case color.FgRed:
		c.Set().Add(color.FgHiRed)
	case color.FgGreen:
		c.Set().Add(color.FgHiGreen)
	default:
		c.Set().Add(attr)
	}

	return c.SprintFunc()(txt)
}

func padding(txt string, val int) string {
	l := utf8.RuneCountInString(txt)
	if val-l < 0 {
		err := errors.Errorf("Negative repeat count (%d) for [%s]", val-l, txt)
		fail(err)
	}
	return strings.Repeat(" ", val-l)
}
