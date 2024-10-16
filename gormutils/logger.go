package gormutils

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/major1201/kubetrack/log"
	gormLogger "gorm.io/gorm/logger"
)

var logger = log.L.WithName("dao")

type logrAdapter struct {
	logr logr.Logger
}

func NewLogrAdapter(l logr.Logger) gormLogger.Writer {
	return &logrAdapter{logr: l}
}

func (l *logrAdapter) Printf(s string, i ...interface{}) {
	l.logr.Info(fmt.Sprintf(s, i...))
}
