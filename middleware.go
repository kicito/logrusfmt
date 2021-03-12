package logrusfmt

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	HeaderXForwardedFor = "X-Forwarded-For"
	HeaderXRealIP       = "X-Real-IP"
	HeaderXFinnoUserID  = "Finno-User-ID"
	HeaderXRequestID    = "X-Request-ID"
)

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func realIPFromHeader(h http.Header) string {
	// Fall back to legacy behavior
	if ip := h.Get(HeaderXForwardedFor); ip != "" {
		return strings.Split(ip, ", ")[0]
	}
	if ip := h.Get(HeaderXRealIP); ip != "" {
		return ip
	}
	return ""
}

func RequestHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip             string
			userID         string
			uniqueCookieID string
			method         string
			uri            string
			userAgent      string
			requestID      string
			ctx            context.Context = r.Context()
		)

		if ip = realIPFromHeader(r.Header); ip == "" {
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
		}
		userID = r.Header.Get(HeaderXFinnoUserID)
		// uniqueCookieID
		method = r.Method
		uri = r.RequestURI
		userAgent = r.UserAgent()
		requestID = r.Header.Get(HeaderXRequestID)

		if ip != "" {
			ctx = context.WithValue(ctx, CtxKeyIP, ip)
		}
		if userID != "" {
			ctx = context.WithValue(ctx, CtxKeyUserID, userID)
		}
		if uniqueCookieID != "" {
			ctx = context.WithValue(ctx, CtxKeyUniqueCookieID, uniqueCookieID)
		}
		if method != "" {
			ctx = context.WithValue(ctx, CtxKeyMethod, method)
		}
		if uri != "" {
			ctx = context.WithValue(ctx, CtxKeyURI, uri)
		}
		if userAgent != "" {
			ctx = context.WithValue(ctx, CtxKeyUserAgent, userAgent)
		}
		if requestID == "" {
			requestID, _ = randomHex(16)
		}
		ctx = context.WithValue(ctx, CtxKeyRequestID, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoggingHTTPMiddleware(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// use NewRecorder to record the response
			rec := httptest.NewRecorder()
			now := time.Now()
			next.ServeHTTP(rec, r)
			for k, v := range rec.Header() {
				w.Header()[k] = v
			}
			res := rec.Result()
			// write respose to actual response writer
			w.WriteHeader(res.StatusCode)
			_, _ = io.Copy(w, res.Body)
			respCtx := context.WithValue(r.Context(), CtxKeyStatus, res.StatusCode)
			respCtx = context.WithValue(respCtx, CtxKeyLatency, time.Since(now).Nanoseconds())
			msg := fmt.Sprintf("request log: %v %s %s", res.Status, r.Method, r.RequestURI)
			go func(ctx context.Context, message string) {
				logger.WithContext(respCtx).
					Info(msg)
			}(respCtx, msg)
		})
	}
}
