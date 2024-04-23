package tmpl

import (
	"embed"
	"io"
	"os"
	"path/filepath"
	"strings"
	tt "text/template"

	"github.com/meian/atgo/text"
	"github.com/meian/atgo/url"
	"github.com/meian/atgo/util"
	"github.com/pkg/errors"
)

//go:embed templates/*.tmpl templates/**/*.tmpl
var files embed.FS

var fm tt.FuncMap

func init() {
	fm = tt.FuncMap{
		"padding":    text.PadRight,
		"date":       util.FormatTime,
		"nullint":    util.NullIntString,
		"mem":        util.FormatMemory,
		"duration":   util.FormatDuration,
		"stov":       util.StringToVar,
		"contesturl": url.ContestURL,
		"taskurl":    url.TaskURL,
		"shortFunc": func(s string) string {
			if len(s) == 0 {
				return ""
			}
			fs := strings.Split(s, "/")
			return fs[len(fs)-1]
		},
	}
}

func TemplateName(name string) string {
	return name + ".tmpl"
}

func textTemplate(name string) *tt.Template {
	data, err := files.ReadFile("templates/" + TemplateName(name))
	if err != nil {
		panic(errors.Wrapf(err, "failed to read template: name=%s", name))
	}
	return tt.Must(templateWithFuncs(name, string(data)))
}

func CmdTemplate(name string) *tt.Template {
	return textTemplate("cmd/" + name)
}

func LoggerTemplate() *tt.Template {
	return textTemplate("logger")
}

func TaskTemplateBinary(name string) io.Reader {
	data, err := files.Open("templates/task/" + TemplateName(name))
	if err != nil {
		panic(errors.Wrapf(err, "failed to read template: name=%s", name))
	}
	return data
}

func TaskTemplate(filename string) (*tt.Template, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	name := filepath.Base(strings.TrimSuffix(filename, filepath.Ext(filename)))
	return templateWithFuncs(name, string(data))
}

func templateWithFuncs(name string, text string) (*tt.Template, error) {
	return tt.New(name).Funcs(fm).Parse(text)
}
