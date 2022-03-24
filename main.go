package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()
	input := flag.Arg(0)
	var inputStruct currencies
	if err := json.Unmarshal([]byte(input), &inputStruct); err != nil {
		exitWithErr(err)
	}

	currencies := mapCurrencyPairs(inputStruct.Poloniex)
	if len(currencies) == 0 {
		exitWithErr(fmt.Errorf("no correct currency pairs provided"))
	}
	eg, errGroupCtx := errgroup.WithContext(context.Background())
	ctx, stop := signal.NotifyContext(errGroupCtx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	for _, currencyChannel := range currencies.normalized() {
		currencyChannel := currencyChannel
		eg.Go(func() error {
			conn, err := connect(ctx)
			if err != nil {
				return err
			}
			defer conn.Close()
			if err := subscribe(conn, currencyChannel); err != nil {
				return err
			}
			for {
				select {
				case <-ctx.Done():
					return unsubscribe(conn, currencyChannel)
				default:
					data, err := read(conn)
					if err != nil {
						return err
					}
					if err := printTransactionInfo(data, currencies[currencyChannel]); err != nil {
						if errors.Is(err, ErrNoTransactionsInfo) || errors.Is(err, ErrHeartBeatMessage) {
							continue
						}
						return err
					}
				}
			}
		})
	}
	if err := eg.Wait(); err != nil {
		exitWithErr(err)
	}
	fmt.Println("bye!")
}
