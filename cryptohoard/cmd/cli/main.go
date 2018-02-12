package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cryptohoard/experimental/cryptohoard/internal/betservice"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/linbit/termui"
	"google.golang.org/grpc"
)

func main() {
	// Setup logging
	logger := getLogger(os.Stderr)

	// Setup connection for client
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:10051", opts...)
	if err != nil {
		errMsg := fmt.Sprintf("fail to dial: %v", err)
		level.Error(logger).Log("msg", errMsg)
		panic(errMsg)
	}
	defer conn.Close()

	client := betservice.NewBetServiceClient(conn)
	callOpts := grpc.FailFast(true)

	var placeReqs []*betservice.PlaceBetRequest
	placeReqs = append(
		placeReqs,
		&betservice.PlaceBetRequest{
			CustomerId: "abcd",
			Product:    betservice.Product_LTCUSD,
			Amount:     100,
		},
	)
	placeReqs = append(
		placeReqs,
		&betservice.PlaceBetRequest{
			CustomerId: "xyz",
			Product:    betservice.Product_BTCUSD,
			Amount:     1000,
		},
	)

	for _, placeReq := range placeReqs {
		resp, err := client.PlaceBet(
			context.Background(),
			placeReq,
			callOpts,
		)
		if err != nil {
			errMsg := fmt.Sprintf("expected no error but got %v\n", err)
			level.Error(logger).Log(
				"msg",
				errMsg,
			)
			panic(errMsg)
		}

		level.Debug(logger).Log("msg", fmt.Sprintf("resp: %v\n", resp))
	}

	err = termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	t := termui.NewTable()
	t.FgColor = termui.ColorWhite
	t.BgColor = termui.ColorDefault
	t.Y = 0
	t.X = 0
	t.Width = 162
	t.Height = 7

	for {
		stream, err := client.ListBets(
			context.Background(),
			&betservice.BetFilter{},
			callOpts,
		)
		if err != nil {
			errMsg := fmt.Sprintf("expected no error but got %v\n", err)
			level.Error(logger).Log(
				"msg",
				errMsg,
			)
			panic(errMsg)
		}

		var rows [][]string
		rows = append(
			rows,
			[]string{
				"ID",
				"Currency",
				"Initial Amount",
				"Current Amount",
				"Profit",
				"Price",
			},
		)

		loop := true
		for loop {
			bet, err := stream.Recv()
			if err == io.EOF {
				loop = false
				continue
			}

			if bet.State == betservice.State_PROCESSING {
				continue
			}

			// level.Debug(logger).Log(
			// 	"msg",
			// 	fmt.Sprintf(
			// 		"ID: %s, product: %d, crypto: %f, initial amount: %f, current amount: %f",
			// 		bet.BetId,
			// 		bet.Product,
			// 		bet.CryptoCurrency,
			// 		bet.InitialAmount,
			// 		bet.CurrentAmount,
			// 	),
			// )

			cur := fmt.Sprintf("%f", bet.CryptoCurrency)
			switch bet.Product {
			case betservice.Product_BTCUSD:
				cur = cur + " BTC"
			case betservice.Product_LTCUSD:
				cur = cur + " LTC"
			}

			rows = append(
				rows,
				[]string{
					bet.BetId,
					cur,
					fmt.Sprintf("%.2f", bet.InitialAmount),
					fmt.Sprintf("%.2f", bet.CurrentAmount),
					fmt.Sprintf("%.2f", bet.ProfitPercent),
					fmt.Sprintf("%.2f", bet.Price),
				},
			)
		}

		t.SetRows(rows)
		termui.Render(t)

		time.Sleep(5 * time.Second)
	}
}

func getLogger(w io.Writer) log.Logger {
	logger := log.NewLogfmtLogger(w)
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "name", os.Args[0])
	logger = log.With(logger, "pid", os.Getpid())
	logger = log.With(logger, "caller", log.DefaultCaller)
	return logger
}
