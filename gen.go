package main

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cflags "-I/usr/include/aarch64-linux-gnu" -tags linux lb lb.c
