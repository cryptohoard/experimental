package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cryptohoard/experimental/cryptohoard/internal/betservice"
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

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "name", os.Args[0])
	logger = log.With(logger, "pid", os.Getpid())
	logger = log.With(logger, "mode", mode)
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = level.NewFilter(logger, level.AllowAll())

	level.Debug(logger).Log("msg", fmt.Sprintf("Starting in %s mode", mode))

	var g group.Group

	// Setup gRPC service
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("failed to listen: %v", err))
		os.Exit(-1)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// Initialize and register Bet service
	betsvc := betservice.NewBetSvc(logger)
	betservice.RegisterBetServiceServer(grpcServer, betsvc)

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
