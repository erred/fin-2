package submit

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"net/http"
	"net/url"
	"os"

	"github.com/google/subcommands"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.seankhliao.com/fin/v3/findata"
	"go.seankhliao.com/svcrunner/v2/observability"
)

type Cmd struct {
	o    observability.Config
	addr *url.URL
	file []byte
}

func (c *Cmd) Name() string     { return `submit` }
func (c *Cmd) Synopsis() string { return `submit a local file` }
func (c *Cmd) Usage() string {
	return `view [options...]

Send a local file to a server

Flags:
`
}

func (c *Cmd) SetFlags(f *flag.FlagSet) {
	c.o.SetFlags(f)
	f.Func("file", "file to process", func(s string) error {
		var err error
		c.file, err = os.ReadFile(s)
		return err
	})
	f.Func("addr", "server address to post to", func(s string) error {
		var err error
		c.addr, err = url.Parse(s)
		return err
	})
}

func (c *Cmd) Execute(ctx context.Context, f *flag.FlagSet, args ...any) subcommands.ExitStatus {
	o := observability.New(c.o)

	if c.addr == nil {
		o.Err(ctx, "no address provided", errors.New("no server"))
		return subcommands.ExitFailure
	}

	out, err := findata.DecodeOne(c.file)
	if err != nil {
		o.Err(ctx, "decode", err)
		return subcommands.ExitFailure
	}

	c.addr.Path = "/" + out.Currency

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	res, err := client.Post(c.addr.String(), "application/cue", bytes.NewReader(c.file))
	if err != nil {
		o.Err(ctx, "post", err)
		return subcommands.ExitFailure
	} else if res.StatusCode >= 400 {
		o.Err(ctx, "unexpected response", errors.New(res.Status))
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
