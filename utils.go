package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type currencies struct {
	Poloniex []string `json:"poloniex"`
}

const poloniexWsApiUrl = "wss://api2.poloniex.com"

type currencyPairsMap map[string]string

func mapCurrencyPairs(pairs []string) currencyPairsMap {
	res := make(currencyPairsMap, len(pairs))
	for _, p := range pairs {
		currencies := strings.Split(p, "_")
		if len(currencies) < 2 {
			continue
		}
		normalized := currencies[1] + "_" + currencies[0]
		res[normalized] = p
	}

	return res
}

func (p currencyPairsMap) normalized() []string {
	res := make([]string, 0, len(p))
	for k := range p {
		res = append(res, k)
	}
	return res
}

func exitWithErr(err error) {
	fmt.Printf("error occurred: %s\n", err.Error())
	os.Exit(1)
}

var ErrHeartBeatMessage = errors.New("no transactions info, just heartbeat message")
var ErrNoTransactionsInfo = errors.New("no transactions info")

type externalError struct {
	Error string `json:"error"`
}
