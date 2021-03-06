package cmd

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

// Question struct
type Question struct {
	Greet     string
	UserInput UserInput
	Result    string
	MetaInfo  MetaInfo
}

type MetaInfo struct {
	Hash     string
	UserName string
	InitDate string
	Subject  string
}

// NewQuestion returns pointer of Question struct that made by options
func NewQuestion(userInput UserInput, name string) (*Question, error) {
	question := &Question{
		Greet:     "お疲れ様です。",
		UserInput: userInput,
		MetaInfo: MetaInfo{
			Hash:     "",
			UserName: name,
			InitDate: time.Now().String(),
			Subject:  userInput.RequireUserInput.Subject,
		},
	}
	return question, nil
}

func (question *Question) Execute() (err error) {
	s := `{{.Greet}} {{.MetaInfo.UserName}}です。
{{.UserInput.RequireUserInput.Subject}}についてご質問させていただきたいです。

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
	question.Result = doc.String()
	return
}

// 生成された質問を保存する
func (question *Question) Save() (err error) {
	// .questionフォルダを作成する
	home_dir := os.Getenv("HOME")
	qpath := home_dir + "/.question"
	if !fileExists(qpath) {
		err = os.Mkdir(qpath, 0777)
		if err != nil {
			return
		}
	}

	// 質問文保存フォルダを作成する
	objpath := qpath + "/objects"
	if !fileExists(objpath) {
		err = os.Mkdir(objpath, 0777)
		if err != nil {
			return
		}
	}

	// 質問文のメタ情報を作成する
	err = initMetaInfo(question)
	if err != nil {
		return
	}

	// 質問文を保存する
	data := []byte(question.Result)
	err = ioutil.WriteFile(objpath+"/"+question.MetaInfo.Hash, data, 0777)
	if err != nil {
		return
	}

	//logsフォルダを作成する
	logpath := qpath + "/logs"
	if !fileExists(logpath) {
		err = os.Mkdir(logpath, 0777)
		if err != nil {
			return
		}
	}

	// ログ・ファイルにメタ情報を追加する
	log_data, err := init_log_data_from_meta(question.MetaInfo)
	if err != nil {
		return
	}
	f, err := os.OpenFile(logpath+"/HEAD", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	defer f.Close()
	if _, err := f.WriteString(log_data); err != nil {
		log.Println(err)
	}
	return
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// 質問文のメタ情報を作成する
func initMetaInfo(question *Question) (err error) {
	// ハッシュ
	subject := question.UserInput.RequireUserInput.Subject
	name := question.UserInput.RequireUserInput.Subject
	initDate := question.MetaInfo.InitDate
	hash, err := initHash(subject + name + initDate)
	question.MetaInfo.Hash = hash
	return
}

// 文字列からハッシュ値を計算する
func initHash(str string) (hash string, err error) {
	p := []byte(str)
	hash = fmt.Sprintf("%x", md5.Sum(p))
	return
}

// ログ・ファイルに入力するデータを作る
func init_log_data_from_meta(meta MetaInfo) (log_data string, err error) {
	hash := meta.Hash
	userName := meta.UserName
	initDate := meta.InitDate
	subject := meta.Subject
	log_data = hash + " " + userName + " " + initDate + " " + subject + "\n"
	return
}
