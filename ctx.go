package logrusfmt

import (
	"context"
)

type contextKey string

var (
	CtxKeyIP             = contextKey("ip")
	CtxKeyUserID         = contextKey("user_id")
	CtxKeyUniqueCookieID = contextKey("unique_cookie_id")
	CtxKeyMethod         = contextKey("method")
	CtxKeyStatus         = contextKey("status")
	CtxKeyURI            = contextKey("uri")
	CtxKeyUserAgent      = contextKey("user_agent")
	CtxKeyLatency        = contextKey("latency")
	CtxKeyRequestID      = contextKey("request_id")

	AllCtxKey = []contextKey{
		CtxKeyIP,
		CtxKeyUserID,
		CtxKeyUniqueCookieID,
		CtxKeyMethod,
		CtxKeyStatus,
		CtxKeyURI,
		CtxKeyUserAgent,
		CtxKeyLatency,
		CtxKeyRequestID,
	}
)

func setRequestContext(ctx context.Context, data map[string]interface{}) {
	for _, ctxKey := range AllCtxKey {
		value := ctx.Value(ctxKey)
		if value != nil {
			data[string(ctxKey)] = value
		}
	}
}
