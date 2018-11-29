#!/bin/bash
GOOS=windows GOARCH=386 go build  -o wsatraining_win.exe -ldflags "-s -w" && upx wsatraining_win.exe && mv wsatraining_win.exe app/.
