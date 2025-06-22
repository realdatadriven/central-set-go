package response

import "net/http"

type MetricsResponseWriter struct {
	StatusCode    int
	BytesCount    int
	headerWritten bool
	wrapped       http.ResponseWriter
}

func NewMetricsResponseWriter(w http.ResponseWriter) *MetricsResponseWriter {
	return &MetricsResponseWriter{
		StatusCode: http.StatusOK,
		wrapped:    w,
	}
}

func (mw *MetricsResponseWriter) Header() http.Header {
	return mw.wrapped.Header()
}

func (mw *MetricsResponseWriter) WriteHeader(statusCode int) {
	mw.wrapped.WriteHeader(statusCode)

	if !mw.headerWritten {
		mw.StatusCode = statusCode
		mw.headerWritten = true
	}
}

func (mw *MetricsResponseWriter) Write(b []byte) (int, error) {
	mw.headerWritten = true

	n, err := mw.wrapped.Write(b)
	mw.BytesCount += n
	return n, err
}

func (mw *MetricsResponseWriter) Unwrap() http.ResponseWriter {
	return mw.wrapped
}
