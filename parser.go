package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func printTransactionInfo(data []byte, currencyPairName string) error {
	var message []interface{}
	if err := json.Unmarshal(data, &message); err != nil {
		return err
	}
	if len(message) < 2 {
		return ErrHeartBeatMessage
	}
	exchangeUpdates, ok := message[2].([]interface{})
	if !ok {
		return ErrNoTransactionsInfo
	}
	for _, v := range exchangeUpdates {
		transactionInfo, ok := v.([]interface{})
		if !ok {
			return fmt.Errorf("unexpected format")
		}
		infoType := transactionInfo[0]
		if infoType != tradeInfoCode {
			continue
		}
		tradeInfo, err := newRecentTrade(transactionInfo, currencyPairName)
		if err != nil {
			return err
		}
		tradeInfoJson, err := json.Marshal(tradeInfo)
		if err != nil {
			return err
		}
		fmt.Println(string(tradeInfoJson))
	}
	return nil
}

func newRecentTrade(transactionInfo []interface{}, currencyPairName string) (recentTrade, error) {
	// ["t", "<trade id>", <1 for buy 0 for sell>, "<price>", "<size>", <timestamp>, "<epoch_ms>"]
	// ["t", "42706057", 1, "0.05567134", "0.00181421", 1522877119, "1522877119341"]
	price, err := strconv.ParseFloat(transactionInfo[3].(string), 64)
	if err != nil {
		return recentTrade{}, err
	}
	amount, err := strconv.ParseFloat(transactionInfo[4].(string), 64)
	if err != nil {
		return recentTrade{}, err
	}
	unixMilli, err := strconv.ParseInt(transactionInfo[6].(string), 10, 64)
	if err != nil {
		return recentTrade{}, err
	}
	tradeInfo := recentTrade{
		Id:        transactionInfo[1].(string),
		Pair:      currencyPairName,
		Price:     price,
		Amount:    amount,
		Side:      mustMapSide(transactionInfo[2].(float64)),
		Timestamp: time.UnixMilli(unixMilli),
	}

	return tradeInfo, nil
}

func mustMapSide(s float64) string {
	switch s {
	case 0:
		return "sell"
	case 1:
		return "buy"
	}
	panic(fmt.Sprintf("unexpected side value: %f", s))
}

type recentTrade struct {
	Id        string    `json:"id"`        // ID транзакции
	Pair      string    `json:"pair"`      // Торговая пара (из списка выше)
	Price     float64   `json:"price"`     // Цена транзакции
	Amount    float64   `json:"amount"`    // Объём транзакции
	Side      string    `json:"side"`      // Как биржа засчитала эту сделку (как buy или как sell)
	Timestamp time.Time `json:"timestamp"` // Время транзакции
}

const tradeInfoCode = "t"
