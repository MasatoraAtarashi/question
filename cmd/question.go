package cmd

import "fmt"

// Question struct
type Question struct {
	Name        string
	Greet       string
	UserInput   UserInput
	Overview    string
	Probrem     string
	Procedure   string
	Expected    string
	Reference   string
	TriedAction string
	Log         string
	Env         string
}

// NewQuestion returns pointer of Question struct that made by options
func NewQuestion(userInput UserInput, name string) (*Question, error) {
	question := &Question{
		Name:      name,
		Greet:     "",
		UserInput: userInput,
	}
	return question, nil
}

func (question *Question) render() error {
	//content := "お疲れ様です。\n"
	//if question.name != "" {
	//	content += question.name + "です。\n"
	//}
	//content += question.situation + "\n\n"
	//content += "よろしくお願いいたします。"
	fmt.Println(question.UserInput.Overview)
	return nil
}
