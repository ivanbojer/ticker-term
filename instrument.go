package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type Instrument struct {
	Row struct {
		Last    string `json:"last"`
		Ma      string `json:"ma"`
		MaClass string `json:"ma_class"`
		Clock   string `json:"clock"`
	} `json:"row"`
	SummaryLast           string `json:"summaryLast"`
	SummaryName           string `json:"summaryName"`
	SummaryNameAlt        string `json:"summaryNameAlt"`
	SummaryChange         string `json:"summaryChange"`
	SummaryChangeClass    string `json:"summaryChangeClass"`
	TechnicalSummary      string `json:"technicalSummary"`
	TechnicalSummaryClass string `json:"technicalSummaryClass"`
	MaBuy                 int    `json:"maBuy"`
	MaSell                int    `json:"maSell"`
	TiBuy                 int    `json:"tiBuy"`
	TiSell                int    `json:"tiSell"`
}

type InstrumentMap map[string]*Instrument

func newInstrumentMap(in []byte) (imap InstrumentMap, err error) {
	imap = make(InstrumentMap)
	err = json.Unmarshal(in, &imap)
	if err != nil {
		return nil, err
	}
	return
}

func (im InstrumentMap) Color(sym string) color.Attribute {
	if id, ok := Symbols[sym]; ok {
		sym = id
	}
	if s, ok := im[sym]; ok {
		if s.SummaryChangeClass == "redFont" {
			return color.FgRed
		}
		return color.FgGreen
	}
	return color.Concealed
}

func (im InstrumentMap) Last(sym string) (str string, val float64, err error) {
	if id, ok := Symbols[sym]; ok {
		sym = id
	}
	if s, ok := im[sym]; ok {
		str = strings.ReplaceAll(s.Row.Last, ",", "")
		if len(str) == 0 {
			return "", 0, errors.New("InstrumentMap.Last @ get Symbol.Row.Last -- empty string")
		}

		val, err = strconv.ParseFloat(str, 32)
		if err != nil {
			return
		}
	}
	return
}

func (im InstrumentMap) Title(sym string) (val string) {
	id := sym
	var useSym bool

	if lbl, ok := Display[sym]; ok {
		return lbl
	}

	if pairID, ok := Symbols[sym]; ok {
		id = pairID
		useSym = true
	}
	if s, ok := im[id]; ok {
		if useSym {
			val = sym
			return
		}
		val = s.SummaryName
	}
	return
}

func (im InstrumentMap) Change(sym string) (change, changePct string) {
	if id, ok := Symbols[sym]; ok {
		sym = id
	}
	split := strings.Split(im[sym].SummaryChange, " ")
	return split[0], split[1]
}

func (im InstrumentMap) Technical(sym string, titleAttr color.Attribute) (val string) {
	if id, ok := Symbols[sym]; ok {
		sym = id
	}
	if s, ok := im[sym]; ok {
		ma := tline(s.Row.Ma, s.MaBuy, s.MaSell, titleAttr, 10)
		ti := tline(s.TechnicalSummary, s.TiBuy, s.TiSell, titleAttr, 4)

		return fmt.Sprintf("%s%s", ma, ti)
	}
	return
}

/////////////////////////////////////////

func tline(class string, nBuy, nSell int, titleAttr color.Attribute, pad int) string {
	offset := 6

	var (
		arrow string
		clr   color.Attribute
	)
	var lbl string
	switch class {
	case "Buy":
		arrow = up
		clr = color.FgGreen
	case "Strong Buy":
		arrow = up + up
		clr = color.FgGreen
		if nSell == 0 {
			arrow = arrow + up
		}
	case "Sell":
		arrow = down
		clr = color.FgRed
	case "Strong Sell":
		arrow = down + down
		clr = color.FgRed
		if nBuy == 0 {
			arrow = arrow + down
		}
	case "Neutral":
		arrow = left + right
		clr = color.FgWhite
	}

	lbl = toColor(fmt.Sprintf("%s%s", arrow, padding(arrow, offset)), clr, titleAttr)

	buyStr, sellStr := strconv.Itoa(nBuy), strconv.Itoa(nSell)
	sep := toColor("/", color.Concealed, titleAttr)
	bs := toColor(buyStr, color.FgGreen, titleAttr)
	ss := toColor(sellStr, color.FgRed, titleAttr)

	cntStr := bs + sep + ss

	return fmt.Sprintf("%s %s%s",
		lbl,
		cntStr,
		padding(buyStr+sellStr, pad),
	)
}
