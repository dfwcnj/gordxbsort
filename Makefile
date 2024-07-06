.DEFAULT_GOAL := build

.PHONY:fmt vet build

fmt:
	go fmt *.go

vet: fmt
	go vet *.go

build: vet
	go build -o rdxbin  main.go klrsort2a.go merge.go input.go pqread.go

profile:
	go test -blockprofile block.prof -cpuprofile cpu.prof -memprofile mem.prof -mutexprofile mutex.prof -bench .

test:
	go test

clean:
	/bin/rm -f rdxbin *.test
	/bin/rm -f *.prof *.out

clobber: clean
	/bin/rm -f *.pdf
