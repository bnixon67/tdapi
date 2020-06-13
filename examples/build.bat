@echo off

FOR %%F IN (*.go) DO (

	SET GOOS=linux
	SET GOARCH=amd64
	echo go build -o bin\%%~nF.amd64 %%F
	go build -o bin\%%~nF.amd64 %%F

	SET GOOS=linux
	SET GOARCH=arm
	echo go build -o bin\%%~nF.arm %%F
	go build -o bin\%%~nF.arm %%F

	SET GOOS=windows
	SET GOARCH=amd64
	echo go build -o bin\%%~nF.exe %%F
	go build -o bin\%%~nF.exe %%F
)
