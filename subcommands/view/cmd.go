package view

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
	"go.seankhliao.com/fin/v3/findata"
	"go.seankhliao.com/svcrunner/v2/observability"
)

type Cmd struct {
	o    observability.Config
	view findata.View
	file string
}

func (c *Cmd) Name() string     { return `view` }
func (c *Cmd) Synopsis() string { return `view a local file` }
func (c *Cmd) Usage() string {
	return `view [options...]

Render the contents of a local file

Flags:
`
}

func (c *Cmd) SetFlags(f *flag.FlagSet) {
	c.o.SetFlags(f)
	f.StringVar(&c.file, "file", "gbp.cue", "file to process")
	f.Func("view", "holdings|incomes|expenses", func(s string) error {
		switch s {
		case "holdings":
			c.view = findata.ViewHoldings
		case "incomes":
			c.view = findata.ViewIncomes
		case "expenses":
			c.view = findata.ViewExpenses
		default:
			return errors.New("unknown view")
		}
		return nil
	})
}

func (c *Cmd) Execute(ctx context.Context, f *flag.FlagSet, args ...any) subcommands.ExitStatus {
	o := observability.New(c.o)
	b, err := os.ReadFile(c.file)
	if err != nil {
		o.Err(ctx, "read file", err)
		return subcommands.ExitFailure
	}

	out, err := findata.DecodeOne(b)
	if err != nil {
		o.Err(ctx, "decode", err)
		return subcommands.ExitFailure
	}

	b = out.TabTable(c.view)
	fmt.Println(string(b))

	return subcommands.ExitSuccess
}
