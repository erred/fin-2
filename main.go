package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/subcommands"
	"go.seankhliao.com/fin/v3/subcommands/serve"
	"go.seankhliao.com/fin/v3/subcommands/submit"
	"go.seankhliao.com/fin/v3/subcommands/view"
)

func main() {
	name := "fin"
	fset := flag.NewFlagSet(name, flag.ExitOnError)
	cmdr := subcommands.NewCommander(fset, name)
	cmdr.Register(&serve.Cmd{}, "server")

	cmdr.Register(&view.Cmd{}, "client")
	cmdr.Register(&submit.Cmd{}, "client")

	cmdr.Register(cmdr.HelpCommand(), "other")

	fset.Parse(os.Args[1:])

	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	os.Exit(int(cmdr.Execute(ctx)))
}
