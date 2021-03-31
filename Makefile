# meta info
NAME := question
VERSION := $(gobump show -r)

.DEFAULT_GOAL := build
export GO111MODULE=on

## Install dependencies
.PHONY: deps
deps:
	go get -v -d

# 開発に必要な依存をインストールする
## Setup
.PHONY: deps
devel-deps: deps
	GO111MODULE=off go get \
		github.com/motemen/gobump/cmd/gobump \
 		github.com/Songmu/make2help/cmd/make2help \
 		github.com/nsf/termbox-go \
 		github.com/sirupsen/logrus

## build binaries ex. make bin/question
bin/%: main.go deps
	go build -ldflags "$(LDFLAGS)" -o $@ $<

## build binary
.PHONY: build
build: bin/question

## Show help
.PHONY: help
help:
	@make2help $(MAKEFILE_LIST)
