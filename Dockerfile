FROM golang:1.11-alpine
LABEL maintainer="Tobias Gesellchen <tobias@gesellix.de> (@gesellix)"

RUN apk add --no-cache git && \
    mkdir -p /go/src/github.com/gesellix/go-npipe/

WORKDIR /go/src/github.com/gesellix/go-npipe/

# we don't really need to run this image, but we add a CMD
# to make it run more convenient
CMD /npipe.exe

ENV CGO_ENABLED 0
ENV GOARCH amd64
ENV GOOS windows

COPY *.go /go/src/github.com/gesellix/go-npipe/
RUN go get ./... && go build -o /npipe.exe
