#!/bin/bash

for SRC in *.go
do
	echo ${SRC}
	BASE=${SRC%.*}
	GOOS=linux GOARCH=amd64 go build -o bin/${BASE}.amd64 ${SRC}
	GOOS=linux GOARCH=arm go build -o bin/${BASE}.arm ${SRC}
	GOOS=windows GOARCH=amd64 go build -o bin/${BASE}.exe ${SRC}
done
