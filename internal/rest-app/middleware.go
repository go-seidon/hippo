package rest_app

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-seidon/hippo/internal/auth"
	"github.com/go-seidon/hippo/internal/datetime"
	"github.com/go-seidon/hippo/internal/logging"
	"github.com/go-seidon/hippo/internal/serialization"
	"github.com/go-seidon/hippo/internal/status"
)

type DefaultMiddlewareParam struct {
	CorrelationIdHeaderKey string
	CorrelationIdCtxKey    ContextKey
}

func NewDefaultMiddleware(p DefaultMiddlewareParam) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			correlationId := r.Header.Get(p.CorrelationIdHeaderKey)
			ctx := r.Context()
			ctx = context.WithValue(ctx, p.CorrelationIdCtxKey, correlationId)

			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	}
}

func NewBasicAuthMiddleware(a auth.BasicAuth, s serialization.Serializer) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authTokens := strings.Split(r.Header.Get("Authorization"), "Basic ")
			if len(authTokens) != 2 {
				Response(
					WithWriterSerializer(w, s),
					WithMessage("credential is not specified"),
					WithHttpCode(http.StatusUnauthorized),
					WithCode(status.ACTION_FORBIDDEN),
				)
				return
			}

			res, err := a.CheckCredential(context.Background(), auth.CheckCredentialParam{
				AuthToken: authTokens[1],
			})
			if err != nil {
				Response(
					WithWriterSerializer(w, s),
					WithHttpCode(http.StatusUnauthorized),
					WithCode(status.ACTION_FORBIDDEN),
					WithMessage("failed check credential"),
				)
				return
			}
			if !res.IsValid() {
				Response(
					WithWriterSerializer(w, s),
					WithMessage("credential is invalid"),
					WithHttpCode(http.StatusUnauthorized),
					WithCode(status.ACTION_FORBIDDEN),
				)
				return
			}

			h.ServeHTTP(w, r)
		})
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

type RequestLogMiddlewareParam struct {
	// required logger
	Logger logging.Logger

	// optional clock
	Clock datetime.Clock

	// key = uri
	// value = set `true` to ignore the uri being logged
	IgnoreURI map[string]bool

	// key = header key
	// value = log key
	Header map[string]string
}

func NewRequestLogMiddleware(p RequestLogMiddlewareParam) (func(h http.Handler) http.Handler, error) {
	if p.Logger == nil {
		return nil, fmt.Errorf("logger is not specified")
	}

	clock := p.Clock
	if p.Clock == nil {
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

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ignoreUri[r.RequestURI] {
				h.ServeHTTP(w, r)
				return
			}

			startTime := clock.Now()
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

			for key, val := range header {
				httpRequest[val] = r.Header.Get(key)
			}

			logger := p.Logger.WithFields(map[string]interface{}{
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
	}, nil
}
