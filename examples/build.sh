#!/bin/bash

for SRC in *.go
do
	echo ${SRC}
	BASE=${SRC%.*}
	GOOS=linux GOARCH=amd64 go build -o bin/${BASE} ${SRC}
	GOOS=windows GOARCH=amd64 go build -o bin/${BASE}.exe ${SRC}
done
