# Windows Named Pipe Http Echo Server.

## About

This little tool creates a Windows Named Pipe and listens for HTTP requests.
Request bodies will be returned with a prefix `[echo]`.

Requests on `/exit` will stop the tool and remove the pipe.

The request method will be ignored, so you can use any of `GET`, `POST` or whatever.

### Usage in integration tests

This tool comes in handy when running integration tests. There's no built-in way to _create_ a named pipe in the JVM,
 so you end up trying to create some JNI/JNA implementation
 to make [the syscall](https://msdn.microsoft.com/en-us/library/windows/desktop/aa365150(v=vs.85).aspx) for you...
 or you use a native binary to handle everything for you and simply call it via `exec`.

You can find my own use case in the [integration tests for the Docker Client](https://github.com/gesellix/docker-client/blob/94ea21ff5620235d51a2adbc4b4106d55e0b0887/client/src/integrationTest/groovy/de/gesellix/docker/client/filesocket/HttpOverNamedPipeIntegrationTest.groovy#L55).
As soon as the project is built on a Windows system, the integration tests verify
 the basic ability to use a named pipe socket for communication with the Docker engine.  

### Yeah, Windows only, but who cares? ¯\_(ツ)_/¯

Since Windows Named Pipes are a Windows only concept (similar to Unix Domain Sockets), you would probably expect
 Windows specific sources or build configurations. Thanks to Golang this isn't necessary: the cross compiler 
 works very well and creates a nice Windows native executable. I'm a fan of automation, which is why
 you'll always find the most recent version of this tool on the Docker Hub, packaged as Docker image.
 The Docker image isn't expected to be run, but it serves to leverage the build _and_ distribution of this tool.
 See below for details!

## Build/Install

### Command Line

    go get -d github.com/gesellix/go-npipe
    cd src
    CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o npipe.exe main.go

### Docker :whale:

    docker build -t npipe:local -f linux/Dockerfile .
    docker create --name npipe npipe:local
    # Or from the Docker Hub (linux):
    #docker create --name npipe gesellix/npipe
    # Or from the Docker Hub (windows):
    #docker create --name npipe gesellix/npipe:windows
    docker cp npipe:/npipe.exe .
    docker rm npipe

## Run

    npipe.exe \\.\pipe\the_pipe

## List pipes (PowerShell)

    get-childitem \\.\pipe\

### Contributing/Feedback

If you have any kind of feedback or would like to improve this little tool,
 please [submit an issue](https://github.com/gesellix/go-npipe/issues) or [create a pull request](https://github.com/gesellix/go-npipe/pulls)!
