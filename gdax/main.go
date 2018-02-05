package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	gdax "github.com/preichenberger/go-gdax"
)

const (
	BTCUSD = "BTC-USD"
	LTCUSD = "LTC-USD"
)

func main() {
	//server := "ws-feed-public.sandbox.gdax.com"
	server := "ws-feed.gdax.com"
	u := url.URL{Scheme: "wss", Host: server}
	log.Printf("connecting to %s", u.String())
	ws, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("Failed to connect to %v, error: %v", u, err)
	}
	log.Printf("resp: %v", resp)

	defer func() {
		if err := ws.Close(); err != nil {
			log.Fatalf("Failed to Close, error: %v", err)
		}
	}()

	subscribe := gdax.Message{
		Type: "subscribe",
		Channels: []gdax.MessageChannel{
			gdax.MessageChannel{
				Name: "ticker",
				ProductIds: []string{
					BTCUSD,
					LTCUSD,
				},
			},
		},
	}

	if err := ws.WriteJSON(subscribe); err != nil {
		log.Fatalf("Failed to write JSON, error: %v", err)
	}

	message := gdax.Message{}
	for i := 0; i < 10; i++ {
		if err := ws.ReadJSON(&message); err != nil {
			log.Fatalf("Failed to read JSON, error: %v", err)
		}
		log.Printf("message: %+v", message)
	}
}
