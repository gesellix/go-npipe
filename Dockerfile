FROM golang:1.6-alpine
MAINTAINER Tobias Gesellchen <tobias@gesellix.de> (@gesellix)

RUN apk add --no-cache git && \
    mkdir /go/src/github.com/gesellix/go-npipe/

WORKDIR /go/src/github.com/gesellix/go-npipe/

# we don't really need to run this image
CMD /npipe.exe

ENV CGO_ENABLED 0
ENV GOARCH amd64
ENV GOOS windows

COPY *.go /go/src/github.com/gesellix/go-npipe/
RUN go get ./... && go build -o /npipe.exe
