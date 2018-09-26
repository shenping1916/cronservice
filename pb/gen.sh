#/bin/sh
protoc --go_out=plugins=irpc:. *.proto
