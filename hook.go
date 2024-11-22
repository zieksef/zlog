package log

import (
	"context"
)

type CtxHook func(ctx context.Context) map[string]interface{}

func withCtx(contexts []context.Context) context.Context {
	if len(contexts) == 0 {
		return nil
	}

	return contexts[0]
}
