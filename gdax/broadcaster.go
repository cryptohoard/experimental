package main

import "context"

const (
	BTCUSD = "BTC-USD"
	BCHUSD = "BCH-USD"
	ETHUSD = "ETH-USD"
	LTCUSD = "LTC-USD"
)

type Price struct {
	ProductID string
	Price     float64
}

type broadcast struct {
	c chan broadcast
	v Price
}

type Broadcaster struct {
	// private fields:
	productID string
	Listenc   chan chan (chan broadcast)
	Sendc     chan<- Price
}

type PriceReceiver struct {
	// private fields:
	C chan broadcast
}

// create a new broadcaster object.
func NewBroadcaster(ctx context.Context, productID string) *Broadcaster {
	listenc := make(chan (chan (chan broadcast)))
	sendc := make(chan Price)
	go func() {
		currc := make(chan broadcast, 1)
		for {
			select {
			case v := <-sendc:
				if v.ProductID != productID {
					continue
				}
				c := make(chan broadcast, 1)
				b := broadcast{c: c, v: v}
				currc <- b
				currc = c
			case r := <-listenc:
				r <- currc
			case <-ctx.Done():
				currc <- broadcast{}
				return
			}
		}
	}()
	return &Broadcaster{
		Listenc: listenc,
		Sendc:   sendc,
	}
}

// start listening to the broadcasts.
func (b *Broadcaster) Register() *PriceReceiver {
	c := make(chan chan broadcast, 0)
	b.Listenc <- c
	return &PriceReceiver{<-c}
}

// broadcast a value to all listeners.
func (b *Broadcaster) Write(v Price) { b.Sendc <- v }

// read a value that has been broadcast,
// waiting until one is available if necessary.
func (r *PriceReceiver) Read() Price {
	b := <-r.C
	v := b.v
	r.C <- b
	r.C = b.c
	return v
}
