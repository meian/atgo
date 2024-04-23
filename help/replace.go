package help

import (
	"bytes"
	"text/template"

	"github.com/spf13/cobra"
)

type Replacer struct {
	Key     string
	Command *cobra.Command
}

// Replace はコマンドのヘルプメッセージ中のプレースホルダをコマンドパスで置換する。
// この処理は cmd 配下の init() で呼び出されるが、その時点では各コマンドの CommandPath() は初期化されていない。
// そのためヘルプを呼び出すタイミングで初期化させるために、SetHelpFunc() を使っている。
func Replace(cmd *cobra.Command, replacers []Replacer) {
	orgFunc := cmd.HelpFunc()
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		m := map[string]any{}
		for _, r := range replacers {
			m[r.Key] = r.Command.CommandPath()
		}
		t := template.Must(template.New("help").Parse(cmd.Short))
		var b bytes.Buffer
		t.Execute(&b, m)
		cmd.Short = b.String()
		t = template.Must(template.New("help").Parse(cmd.Long))
		b.Reset()
		t.Execute(&b, m)
		cmd.Long = b.String()
		orgFunc(cmd, args)
	})
}
