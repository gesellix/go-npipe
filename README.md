# Windows Named Pipes helper

## Build

    CGO_ENABLED=0 GOARCH=386 GOOS=windows go build -o npipe.exe main.go 

## Run

    npipe \\.\pipe\the_pipe

## List pipes (PowerShell)

    get-childitem \\.\pipe\
