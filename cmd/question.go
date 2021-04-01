package cmd

import "fmt"

// Question struct
type Question struct {
	name 	string
	content string
}

// NewQuestion returns pointer of Question struct that made by options
func NewQuestion(name string) (*Question, error) {
	question := &Question{
		name: name,
		content: "content of question",
	}
	return question, nil
}

func (question *Question) render() error {
	content := "お疲れ様です。\n"
	if question.name != "" {
		content += question.name + "です。\n"
	}
	content += question.content + "\n\n"
	content += "よろしくお願いいたします。"
	fmt.Println(content)
	return nil
}

