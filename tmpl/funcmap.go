package tmpl

import (
	htmlTemplate "html/template"
	"os"
	"strings"
	textTemplate "text/template"

	"github.com/major1201/kubetrack/utils"
	"github.com/major1201/kubetrack/utils/funcx"
	"github.com/major1201/kubetrack/utils/goutils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// FuncMap a common template tmpl implement
type FuncMap map[string]any

// TextFuncMap convert to text template funcMap
func (fm FuncMap) TextFuncMap() textTemplate.FuncMap {
	return textTemplate.FuncMap(fm)
}

// HTMLFuncMap convert to html template funcMap
func (fm FuncMap) HTMLFuncMap() htmlTemplate.FuncMap {
	return htmlTemplate.FuncMap(fm)
}

// GetFuncMap returns the common funcMap
func GetFuncMap() FuncMap {
	return fm
}

var fm = FuncMap{
	// numbers
	"int":   funcx.Partial1Of2FromTail(utils.IntDefault, 0),
	"intdv": utils.IntDefault,
	"inc":   inc,
	"add":   add,
	"sub":   sub,
	"mul":   mul,
	"div":   div,
	"mod":   mod,
	"rand":  random,

	// strings
	"title":      cases.Title(language.English).String,
	"replaceall": strings.ReplaceAll,
	"trim":       goutils.Trim,
	"trimleft":   goutils.TrimLeft,
	"trimright":  goutils.TrimRight,
	"upper":      strings.ToUpper,
	"lower":      strings.ToLower,
	"join":       goutils.PartialSwap(strings.Join),
	"split":      goutils.PartialSwap(strings.Split),
	"hasprefix":  goutils.PartialSwap(strings.HasPrefix),
	"hassuffix":  goutils.PartialSwap(strings.HasSuffix),
	"contains":   goutils.PartialSwap(strings.Contains),
	"indent":     indent,
	"uuid":       goutils.UUID,
	"filesize":   goutils.FileSize[int],
	"leftpad":    goutils.LeftPad,
	"rightpad":   goutils.RightPad,

	// bool
	"bool":   funcx.Partial1Of2FromTail(utils.BoolDefault, false),
	"booldv": utils.BoolDefault,

	// encoding
	"base64en":   encodeBase64,
	"base64de":   decodeBase64,
	"md5":        encodeMd5,
	"sha1":       encodeSha1,
	"sha224":     encodeSha224,
	"sha256":     encodeSha256,
	"sha512":     encodeSha512,
	"json":       toJSON,
	"prettyjson": toPrettyJSON,
	"yaml":       toYAML,

	// system
	"debug": debug,
	"env":   os.Getenv,
	"idx":   index,
}
