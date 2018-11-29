#!/bin/bash
#go build  -o wsatraining_mac.app -ldflags "-s -w" && mv wsatraining_mac.app app/.
#go build  -o wsatraining_mac.app -ldflags "-s -w" && upx "-9" wsatraining_mac.app && mv wsatraining_mac.app app/.
go build  -o wsatraining_mac.app -ldflags "-s -w"
