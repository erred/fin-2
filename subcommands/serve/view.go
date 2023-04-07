package serve

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"go.seankhliao.com/fin/v3/findata"
	"go.seankhliao.com/webstyle"
)

func (s *Server) hView(cur string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx, span := s.o.T.Start(r.Context(), "hView")
		defer span.End()

		var out findata.Currency
		switch r.Method {
		case http.MethodPost:
			b, err := io.ReadAll(r.Body)
			if err != nil {
				s.o.HTTPErr(ctx, "read request", err, rw, http.StatusBadRequest)
				return
			}
			out, err = findata.DecodeOne(b)
			if err != nil {
				s.o.HTTPErr(ctx, "decode data", err, rw, http.StatusBadRequest)
				return
			}

			err = os.WriteFile(filepath.Join(s.dir, cur+".cue"), b, 0o644)
			if err != nil {
				s.o.HTTPErr(ctx, "save data", err, rw, http.StatusInternalServerError)
				return
			}

		case http.MethodGet:
			b, err := os.ReadFile(filepath.Join(s.dir, cur+".cue"))
			if err != nil {
				s.o.HTTPErr(ctx, "read data file", err, rw, http.StatusInternalServerError)
				return
			}
			out, err = findata.DecodeOne(b)
			if err != nil {
				s.o.HTTPErr(ctx, "decode data", err, rw, http.StatusNotFound)
				return
			}

		default:
			s.o.HTTPErr(ctx, "GET or POST", errors.New("bad method"), rw, http.StatusMethodNotAllowed)
			return
		}

		var buf bytes.Buffer
		buf.WriteString("# ")
		buf.WriteString(cur)
		buf.WriteString("\n\n## currency view\n\n### _")
		buf.WriteString(cur)
		buf.WriteString("_\n\n")

		buf.WriteString("\n#### _holdings_\n\n")
		buf.Write(out.MarkdownTable(findata.ViewHoldings))
		buf.WriteString("\n#### _expenses_\n\n")
		buf.Write(out.MarkdownTable(findata.ViewExpenses))
		buf.WriteString("\n#### _income_\n\n")
		buf.Write(out.MarkdownTable(findata.ViewIncomes))

		err := s.render.Render(rw, &buf, webstyle.Data{})
		if err != nil {
			s.o.HTTPErr(ctx, "render", err, rw, http.StatusInternalServerError)
			return
		}
	})
}
