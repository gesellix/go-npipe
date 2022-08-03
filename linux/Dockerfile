FROM golang:1.19.0-alpine AS build-backend

ENV CGO_ENABLED 0
ENV GO111MODULE=on
ENV GOARCH amd64
ENV GOOS windows

COPY ./src /app
RUN cd /app && go get -d \
 && go build \
    -a \
    -ldflags '-s -w -extldflags "-static"' \
    -o /npipe.exe main.go

FROM alpine:3.16.1

ENV LANG=C.UTF-8

# we don't really need to run this image, but we add a CMD
# to make it run more convenient
CMD /npipe.exe

COPY --from=build-backend /npipe.exe /npipe.exe
