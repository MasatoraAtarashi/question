package cmd

// Question struct
type Question struct {
	name 	string
	content string
}

// NewQuestion returns pointer of Question struct that made by options
func NewQuestion(name string) (*Question, error) {
	question := &Question{
		name: name,
		content: "To: Perlの技術ML\nFrom: とむら <tomura@example.org>\nSubject: CGIの文字化け\n\n先月からPerlでCGIを書きはじめた、とむら(仮名)と申します。\nCのプログラミング経験は3年ほどありますが、PerlやCGIははじめてです。\n\n市販されているXXXXスクリプトを修正して実行すると、\n文字化けを起こすという問題で困っています。\n原因または解決策をご存知の方はいらっしゃいませんか。\n\n私の行った手順は以下です。\n(1) XXXX.cgiをhttp://xxxx からダウンロードする。\n(2) 管理ファイル名$adminの値をadmin.datからadmin.txtに変更する。\n(3) 見出しを「一覧」から「ソフトウェアの表示」に変更する。\n(4) XXXX.cgiをサーバにFTPで転送する。\n(5) XXXX.cgiのパーミッションを 755(rwxr-xr-x)にする。\n(6) ブラウザからXXXXにアクセスする。\n\nすると、\n\n「ソフトウェアの表示」\n\nと表示されるべき部分が、\n\n「ャtトウェアの侮ヲ」\n\nと表示されてしまいます\n（最後の「ヲ」は半角文字です。\nそのままメールすると読めなくなると思い、ここだけ全角で書きなおしました。\nあとは表示されたものをそのままコピー＆ペーストしています）。\n\n手順の(3)を行わないと、正しく「一覧」と表示されます。\n\n関連すると思われるスクリプトは以下のようになっています。\n(これより前に100行ほどあるのですが、多すぎるのでカットしました)\n\nprint \"<head><title>$title</title>\";\nprint \"<body>\";\nprint \"<h1>ソフトウェアの表示</h1>\";\nprint \"</body>\";\nprint \"</html>\";\n\nローカルマシンのコマンドラインから、\n-cwで文法チェックしても異常は見つかりません(以下のようになります)。\n\nperl -cw XXXX.cgi\nXXXX.cgi syntax OK\n\nなお、私の環境は以下の通りです。\n(ローカルマシン)\n・Windows NT 4.0\n・Internet Explorer 5.0\n(サーバ)\n・自作CGIが使えるXXXXプロバイダのサーバ\n・サーバはApacheですが、バージョンは不明です。\n\n検索エンジンで「文字化け CGI」を検索しましたが、\n解決に役立つ情報は見つかりませんでした。\n過去ログも読みましたが、探し方がまずいのか、\n関連する情報を見つけることはできませんでした。\n\n----\nとむら(仮名) tomura@example.org",
	}
	return question, nil
}
