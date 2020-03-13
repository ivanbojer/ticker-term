package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/zpkg/ticker-term/internal/prand"
)

func randUserAgent() string {
	i := prand.IntN(len(userAgents) - 1)
	return userAgents[i]
}

const (
	uri = "https://ssltsw.forexprostools.com/api.php"
	qp  = "?action=refresher&timeframe=900"
)

func fetch(responses chan *http.Response, symbolIDs ...string) {

	qstr := qp + "&pairs=" + strings.Join(symbolIDs, ",")
	defaultReq, _ := http.NewRequest("GET", uri+qstr, nil)

	for {
		req := *defaultReq
		req.Header.Set("User-Agent", randUserAgent())
		res, err := http.DefaultClient.Do(&req)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		responses <- res
	}
}
