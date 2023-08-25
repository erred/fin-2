package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.seankhliao.com/svcrunner/v2/observability"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	conf := &Config{}

	// flags
	fset := flag.NewFlagSet("fin", flag.ExitOnError)
	conf.SetFlags(fset)
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "missing subcommand view|submit|serve")
		os.Exit(1)
	}
	fset.Parse(os.Args[2:])

	// observability
	o := observability.New(conf.o)

	// context
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var err error
	switch os.Args[1] {
	case "view":
		err = View(ctx, o, conf)
	case "submit":
		err = Submit(ctx, o, conf)
	case "serve":
		err = run(ctx, o, conf)
	default:
		fmt.Fprintln(os.Stderr, "unknown subcommand view|submit|serve:", os.Args[1])
		os.Exit(1)
	}
	if err != nil {
		o.Err(ctx, "exit", err)
	}
}

func run(ctx context.Context, o *observability.O, conf *Config) error {
	mux := http.NewServeMux()
	h2svr := &http2.Server{}
	svr := &http.Server{
		Handler:  h2c.NewHandler(mux, h2svr),
		ErrorLog: slog.NewLogLogger(o.H, slog.LevelWarn),
	}

	app := New(ctx, o, conf)
	app.Register(mux)

	addr := conf.address
	o.L.LogAttrs(ctx, slog.LevelInfo, "starting listener", slog.String("address", addr))
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return o.Err(ctx, "failed to listen", err, slog.String("address", addr))
	}

	go func() {
		<-ctx.Done()
		o.L.LogAttrs(ctx, slog.LevelInfo, "starting shutdown", slog.String("reason", context.Cause(ctx).Error()))
		err := lis.Close()
		if err != nil {
			o.Err(ctx, "failed to close listener", err, slog.String("address", addr))
		}
		err = svr.Shutdown(context.Background())
		if err != nil {
			o.Err(ctx, "failed to close server", err, slog.String("address", addr))
		}
	}()

	err = svr.Serve(lis)
	if err != nil {
		return o.Err(ctx, "failed to serve", err)
	}
	return nil
}
