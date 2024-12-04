FROM golang:1.23.4-nanoserver AS build-backend

ENV CGO_ENABLED 0
ENV GO111MODULE=on
ENV GOARCH amd64
ENV GOOS windows

COPY ./src /app
RUN cd /app && go get -d \
 && go build \
    -a \
    -ldflags "-s -w -extldflags \"-static\"" \
    -o /bin/npipe.exe main.go

FROM mcr.microsoft.com/windows/nanoserver:ltsc2022

ENV LANG=C.UTF-8

# we don't really need to run this image, but we add a CMD
# to make it run more convenient
CMD /npipe.exe

COPY --from=build-backend /bin/npipe.exe /npipe.exe
