package cmd

import (
	"os"
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

func (question *Question) render() (err error) {
	tmpl, err := template.New("question").Parse(`
{{.Greet}} {{.Name}}です。
{{.UserInput.Overview}}という問題で困っています。

よろしくお願いいたします。	
`)
	if err != nil {
		return
	}
	err = tmpl.Execute(os.Stdout, question)
	return
}
