package noop

import "net/http"

// Exporter does nothing but implements the Exporter interface.
type Exporter struct{}

// New creates new noop.Exporter.
func New() *Exporter {
	return &Exporter{}
}

// CaptureError send error to nowhere.
func (*Exporter) CaptureError(err error, tags map[string]string, context map[string]interface{}) {}

// CaptureHttpError send error to nowhere.
func (*Exporter) CaptureHttpError(err error, r http.Request, tags map[string]string, context map[string]interface{}) {
}

// Close does nothing.
func (*Exporter) Close() {}
