.PHONY: all deps

all: deps

deps:
	go get github.com/veandco/go-sdl2/sdl{,_mixer,_image,_ttf}

build:
	go build -o bin/flappy