package pricetracker

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/websocket"
	gdax "github.com/preichenberger/go-gdax"
)

type PriceTracker struct {
	logger         log.Logger
	ws             *websocket.Conn
	btcBroadcaster *Broadcaster
	ltcBroadcaster *Broadcaster
	initialized    bool
}

func NewPriceTracker(logger log.Logger) *PriceTracker {
	return &PriceTracker{logger: logger, initialized: false}
}

func (pt *PriceTracker) Register(productID string) *PriceReceiver {
	ctx := context.Background()

	pt.server()

	switch productID {
	case BTCUSD:
		if pt.btcBroadcaster == nil {
			pt.btcBroadcaster = NewBroadcaster(ctx, BTCUSD)
		}
		return pt.btcBroadcaster.Register()
	case LTCUSD:
		if pt.ltcBroadcaster == nil {
			pt.ltcBroadcaster = NewBroadcaster(ctx, LTCUSD)
		}
		return pt.ltcBroadcaster.Register()
	}

	return &PriceReceiver{}
}

func (pt *PriceTracker) server() {
	if pt.initialized {
		return
	}

	pt.initialized = true

	method := "pricetracker.server"

	//server := "ws-feed-public.sandbox.gdax.com"
	server := "ws-feed.gdax.com"
	u := url.URL{Scheme: "wss", Host: server}
	level.Debug(pt.logger).Log(
		"method",
		method,
		"msg",
		fmt.Sprintf("connecting to %s", u.String()),
	)
	ws, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		errMsg := fmt.Sprint("Failed to connect to %v, error: %v", u, err)
		level.Error(pt.logger).Log("method", method, "msg", errMsg)
		panic(errMsg)
	}
	level.Debug(pt.logger).Log(
		"method",
		method,
		"msg",
		fmt.Sprintf("resp: %v", resp),
	)
	pt.ws = ws

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

	if err := pt.ws.WriteJSON(subscribe); err != nil {
		errMsg := fmt.Sprintf("Failed to write JSON, error: %v", err)
		level.Error(pt.logger).Log("method", method, "msg", errMsg)
		panic(errMsg)
	}

	go pt.gdaxReader()
}

func (pt *PriceTracker) gdaxReader() {
	method := "pricetracker.gdaxReader"

	type priceinfo struct {
		productID string
		price     float64
		sequence  int64
	}

	c := make(chan priceinfo, 100)

	go func() {
		message := gdax.Message{}
		for {
			if err := pt.ws.ReadJSON(&message); err != nil {
				errMsg := fmt.Sprintf("Failed to read JSON, error: %v", err)
				level.Error(pt.logger).Log("method", method, "msg", errMsg)
				panic(errMsg)
			}
			level.Debug(pt.logger).Log(
				"method",
				method,
				"msg",
				fmt.Sprintf("message: %+v", message),
			)
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
			if pt.btcBroadcaster != nil {
				pt.btcBroadcaster.Write(btcPrice)
			}
			if pt.ltcBroadcaster != nil {
				pt.ltcBroadcaster.Write(ltcPrice)
			}
		}
	}
}
