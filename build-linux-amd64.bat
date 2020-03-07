
@echo off

set GOARCH=amd64
set GOOS=linux
set GOHOSTARCH=amd64
set GOHOSTOS=windows

for %%i in ("%~dp0\.") do (
	set CurrentDir=%%~ni
)

set TARGET=bin\%CurrentDir%-%GOOS%-%GOARCH%
set TARGETDIR=bin\%GOOS%-%GOARCH%

mkdir %TARGETDIR%

echo building %TARGET% ..

go build -o %TARGET%

copy %TARGET% %TARGETDIR%

echo "OKay"

pause