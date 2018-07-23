.DEFAULT_GOAL := build

build: build.all

build.all: build.macosx build.linux

build.macosx:
	GOOS=darwin \
	GOARCH=amd64 \
	go build \
	-ldflags="-s -w" \
	-o bin/cwametrics \
	main.go

build.linux:
	GOOS=linux \
	GOARCH=amd64 \
	go build \
	-ldflags="-s -w" \
	-o bin/cwametric.l.txt \
	main.go

-include .private/mks/*.mk