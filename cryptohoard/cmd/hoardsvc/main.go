package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cryptohoard/experimental/cryptohoard/internal/betservice"
	"github.com/cryptohoard/experimental/cryptohoard/internal/pricetracker"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
)

func main() {
	mode := "production"

	port := flag.Int("port", 10051, "The server port")
	test := flag.Bool("test", false, "run in test mode")

	flag.Parse()

	if *test {
		mode = "test"
	}

	// Setup logging
	logger := getLogger(os.Stderr, mode)

	level.Debug(logger).Log("msg", fmt.Sprintf("Starting in %s mode", mode))

	var g group.Group

	// Setup gRPC service
	lis, grpcServer, err := setupGRPCServer(*port)
	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("failed to listen: %v", err))
		os.Exit(-1)
	}

	// Register services to grpcServer
	registerServices(grpcServer, logger)

	// Start gRPC server
	g.Add(
		func() error { return grpcServer.Serve(lis) },
		func(error) { grpcServer.GracefulStop() },
	)

	// Setup signal handler
	stop := make(chan struct{})
	g.Add(
		func() error { return signalHandler(stop) },
		func(error) { close(stop) },
	)
	level.Debug(logger).Log("exit", g.Run())
}

func getLogger(w io.Writer, mode string) log.Logger {
	logger := log.NewLogfmtLogger(w)
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "name", os.Args[0])
	logger = log.With(logger, "pid", os.Getpid())
	logger = log.With(logger, "mode", mode)
	logger = log.With(logger, "caller", log.DefaultCaller)
	return logger
}

func setupGRPCServer(port int) (net.Listener, *grpc.Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return nil, nil, err
	}

	var opts []grpc.ServerOption
	return lis, grpc.NewServer(opts...), nil
}

func registerServices(grpcServer *grpc.Server, logger log.Logger) {
	// Create PriceTracker
	pt := pricetracker.NewPriceTracker(logger)

	// Initialize and register Bet service
	betsvc := betservice.NewBetSvc(logger, pt)
	betservice.RegisterBetServiceServer(grpcServer, betsvc)

	//go placebets(logger)
}

func placebets(logger log.Logger) {
	time.Sleep(5 * time.Second)

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

			level.Debug(logger).Log(
				"msg",
				fmt.Sprintf(
					"ID: %s, product: %d, crypto: %f, initial amount: %f, current amount: %f",
					bet.BetId,
					bet.Product,
					bet.CryptoCurrency,
					bet.InitialAmount,
					bet.CurrentAmount,
				),
			)
		}

		time.Sleep(5 * time.Second)
	}
}

func signalHandler(stop chan struct{}) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-c:
		return fmt.Errorf("received signal %s", sig)
	case <-stop:
		return nil
	}
}
