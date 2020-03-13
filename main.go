package main

import (
	"fmt"
	"net/http"
	"os"

	term "github.com/buger/goterm"
	"github.com/pkg/errors"
)

func (p *Parser) Start(responses chan *http.Response, symbols ...string) (err error) {
	p.symbols = symbols

	for {
		p.count++
		term.MoveCursor(0, 0)

		imap := unwrap(<-responses)

		for symIdx, sym := range symbols {
			title := imap.Title(sym)

			if sym == "0000" {
				term.Print(title)
				p.countUnchanged++
				continue
			}

			last, lastf, err := imap.Last(sym)
			if err != nil {
				return errors.Wrap(err, "Parser.Start @ imap.Last")
			}

			change, changePct := imap.Change(sym)

			curColor, titleColor, defaultColor := p.getColors(sym, lastf, imap)
			titlePad, lastPad, changePctPad, changePad := 22, 10, 10, 15

			ln := fmt.Sprintf("%s%s%s%s%s%s%s%s %s",
				toColor(title, titleColor, titleColor),
				padding(title, titlePad),
				toColor(last, curColor, titleColor),
				padding(last, lastPad),
				toColor(changePct, defaultColor, titleColor),
				padding(changePct, changePctPad),
				toColor(change, defaultColor, titleColor),
				padding(change, changePad),
				imap.Technical(sym, titleColor),
			)

			if symIdx == 0 {
				p.header(titlePad, lastPad, changePctPad, changePad)
			}
			term.Println(ln)

		}

		p.status()
		p.flush()
	}
}

func main() {

	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	term.Clear()

	responses := make(chan *http.Response, 2)

	go fetch(responses, defaultPairsSlice...)

	p := newParser()
	err := p.Start(responses, defaultPairsSlice...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
