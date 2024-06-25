.DEFAULT_GOAL := build

.PHONY:fmt vet build

fmt:
	go fmt main.go
	go fmt klrsort2a.go

vet: fmt
	go vet main.go klrsort2a.go input.go merge.go
	go vet randomdata.go klrsort2a.go klrsort2a_test.go

build: vet
	go build -o rdxbin  main.go klrsort2a.go merge.go input.go

profile:
	go test -cpuprofile cpu.prof -memprofile mem.prof -bench .

test:
	go test

clean:
	/bin/rm -f rdxbin *.test
	/bin/rm -f *.prof *.out

