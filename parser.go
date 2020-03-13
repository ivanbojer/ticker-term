package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	term "github.com/buger/goterm"
	"github.com/fatih/color"
)

type Parser struct {
	last        map[string]float64
	titleColors map[string]color.Attribute
	lastColors  map[string]color.Attribute
	unch        map[string]int

	ids   []string
	pairs string

	countUnchanged   int
	countShouldPause int
	sleep            bool
	shouldClear      bool

	start time.Time
	count int

	t30      time.Time
	t30Count int64

	t30Avg time.Duration

	symbols []string
}

func newParser() *Parser {
	p := &Parser{
		last:        make(map[string]float64),
		titleColors: make(map[string]color.Attribute),
		lastColors:  make(map[string]color.Attribute),
		unch:        make(map[string]int),
		start:       time.Now(),
		t30:         time.Now(),
	}

	for _, sym := range defaultPairsSlice {
		if id, ok := Symbols[sym]; ok {
			sym = id
		}
		if sym == "0000" {
			continue
		}
		p.ids = append(p.ids, sym)
	}
	p.pairs = strings.Join(p.ids, ",")

	return p
}

func (p *Parser) getLast(sym string) float64 {
	pl, _ := p.last[sym]
	return pl
}

func (p *Parser) getLastColor(sym string) color.Attribute {
	if lc, ok := p.lastColors[sym]; ok {
		return lc
	}
	return color.Concealed
}

func (p *Parser) getLastTitleColor(sym string) color.Attribute {
	if lc, ok := p.titleColors[sym]; ok {
		return lc
	}
	return color.Concealed
}

func (p *Parser) header(titlePad, lastPad, changePctPad, changePad int) {
	var (
		hsym  = "Symbol"
		hlast = "Last"
		hpct  = "%"
		hchg  = "+/-"
		hma   = "MA"
		hti   = "Technical"
	)

	head := fmt.Sprintf("%s%s%s%s%s%s%s%s %s%s%s%s",
		hsym,
		padding(hsym, titlePad),
		hlast,
		padding(hlast, lastPad),
		hpct,
		padding(hpct, changePctPad),
		hchg,
		padding(hchg, changePad),
		hma,
		padding(hma, 18),
		hti,
		padding(hti, 12),
	)
	head = color.New(color.BgBlack).Set().Add(color.FgWhite).SprintFunc()(head)
	term.Println(head)
}

func (p *Parser) status() {
	now := time.Now()
	avg := time.Duration(now.Sub(p.start).Milliseconds()/int64(p.count)) * time.Millisecond

	p.t30Count++
	if p.count%30 == 0 {
		p.t30Avg = time.Duration(now.Sub(p.t30).Milliseconds()/p.t30Count) * time.Millisecond
		p.t30Count = 0
		p.t30 = time.Now()
	}

	if p.count == 20000 {
		p.start = now.Add(-time.Duration(avg * 1000))
		p.count = 1000
	}

	status := fmt.Sprintf("\nTarget: 180ms  Avg: %s %s Trailing 30: %s",
		avg,
		padding(avg.String(), 12),
		p.t30Avg,
	)
	status = status + padding(status, 90)
	status = color.New(color.BgBlack).Set().Add(color.FgWhite).SprintFunc()(status)
	term.Println(toColor(status, color.Concealed, color.Concealed))
}

func (p *Parser) getColors(sym string, lastf float64, imap InstrumentMap) (current, title, defaultc color.Attribute) {

	prev := p.getLast(sym)
	p.last[sym] = lastf

	switch {
	case prev == 0:
		current, title = color.Concealed, color.Concealed
		p.lastColors[sym], p.titleColors[sym] = current, title
	case lastf > prev:
		current, title = color.Reset, color.Reset
		p.lastColors[sym], p.titleColors[sym] = color.FgGreen, title
		p.zeroUnch(sym)
	case lastf < prev:
		current, title = color.Reset, color.Reset
		p.lastColors[sym], p.titleColors[sym] = color.FgRed, title
		p.zeroUnch(sym)
	default:
		current = p.getLastColor(sym)
		title = p.getLastTitleColor(sym)
		p.lastColors[sym], p.titleColors[sym] = current, title

		p.incrUnch(sym)
	}

	defaultc = imap.Color(sym)
	return
}

///////////////////////////////

func (p *Parser) incrUnch(sym string) {
	p.unch[sym] = p.unch[sym] + 1
	if p.unch[sym] > 1000 {
		p.lastColors[sym], p.titleColors[sym] = color.Concealed, color.Concealed
	}

	p.countUnchanged++
	if p.countUnchanged >= len(p.symbols) {
		p.countShouldPause++
	}
}

func (p *Parser) zeroUnch(sym string) {
	p.unch[sym] = 0
	p.countShouldPause = 0
}

func (p *Parser) pause() {

	p.countShouldPause = 0
	term.Println(" ")
	term.Println("The last 250 requests returned the same value for each symbol.")
	term.Println("Switching to sleep mode. Will check for updated results every 60 seconds...")
	p.sleep = true

	for k, _ := range p.lastColors {
		p.lastColors[k] = color.Concealed
		p.titleColors[k] = color.Concealed
	}
}

func (p *Parser) flush() {
	if p.countShouldPause > 250 {
		p.pause()
	}

	if p.shouldClear {
		p.shouldClear = false
		term.Clear()
	}

	term.Flush()
	if p.sleep {
		time.Sleep(time.Second * 60)
		term.Clear()
	}

	p.countUnchanged = 0
}

/////////////////////////////

func unwrap(res *http.Response) InstrumentMap {
	wrap := map[string]interface{}{}

	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fail(err)
	}

	err = json.Unmarshal(out, &wrap)
	if err != nil {
		fail(err)
	}

	delete(wrap, "time")
	out, err = json.Marshal(&wrap)
	if err != nil {
		fail(err)
	}

	imap, err := newInstrumentMap(out)
	if err != nil {
		fail(err)
	}

	return imap
}
