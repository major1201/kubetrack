package tmpl

import (
	"strings"

	"github.com/major1201/kubetrack/utils/goutils"
)

func indent(indent int, s string) string {
	return goutils.Indent(s, strings.Repeat(" ", indent))
}
