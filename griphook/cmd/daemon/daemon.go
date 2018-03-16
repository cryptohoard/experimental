package daemon

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	coinbase "github.com/cryptohoard/go-coinbase"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/julienschmidt/httprouter"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/api/iterator"
)

var (
	clientID     = os.GetEnv("GRIPHOOK_CLIENTID")
	clientSecret = os.GetEnv("GRIPHOOK_CLIENT_SECRET")
	redirectURL  = "https://104.196.230.1:443/callback"
	logger       = getLogger(os.Stderr, "test")
)

func Run() {
	mode := "test"

	// Setup logging
	//logger := getLogger(os.Stderr, mode)

	level.Debug(logger).Log("msg", fmt.Sprintf("Starting in %s mode", mode))

	var g group.Group

	// Setup signal handler
	stop := make(chan struct{})
	g.Add(
		func() error { return signalHandler(stop) },
		func(error) { close(stop) },
	)

	// Setup handlers
	r := httprouter.New()
	r.GET("/login", login)
	r.GET("/callback", callback)
	s := http.Server{
		Addr:    ":443",
		Handler: r,
	}
	g.Add(func() error {
		return s.ListenAndServeTLS("server.crt", "server.key")
	}, func(error) {
		s.Close()
	})

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

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := coinbase.NewOAuthClient(clientID, clientSecret, redirectURL)

	loginURL := c.CreateAuthorizeUrl(
		[]string{
			"wallet:accounts:read",
			"wallet:buys:read",
			"wallet:deposits:read",
			"wallet:orders:read",
			"wallet:sells:read",
			"wallet:transactions:read",
			"wallet:user:read",
		},
		"",
	)
	fmt.Fprintf(w, "%s\n", loginURL)
}

func callback(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	v := r.URL.Query()

	level.Debug(logger).Log("msg", fmt.Sprintf("GET params were: %v", v))

	code := v.Get("code")
	level.Debug(logger).Log("msg", fmt.Sprintf("code: %s", code))

	state := v.Get("state")
	level.Debug(logger).Log("msg", fmt.Sprintf("state: %s", state))

	if code != "" {
		c := coinbase.NewOAuthClient(clientID, clientSecret, redirectURL)

		tokens, err := c.Tokens(code)
		if err != nil {
			level.Error(logger).Log("msg", fmt.Sprintf("error: %v", err))
			http.Error(w, "couldn't get tokens", http.StatusInternalServerError)
			return
		}

		c.SetToken(tokens.AccessToken)

		it := c.Accounts()
		p := iterator.NewPager(it, 25, "")
		for {
			var accounts []coinbase.Account
			nextPageToken, err := p.NextPage(&accounts)
			if err != nil {
				level.Error(logger).Log("msg", fmt.Sprintf("error: %v", err))
				http.Error(w, "couldn't get accounts", http.StatusInternalServerError)
				return
			}
			for _, a := range accounts {
				level.Debug(logger).Log("msg", fmt.Sprintf("account: %+v", a))
			}
			if nextPageToken == "" {
				break
			}
		}
	}
}
