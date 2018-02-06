package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/oklog/oklog/pkg/group"
)

func main() {
	mode := "production"

	test := flag.Bool("test", false, "run in test mode.")

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

	logger.Log("msg", fmt.Sprintf("Starting in %s mode", mode))

	var g group.Group
	// Setup signal handler
	stop := make(chan struct{})
	g.Add(
		func() error { return signalHandler(stop) },
		func(error) { close(stop) })
	logger.Log("exit", g.Run())
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
