package hlog

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

var logger *zerolog.Logger

func init() {
	host, _ := os.Hostname()
	l := zerolog.New(os.Stdout).With().Timestamp().Str("host", host).Logger()
	logger = &l
}

// Logger is a middleware that logs end of each request, along with
// some useful data about what was requested, what the response status was,
// and how long it took to return.
func Logger(h http.Handler) http.Handler {
	handler := hlog.NewHandler(*logger)
	accessHandler := hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	})
	remoteAddrHandler := hlog.RemoteAddrHandler("ip")
	userAgentHandler := hlog.UserAgentHandler("user_agent")
	refererHandler := hlog.RefererHandler("referer")
	requestIDHandler := RequestIDHeaderHandler("req_id", "X-Request-Id")
	return handler(accessHandler(remoteAddrHandler(userAgentHandler(refererHandler(requestIDHandler(h))))))
}

type idKey struct{}

// RequestIDHeaderHandler adds given header from request's header as a field to
// the context's logger using fieldKey as field key. Returns a handler setting a unique
// id to the request which can be gathered using IDFromRequest(req). If the header does
// not exists this generated id is added as a field to the logger using the passed
// fieldKey as field name. The id is also added as a response header if the headerName
// is not empty.
//
// The generated id is a URL safe base64 encoded mongo object-id-like unique id.
// Mongo unique id generation algorithm has been selected as a trade-off between
// size and ease of use: UUID is less space efficient and snowflake requires machine
// configuration.
func RequestIDHeaderHandler(fieldKey, headerName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			idStr := r.Header.Get(headerName)
			if idStr == "" {
				id, ok := hlog.IDFromRequest(r)
				if !ok {
					id = xid.New()
					ctx = context.WithValue(ctx, idKey{}, id)
					r = r.WithContext(ctx)
				}
				idStr = id.String()
			}
			if fieldKey != "" {
				log := zerolog.Ctx(ctx)
				log.UpdateContext(func(c zerolog.Context) zerolog.Context {
					return c.Str(fieldKey, idStr)
				})
			}
			if headerName != "" {
				w.Header().Set(headerName, idStr)
			}
			next.ServeHTTP(w, r)
		})
	}
}