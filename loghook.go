package hooks

import (
	"bytes"
	"context"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/uptrace/bun"
	"github.com/valyala/fasttemplate"
)

var _ bun.QueryHook = (*LogHook)(nil)

var DefaultLogHook = NewLogHook(os.Stdout, time.RFC3339, false,
	``,
	`{"ts": "${time}", "operation": "${operation}", "status": "${status}", "latency": "${duration}}"`+"\n",
	`{"ts": "${time}", "operation": "${operation}", "status": "${status}", "error": "${error}", "latency": "${duration}}"`+"\n",
)

type LogHook struct {
	TimeFormat     string
	ErrorFormat    string
	BeforeFormat   string
	AfterFormat    string
	errorTemplate  *fasttemplate.Template
	beforeTemplate *fasttemplate.Template
	afterTemplate  *fasttemplate.Template
	Output         io.Writer
	OnlyFailed     bool
	pool           *sync.Pool
}

func NewLogHook(output io.Writer, timeFormat string, onlyFailed bool, beforeFormat, afterFormat, errorFormat string) *LogHook {
	hook := new(LogHook)
	hook.Output = output
	hook.OnlyFailed = onlyFailed
	hook.TimeFormat = timeFormat
	hook.errorTemplate = fasttemplate.New(errorFormat, "${", "}")

	if !onlyFailed {
		hook.beforeTemplate = fasttemplate.New(beforeFormat, "${", "}")
		hook.afterTemplate = fasttemplate.New(afterFormat, "${", "}")
	}

	hook.pool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 256))
		},
	}

	return hook
}

func (h *LogHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	if h.beforeTemplate != nil {
		buf := h.pool.Get().(*bytes.Buffer)
		buf.Reset()
		defer h.pool.Put(buf)

		h.parse(h.beforeTemplate, buf, event)
		h.Output.Write(buf.Bytes())
	}
	return ctx
}

func (h *LogHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	buf := h.pool.Get().(*bytes.Buffer)
	buf.Reset()
	defer h.pool.Put(buf)

	if h.errorTemplate != nil && event.Err != nil {
		h.parse(h.errorTemplate, buf, event)

	} else if h.afterTemplate != nil {
		h.parse(h.afterTemplate, buf, event)
	}

	h.Output.Write(buf.Bytes())
}

func (h *LogHook) parse(template *fasttemplate.Template, buf *bytes.Buffer, event *bun.QueryEvent) error {
	_, err := template.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
		switch tag {
		case "time":
			return buf.WriteString(time.Now().Format(h.TimeFormat))
		case "query":
			return buf.WriteString(event.Query)
		case "affected_rows":
			if n, err := event.Result.RowsAffected(); err == nil {
				return buf.WriteString(strconv.FormatInt(n, 10))
			}
		case "error":
			return buf.WriteString(event.Err.Error())
		case "status":
			if event.Err == nil {
				return buf.WriteString("successful")
			}
			return buf.WriteString("failed")
		case "duration":
			delta := time.Since(event.StartTime)
			return buf.WriteString(delta.String())
		case "duration_sec":
			delta := time.Since(event.StartTime)
			return buf.WriteString(strconv.FormatFloat(delta.Seconds(), 'f', 4, 64))
		case "operation":
			return buf.WriteString(event.Operation())
		default:
		}

		return 0, nil
	})
	return err
}
