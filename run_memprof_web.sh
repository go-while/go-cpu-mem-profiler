#!/bin/bash
PATH="$PATH:/usr/local/go/bin"
export GOPATH=$(pwd)
#export GO111MODULE=auto
#export GOEXPERIMENT=arenas
#go tool pprof cpu.pprof.webgrab http://127.0.0.1:1234/debug/pprof/profile
file=mem.pprof.out
test ! -z "$1" && file="$1"
ls -lh "$file"
go tool pprof -http=:17172 "$file"
