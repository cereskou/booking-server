@echo off

set GOARCH=amd64
set GOOS=windows
set RELEASE=0
set GENERATE=0
set LINUX=0
set MODULE=booking.exe

if /I "%1"=="release" (
    set RELEASE=1
    echo build release.
)

if /I "%1"=="linux" (
    echo build linux.
    set LINUX=1
    set MODULE=booking
    
    if /I "%2"=="release" (
        set RELEASE=1
        echo build release.
    )
)

if not exist build.json (
    echo not found build.json
    goto :EOF
)

echo clean ...
go clean

if not exist go.mod (
    echo golang mod init...
    go mod init
)

::run go generate?
if not exist resource.syso SET GENERATE=1
if not exist version.go SET GENERATE=1
if %RELEASE% equ 1 SET GENERATE=1

::go generate
if %GENERATE% equ 1 (
    :: go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
    echo generate resource...
    go generate
    if %errorlevel% neq 0 (
        echo failed.
        goto :EOF
    )
)

::echo test ...
::go test

echo build ...
if %LINUX% equ 1 (
    ::linux
    set GOARCH=amd64
    set GOOS=linux
    go build -ldflags "-s -w" -a -o %MODULE% .
) else (
    ::windows
    go build -ldflags "-s -w" -o %MODULE% .
)

if %errorlevel% equ 0 (
    echo done.
) else (
    echo failed.
    goto :EOF
)
