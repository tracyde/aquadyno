package main

import (
	"github.com/tracyde/dummylcd"
)

type msg []string

func main() {
	display := dummylcd.NewDummyLCD(20, 4)
	text := msg{
		"Testing line one",
		"Testing line two",
		"Testing line three",
		"Testing line four",
	}
	display.Print(text)
}
