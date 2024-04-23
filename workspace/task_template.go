package workspace

import (
	"path/filepath"

	"github.com/meian/atgo/tmpl"
)

func TaskTemplate(name string) string {
	return filepath.Join(TemplateDir(), tmpl.TemplateName(name))
}
