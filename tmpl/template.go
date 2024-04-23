package tmpl

import (
	"embed"
	"strings"
	tt "text/template"

	"github.com/pkg/errors"
)

//go:embed templates/*.tmpl templates/**/*.tmpl
var files embed.FS

var fm tt.FuncMap

func init() {
	fm = tt.FuncMap{
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

func templateWithFuncs(name string, text string) (*tt.Template, error) {
	return tt.New(name).Funcs(fm).Parse(text)
}
