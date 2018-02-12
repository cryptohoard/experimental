package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	gdax "github.com/preichenberger/go-gdax"
)

var (
	ctx            = context.Background()
	btcBroadcaster *Broadcaster
	ltcBroadcaster *Broadcaster
)

func RegisterPriceReceiver(productID string) *PriceReceiver {
	switch productID {
	case BTCUSD:
		if btcBroadcaster == nil {
			btcBroadcaster = NewBroadcaster(ctx, BTCUSD)
		}
		return btcBroadcaster.Register()
	case LTCUSD:
		if ltcBroadcaster == nil {
			ltcBroadcaster = NewBroadcaster(ctx, LTCUSD)
		}
		return ltcBroadcaster.Register()
	}

	return &PriceReceiver{}
}

func server() *websocket.Conn {
	//server := "ws-feed-public.sandbox.gdax.com"
	server := "ws-feed.gdax.com"
	u := url.URL{Scheme: "wss", Host: server}
	log.Printf("connecting to %s", u.String())
	ws, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("Failed to connect to %v, error: %v", u, err)
	}
	log.Printf("resp: %v", resp)

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

	return ws
}

func gdaxReader(ws *websocket.Conn) {
	type priceinfo struct {
		productID string
		price     float64
		sequence  int64
	}

	c := make(chan priceinfo, 100)

	go func() {
		message := gdax.Message{}
		for {
			if err := ws.ReadJSON(&message); err != nil {
				log.Fatalf("Failed to read JSON, error: %v", err)
			}
			//log.Printf("message: %+v", message)
			c <- priceinfo{
				productID: message.ProductId,
				price:     message.Price,
				sequence:  message.Sequence,
			}
		}
	}()

	btcPrice := Price{ProductID: BTCUSD}
	btcSeq := int64(0)
	ltcPrice := Price{ProductID: LTCUSD}
	ltcSeq := int64(0)

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case p := <-c:
			switch p.productID {
			case BTCUSD:
				if p.sequence > btcSeq {
					btcPrice.Price = p.price
				}
			case LTCUSD:
				if p.sequence > ltcSeq {
					ltcPrice.Price = p.price
				}
			}
		case <-ticker.C:
			if btcBroadcaster != nil {
				btcBroadcaster.Write(btcPrice)
			}
			if ltcBroadcaster != nil {
				ltcBroadcaster.Write(ltcPrice)
			}
		}
	}
}

func main() {
	ws := server()

	defer func() {
		if err := ws.Close(); err != nil {
			log.Fatalf("Failed to Close, error: %v", err)
		}
	}()

	go gdaxReader(ws)

	btcPriceReceiver := RegisterPriceReceiver(BTCUSD)
	ltcPriceReceiver := RegisterPriceReceiver(LTCUSD)
	for {
		fmt.Printf("%+v\n", btcPriceReceiver.Read())
		fmt.Printf("%+v\n", ltcPriceReceiver.Read())
	}
}
