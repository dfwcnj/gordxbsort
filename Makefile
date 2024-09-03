.DEFAULT_GOAL := build

.PHONY:fmt vet build

fmt:
	go fmt gordxbinsort/*.go
	go fmt rdxbin/main.go

vet: fmt
	go vet gordxbinsort/*.go
	go vet rdxbin/main.go

build: vet
	go build -o rdxbin  rdxbin/main.go gordxbinsort/*.go

profile:
	go test -blockprofile block.prof -cpuprofile cpu.prof -memprofile mem.prof -mutexprofile mutex.prof -bench .

test: vet clobber
	go test

clean:
	/bin/rm -f rdxbin *.test
	/bin/rm -f *.prof *.out

clobber: clean
	/bin/rm -f *.pdf
	/bin/rm -rf /tmp/rdxsort*
