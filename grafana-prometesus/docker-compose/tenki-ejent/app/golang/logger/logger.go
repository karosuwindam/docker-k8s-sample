package logger

import (
	"os"
	"path/filepath"

	"log/slog"
)

func Info(msg string, args ...interface{}) {
	ops := &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: replace,
	}
	tmp := slog.New(slog.NewJSONHandler(os.Stdout, ops))
	tmp.Info(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	ops := &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: replace,
	}
	tmp := slog.New(slog.NewJSONHandler(os.Stdout, ops))
	tmp.Debug(msg, args...)
}

func Error(msg string, args ...interface{}) {
	ops := &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: replace,
	}
	tmp := slog.New(slog.NewJSONHandler(os.Stdout, ops))
	tmp.Error(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	ops := &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: replace,
	}
	tmp := slog.New(slog.NewJSONHandler(os.Stdout, ops))
	tmp.Warn(msg, args...)
}

func replace(groups []string, a slog.Attr) slog.Attr {
	// Remove time.
	// if a.Key == slog.TimeKey && len(groups) == 0 {
	// 	return slog.Attr{}
	// }
	// Remove the directory from the source's filename.
	if a.Key == slog.SourceKey {
		source := a.Value.Any().(*slog.Source)
		source.File = filepath.Base(source.File)
	}
	return a
}
