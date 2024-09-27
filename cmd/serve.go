package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/micamics/extracter/excel"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s excel.Service
	{
		s = excel.NewService()
		s = excel.Logging(logger)(s)
	}

	var h http.Handler
	{
		h = excel.CreateHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		if err := logger.Log("transport", "HTTP", "addr", *httpAddr); err != nil {
			slog.Error("logging", "error", err)
		}
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	if err := logger.Log("exit", <-errs); err != nil {
		slog.Error("logging", "error", err)
	}
}
