# Windows Named Pipe Http Echo Server.

## About

This little tool creates a Windows Named Pipe and listens for HTTP requests.
Request bodies will be returned with a prefix `[echo]`.

Requests on `/exit` will stop the tool and remove the pipe.

The request method will be ignored, so you can use any of `GET`, `POST` or whatever.

## Build/Install

### Command Line

    go get ./...
    CGO_ENABLED=0 GOARCH=386 GOOS=windows go build -o npipe.exe main.go
     
### Docker

    docker create --name npipe gesellix/npipe
    docker cp npipe:/npipe.exe .
    docker rm npipe

## Run

    npipe.exe \\.\pipe\the_pipe

## List pipes (PowerShell)

    get-childitem \\.\pipe\
