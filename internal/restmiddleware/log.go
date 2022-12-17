package restmiddleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-seidon/provider/datetime"
	"github.com/go-seidon/provider/logging"
	"github.com/go-seidon/provider/logging/logrus"
)

type requestLog struct {
	logger    logging.Logger
	clock     datetime.Clock
	ignoreURI map[string]bool
	header    map[string]string
}

func (m *requestLog) Handle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.ignoreURI[r.RequestURI] {
			h.ServeHTTP(w, r)
			return
		}

		startTime := m.clock.Now()
		mw := &metricWriter{ResponseWriter: w}

		h.ServeHTTP(mw, r)

		timeElapsed := time.Since(startTime)
		duration := int64(timeElapsed) / int64(time.Millisecond)
		status, size := mw.CaptureMetric()
		httpRequest := map[string]interface{}{
			"requestMethod": r.Method,
			"requestUrl":    r.URL.String(),
			"requestSize":   size,
			"status":        status,
			"serverIp":      r.Host,
			"remoteAddr":    r.RemoteAddr,
			"protocol":      r.Proto,
			"receivedAt":    startTime.UTC().Format(time.RFC3339),
			"duration":      duration,
		}

		for key, val := range m.header {
			httpRequest[val] = r.Header.Get(key)
		}

		logger := m.logger.WithFields(map[string]interface{}{
			"httpRequest": httpRequest,
		})

		message := fmt.Sprintf("request: %s %s", r.Method, r.RequestURI)
		if status >= 100 && status <= 399 {
			logger.Info(message)
		} else if status >= 400 && status <= 499 {
			logger.Warn(message)
		} else {
			logger.Error(message)
		}
	})
}

type RequestLogParam struct {
	Logger logging.Logger
	Clock  datetime.Clock
	// key = uri
	// value = set `true` to ignore the uri being logged
	IgnoreURI map[string]bool
	// key = header key
	// value = log key
	Header map[string]string
}

func NewRequestLog(p RequestLogParam) *requestLog {
	logger := p.Logger
	if logger == nil {
		logger = logrus.NewLogger()
	}

	clock := p.Clock
	if clock == nil {
		clock = datetime.NewClock()
	}

	ignoreUri := map[string]bool{}
	if p.IgnoreURI != nil {
		ignoreUri = p.IgnoreURI
	}

	header := map[string]string{}
	if p.Header != nil {
		header = p.Header
	}
	header["User-Agent"] = "userAgent"
	header["Referer"] = "referer"
	header["X-Forwarded-For"] = "forwardedFor"

	return &requestLog{
		logger:    logger,
		clock:     clock,
		ignoreURI: ignoreUri,
		header:    header,
	}
}

type metricWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *metricWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *metricWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.size += n
	return n, err
}

func (w *metricWriter) CaptureMetric() (status int, size int) {
	return w.status, w.size
}
