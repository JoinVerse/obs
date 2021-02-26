package hlog

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func decodeIfBinary(out *bytes.Buffer) (string, error) {
	p := out.Bytes()
	if len(p) == 0 || p[0] < 0x7F {
		return out.String(), nil
	}
	return "", fmt.Errorf("unknown")
}

func TestRequestIDFromHeaderHandler(t *testing.T) {
	out := &bytes.Buffer{}
	reqID := "514bbe5bb5251c92bd07a9846f4a1ab6"
	r := &http.Request{
		Header: http.Header{
			"X-Request-Id": []string{reqID},
		},
	}
	h := RequestIDHeaderHandler("id", "X-Request-Id")(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, ok := hlog.IDFromRequest(r)
				if ok {
					t.Fatal("Not missing id in request")
				}
				if want, got := "514bbe5bb5251c92bd07a9846f4a1ab6", w.Header().Get("X-Request-Id"); got != want {
					t.Errorf("Invalid Request-Id header, got: %s, want: %s", got, want)
				}
				l := hlog.FromRequest(r)
				l.Log().Msg("")
				got, err := decodeIfBinary(out)
				if err != nil {
					t.Fatal("Can not transform to string")
				}
				if want := fmt.Sprintf(`{"id":"%s"}`+"\n", reqID); want != got {
					t.Errorf("Invalid log output, got: %s, want: %s", got, want)
				}
			},
		),
	)
	h = hlog.NewHandler(zerolog.New(out))(h)
	h.ServeHTTP(httptest.NewRecorder(), r)
}
