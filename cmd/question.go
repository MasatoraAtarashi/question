package cmd

import (
	"bytes"
	"strings"
	"text/template"
)

// Question struct
type Question struct {
	Name      string
	Greet     string
	UserInput UserInput
}

// NewQuestion returns pointer of Question struct that made by options
func NewQuestion(userInput UserInput, name string) (*Question, error) {
	question := &Question{
		Name:      name,
		Greet:     "お疲れ様です。",
		UserInput: userInput,
	}
	return question, nil
}

func (question *Question) Execute() (result string, err error) {
	s := `{{.Greet}} {{.Name}}です。
{{.UserInput.RequireUserInput.Subject }}についてご質問させていただきたいです。

{{ if .UserInput.Reference -}}
現在、{{.UserInput.Reference}}を参考に{{.UserInput.RequireUserInput.Ideal}}を実現したいと思っております。	
{{ else -}}
現在、{{.UserInput.RequireUserInput.Ideal}}を実現したいと思っております。
{{ end -}}

私の行った手順は以下です。
{{.UserInput.RequireUserInput.Procedure}}

すると、以下のエラーが発生しました。
{{.UserInput.RequireUserInput.Problem}}

{{ if .UserInput.Source -}}
該当のソースコードはこちらです。
#source
{{ .UserInput.Source}}
#source
{{ end -}}

原因を確かめるため、以下を試してみましたが、問題の解決には至りませんでした。
{{ .UserInput.RequireUserInput.TriedAction }}

{{ if .UserInput.Env -}}
なお、私の環境は以下の通りです。
{{ .UserInput.Env }}
{{ end -}}

{{ if .UserInput.Log -}}
詳細なログは以下の通りです。
#source
{{ .UserInput.Log }}
#source
{{ end }}

何卒よろしくお願いいたします。	
`
	s = strings.Replace(s, "#source", "```", 4)
	tmpl, err := template.New("question").Parse(s)
	if err != nil {
		return
	}
	var doc bytes.Buffer
	err = tmpl.Execute(&doc, question)
	result = doc.String()
	return
}
