sources = *.go

all: build

build: playback-haproxy

playback-haproxy: $(sources)
	time go build
