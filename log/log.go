package log

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-logr/logr"
	"github.com/major1201/kubetrack/third_party/glogr"
)

var (
	L   logr.Logger
	Std *log.Logger
)

func initGlog() {
	_ = flag.Set("v", "5")
	_ = flag.Set("logtostderr", "true")
	flag.Parse()
	L = glogr.NewWithOptions(glogr.Options{ErrorStack: glogr.DebugStack})
	Std = NewStd(L)
}

func initSlog() {
	opts := &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// reserve file name only
			if a.Key == slog.SourceKey {
				src := a.Value.Any().(*slog.Source)
				if src != nil {
					src.File = filepath.Base(src.File)
				}
			}
			return a
		},
	}
	// L = logr.FromSlogHandler(NewExtendedJSONHandler(os.Stdout, opts, true))
	L = logr.FromSlogHandler(NewExtendedTextHandler(os.Stdout, opts, true))
	Std = NewStd(L)
}

func init() {
	initSlog()
}

func NewStd(l logr.Logger) *log.Logger {
	return log.New(&writer{l: l}, "", log.Lshortfile)
}

type writer struct {
	l logr.Logger
}

func (w *writer) Write(p []byte) (n int, err error) {
	w.l.Info(string(p))
	return len(p), nil
}
