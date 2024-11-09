package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		// Special handling for error attributes
		if a.Key == "error" {
			switch v := a.Value.Any().(type) {
			case error:
				fields[a.Key] = v.Error()
			case string:
				fields[a.Key] = v
			default:
				fields[a.Key] = fmt.Sprint(v)
			}
			return true
		}
		// Handle other attributes normally
		fields[a.Key] = a.Value.Any()
		return true
	})

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[15:04:05.000]")
	msg := color.CyanString(r.Message)

	var fileStr string
	if r.PC != 0 {
		frames := runtime.CallersFrames([]uintptr{r.PC})
		frame, _ := frames.Next()
		fileStr = color.GreenString(frame.File) + ":" + color.GreenString(fmt.Sprint(frame.Line))
	}

	h.l.Println(timeStr, level, msg, fileStr, color.WhiteString(string(b)))

	return nil
}

func NewPrettyHandler(
	out io.Writer,
	opts PrettyHandlerOptions,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}

func NewPrettyHandlerWithDefaults(level slog.Level) *PrettyHandler {
	return NewPrettyHandler(
		os.Stdout,
		PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			},
		},
	)
}

func NewJSONHandlerWithDefaults(level slog.Level) *slog.JSONHandler {
	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level, AddSource: true})
}
