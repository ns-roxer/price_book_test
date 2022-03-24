package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
)

func read(conn net.Conn) ([]byte, error) {
	msg, opCode, err := wsutil.ReadServerData(conn)
	if err != nil {
		return nil, err
	}
	if opCode != ws.OpText {
		return nil, fmt.Errorf("unexpected data format")
	}
	var extErr externalError
	err = json.Unmarshal(msg, &extErr)
	if err == nil {
		return nil, fmt.Errorf(extErr.Error)
	}
	return msg, nil
}

func connect(ctx context.Context) (net.Conn, error) {
	conn, _, _, err := ws.DefaultDialer.Dial(ctx, poloniexWsApiUrl)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func subscribe(conn net.Conn, channel string) error {
	subscribeCommand := []byte(fmt.Sprintf(`{"command": "subscribe", "channel": %q}`, channel))
	return wsutil.WriteClientMessage(conn, ws.OpText, subscribeCommand)
}

func unsubscribe(conn net.Conn, channel string) error {
	subscribeCommand := []byte(fmt.Sprintf(`{"command": "unsubscribe", "channel": %q}`, channel))
	return wsutil.WriteClientMessage(conn, ws.OpText, subscribeCommand)
}
