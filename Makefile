.DEFAULT_GOAL := build

.PHONY:fmt vet build

fmt:
	go fmt main.go
	go fmt drsort2a.go

vet: fmt
	go vet main.go drsort2a.go

build: vet
	go build -o rdxbin  main.go drsort2a.go

profile:
	go test -cpuprofile cpu.prof -memprofile mem.prof -bench .

test:
	go test

clean:
	/bin/rm -f rdxsort godrdxsort.test
	/bin/rm -f cpu.prof mem.prof profile001.callgraph.out

