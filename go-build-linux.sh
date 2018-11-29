#!/bin/bash
#upx wsatraining_linux.elf &&
GOOS=linux GOARCH=amd64 go build -o wsatraining_linux.elf -ldflags "-s -w" && upx wsatraining_linux.elf && mv wsatraining_linux.elf app/.
