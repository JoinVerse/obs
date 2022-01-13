package hlog

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

var logger *zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	host, _ := os.Hostname()
	l := zerolog.New(os.Stdout).With().Timestamp().Str("host", host).Logger()
	logger = &l
}

// A LoggerZ represents an active logging object that generates lines
// of JSON output to an io.Writer. Each logging operation makes a single
// call to the Writer's Write method. There is no guarantee on access
// serialization to the Writer. If your Writer is not thread safe,
// you may consider a sync wrapper.
type LoggerZ struct {
	zerolog.Logger
}

// New creates a root logger with os.Stdout writer.
func New() LoggerZ {
	return NewWithWriter(os.Stdout)
}

// NewWithWriter creates a root logger with given output writer
func NewWithWriter(w io.Writer) LoggerZ {
	host, _ := os.Hostname()
	l := zerolog.New(w).With().Timestamp().Str("host", host).Logger()
	return LoggerZ{l}
}

func (l *LoggerZ) Handler(h http.Handler) http.Handler {
	zLogger := l.Logger
	handler := hlog.NewHandler(zLogger)
	accessHandler := hlog.AccessHandler(
		func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		},
	)
	requestBodyHandler := RequestBodyHandler("requestBody")
	remoteAddrHandler := hlog.RemoteAddrHandler("ip")
	userAgentHandler := hlog.UserAgentHandler("userAgent")
	refererHandler := hlog.RefererHandler("referer")
	requestIDHandler := RequestIDHeaderHandler("requestId", "X-Request-Id")
	return handler(
		accessHandler(requestBodyHandler(remoteAddrHandler(userAgentHandler(refererHandler(requestIDHandler(h)))))),
	)
}

// RequestBodyHandler adds the requested Body as a field to the context's logger
// using fieldKey as field key.
func RequestBodyHandler(fieldKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Body != nil && r.Body != http.NoBody {
					log := zerolog.Ctx(r.Context())
					// Read the content
					var bodyBytes []byte
					if r.Body != nil {
						bodyBytes, _ = ioutil.ReadAll(r.Body)
					}
					// Restore the io.ReadCloser to its original state
					r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
					log.UpdateContext(
						func(c zerolog.Context) zerolog.Context {
							// Use the content
							if isJSON(bodyBytes) {
								return c.RawJSON(fieldKey, bodyBytes)
							}
							return c.Bytes(fieldKey, bodyBytes)
						},
					)
				}
				next.ServeHTTP(w, r)
			},
		)
	}
}

func isJSON(s []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(s, &js) == nil
}

// Deprecated: Use LoggerZ object instead.
// Logger is a middleware that logs end of each request, along with
// some useful data about what was requested, what the response status was,
// and how long it took to return.
func Logger(h http.Handler) http.Handler {
	l := &LoggerZ{*logger}
	return l.Handler(h)
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
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
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
					log.UpdateContext(
						func(c zerolog.Context) zerolog.Context {
							return c.Str(fieldKey, idStr)
						},
					)
				}
				if headerName != "" {
					w.Header().Set(headerName, idStr)
				}
				next.ServeHTTP(w, r)
			},
		)
	}
}
