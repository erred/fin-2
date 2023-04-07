package serve

import (
	"context"
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.seankhliao.com/svcrunner/v2/observability"
	"go.seankhliao.com/svcrunner/v2/tshttp"
	"go.seankhliao.com/webstyle"
)

type Server struct {
	o *observability.O

	svr *tshttp.Server

	render webstyle.Renderer

	dir string
}

func New(ctx context.Context, c *Cmd) *Server {
	svr := tshttp.New(ctx, c.tshttpConf)
	s := &Server{
		o:   svr.O,
		svr: svr,

		render: webstyle.NewRenderer(webstyle.TemplateCompact),
		dir:    c.dir,
	}

	svr.Mux.Handle("/eur", otelhttp.NewHandler(s.hView("eur"), "hView - eur"))
	svr.Mux.Handle("/gbp", otelhttp.NewHandler(s.hView("gbp"), "hView - gbp"))
	svr.Mux.Handle("/twd", otelhttp.NewHandler(s.hView("twd"), "hView - twd"))
	svr.Mux.Handle("/", otelhttp.NewHandler(http.HandlerFunc(s.hIndex), "hIndex"))
	svr.Mux.HandleFunc("/-/ready", func(rw http.ResponseWriter, r *http.Request) { rw.Write([]byte("ok")) })
	return s
}

func (s *Server) Run(ctx context.Context) error {
	return s.svr.Run(ctx)
}

func (s *Server) hIndex(rw http.ResponseWriter, r *http.Request) {
	ctx, span := s.o.T.Start(r.Context(), "hIndex")
	defer span.End()

	if r.URL.Path != "/" {
		http.Redirect(rw, r, "/", http.StatusFound)
		return
	}

	c := `
# fin

## money

### _fin_

- [GBP](/gbp)
- [EUR](/eur)
- [TWD](/twd)
`

	err := s.render.Render(rw, strings.NewReader(c), webstyle.Data{})
	if err != nil {
		s.o.HTTPErr(ctx, "render", err, rw, http.StatusInternalServerError)
	}
}
