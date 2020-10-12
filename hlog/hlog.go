package hlog

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"net/http"
	"os"
	"time"
)

// Logger is a middleware that logs end of each request, along with
// some useful data about what was requested, what the response status was,
// and how long it took to return.
func Logger(h http.Handler) http.Handler {
	host, _ := os.Hostname()
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("host", host).Logger()
	handler := hlog.NewHandler(logger)
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
	requestIDHandler := hlog.RequestIDHandler("req_id", "X-Request-Id")
	return handler(accessHandler(remoteAddrHandler(userAgentHandler(refererHandler(requestIDHandler(h))))))
}
