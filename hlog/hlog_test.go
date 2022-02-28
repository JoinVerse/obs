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
	if r.Body != nil && r.Body != http.NoBody {
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		fmt.Println("body: ", string(body))
	}
	_, _ = w.Write([]byte("ok"))
}

func TestBodyHandler(t *testing.T) {
	expectedLogBody := []byte(`{"key1":"Value","key2":"42"}`)
	var expectedJSONBody map[string]interface{}
	err := json.Unmarshal(expectedLogBody, &expectedJSONBody)
	assert.Nil(t, err)
	r := &http.Request{
		URL:  &url.URL{Path: "/"},
		Body: ioutil.NopCloser(bytes.NewBuffer(expectedLogBody)),
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
		RequestBody map[string]interface{} `json:"requestBody"`
	}{}
	err = json.Unmarshal(out.Bytes(), &actualLog)
	assert.Nil(t, err)
	assert.Equal(t, expectedJSONBody, actualLog.RequestBody, "Unexpected RequestBody Log")
}

// TestHTTPRequestLogFormat check the format https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#HttpRequest
func TestHTTPRequestLogFormat(t *testing.T) {
	expectedHTTPRequestLog :=
		[]byte(`{"latency":"0.000000s", "protocol":"HTTP/1.1", "referer":"https://example.com/", 
				"remoteIp":"188.26.219.97", "requestMethod":"GET", "requestUrl":"https://example.com/", 
				"responseSize":"2", "status":200, "userAgent":"obs"}`)
	var expectedHTTPRequestJSON map[string]interface{}
	err := json.Unmarshal(expectedHTTPRequestLog, &expectedHTTPRequestJSON)
	assert.Nil(t, err)
	r := &http.Request{
		Method:     http.MethodGet,
		RemoteAddr: "10.132.0.241",
		URL:        &url.URL{Scheme: "https", Host: "example.com", Path: "/"},
		Proto:      "HTTP/1.1",
		Header: map[string][]string{
			"User-Agent":      {"obs"},
			"Referer":         {"https://example.com/"},
			"X-Request-Id":    {"randomRequest123"},
			"X-Forwarded-For": {"188.26.219.97, 10.132.0.241"},
		},
	}

	httpTestHandler := &httpTestHandler{}
	recorder := httptest.NewRecorder()
	out := &bytes.Buffer{}
	logger := NewWithWriter(out)
	h := logger.Handler(httpTestHandler)
	h.ServeHTTP(recorder, r)

	assert.Equal(t, http.StatusOK, recorder.Code, "Response code must be 200")
	actualLog := struct {
		HTTPRequestLog map[string]interface{} `json:"httpRequest"`
	}{}
	err = json.Unmarshal(out.Bytes(), &actualLog)
	assert.Nil(t, err)

	actualLog.HTTPRequestLog["latency"] = "0.000000s" // Just ignore the latency value
	assert.Equal(t, expectedHTTPRequestJSON, actualLog.HTTPRequestLog, "Unexpected HttpRequest Log Format")
}
