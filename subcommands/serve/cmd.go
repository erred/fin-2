package serve

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"go.seankhliao.com/svcrunner/v2/tshttp"
)

type Cmd struct {
	tshttpConf tshttp.Config

	dir string
}

func (c *Cmd) Name() string     { return `serve` }
func (c *Cmd) Synopsis() string { return `start server` }
func (c *Cmd) Usage() string {
	return `serve [options...]

Starts a server managing listening records

Flags:
`
}

func (c *Cmd) SetFlags(f *flag.FlagSet) {
	c.tshttpConf.SetFlags(f)
	f.StringVar(&c.dir, "data.dir", "", "directory to store data")
}

func (c *Cmd) Execute(ctx context.Context, f *flag.FlagSet, args ...any) subcommands.ExitStatus {
	err := New(ctx, c).Run(ctx)
	if err != nil {
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
