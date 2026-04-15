package pkg

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

// @Title        logger.go
// @Description
// @Create       2026-04-15 13:29
// @Update       2026-04-15 13:29

func JsonLogger() *slog.Logger {
	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.RFC3339))
				}
			}
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				shortFile := source.File
				for i := len(source.File) - 1; i > 0; i-- {
					if source.File[i] == '/' {
						shortFile = source.File[i+1:]
						break
					}
				}
				return slog.String("source", fmt.Sprintf("%s:%d", shortFile, source.Line))
			}
			return a
		},
	}))
	return jsonLogger
}
