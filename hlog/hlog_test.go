package hlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/stretchr/testify/assert"
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
				if want, got := reqID, w.Header().Get("X-Request-Id"); got != want {
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

type httpTestHandler struct{}

func (h *httpTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	fmt.Println("body: ", string(body))
	w.Write([]byte("ok"))
}

func TestBodyHandler(t *testing.T) {
	expectedLogBody := "{\"key1\":\"Value\",\"key2\":\"42\"}"
	r := &http.Request{
		URL:  &url.URL{Path: "/"},
		Body: ioutil.NopCloser(bytes.NewBuffer([]byte(expectedLogBody))),
	}

	httpTestHandler := &httpTestHandler{}
	recorder := httptest.NewRecorder()
	out := &bytes.Buffer{}
	logger := NewWithWriter(out)
	h := logger.Handler(httpTestHandler)
	h.ServeHTTP(recorder, r)

	if r.Body == nil || r.Body == http.NoBody {
		t.Fatal("Request body must be present")
	}
	assert.Equal(t, http.StatusOK, recorder.Code, "Response code must be 200")
	actualLog := struct {
		RequestBody string `json:"requestBody"`
	}{}
	err := json.Unmarshal([]byte(out.String()), &actualLog)
	assert.Nil(t, err)
	assert.Equal(t, expectedLogBody, actualLog.RequestBody, "Unexpected RequestBody Log String")
}
