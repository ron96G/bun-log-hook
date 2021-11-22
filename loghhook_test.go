package hooks

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/uptrace/bun"
)

var buf = bytes.NewBuffer(nil)
var hook = NewLogHook(buf, time.RFC3339, false,
	``,
	`{"ts": "${time}", "operation": "${operation}", "status": "${status}", "latency": "${duration}}"`+"\n",
	`{"ts": "${time}", "operation": "${operation}", "status": "${status}", "error": "${error}", "latency": "${duration}}"`+"\n",
)

func BenchmarkLoghook(b *testing.B) {
	b.ReportAllocs()
	buf.Reset()

	ctx := context.Background()
	e := &bun.QueryEvent{
		Query:     "SELECT",
		StartTime: time.Now(),
		Err:       nil,
	}

	for i := 0; i < b.N; i++ {
		hook.BeforeQuery(ctx, e)
		hook.AfterQuery(ctx, e)
	}
}

func BenchmarkLoghookFailed(b *testing.B) {
	b.ReportAllocs()
	buf.Reset()

	ctx := context.Background()
	e := &bun.QueryEvent{
		Query:     "SELECT",
		StartTime: time.Now(),
		Err:       fmt.Errorf("some error"),
	}

	for i := 0; i < b.N; i++ {
		hook.BeforeQuery(ctx, e)
		hook.AfterQuery(ctx, e)
	}
}
