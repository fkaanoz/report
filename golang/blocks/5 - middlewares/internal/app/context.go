package app

import "context"

type contextKey string

var Key contextKey = "context_key"

type Values struct {
	traceID    string
	statusCode int
}

func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(Key).(Values)
	if !ok {
		return ""
	}

	return v.traceID
}
