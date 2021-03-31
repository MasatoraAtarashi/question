question - 技術系質問作成ツール
=======

[![MIT License](http://img.shields.io/badge/license-Apache-blue.svg?style=flat)](LICENSE)

## Description
`question`は「ちゃんとした」質問を簡単に作成することができるツールです。

(「ちゃんとした」質問とは、目的・背景・ログ・試したこと等、質問された側が問題を解決するのに必要な情報を過不足なく盛り込んだ質問のことです。)

コマンドライン上で指示に従って情報を入力していくだけで、「ちゃんとした」質問文が完成されます。

「ちゃんとした」質問を作るのは意外と難しいので、このツールをぜひご利用ください。

## Installation
### Homebrew

	brew tap 
	brew install 

### go get
Install

    $ go get 

Update

    $ go get -u 

## Usage

    %  question init
- gif貼る

### push question to slack channel

    %  question push <slack channel>
- gif貼る

### show question log
    
    %  question log
- キャプチャ貼る

### to show other usage

    %  question --help

- gif貼る

## Options
```
  Options:
  -h,  --help                   print usage and exit
  -n,  --name <yourname>        specify your name
  -p,  --preamble <あいさつ>     specify a greeting
  -m,  --markdown               make the output in markdown format                 
  -t,  --tips                   print tips related to question and exit
  -p,  --procs <num>            split ratio to download file
```

## Author
[MasatoraAtarashi](https://github.com/MasatoraAtarashi)
