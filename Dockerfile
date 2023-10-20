FROM golang:1.21

# Active le comportement de module indépendant
ENV GO111MODULE=on

# Je vais faire une build en 2 étapes
# https://dave.cheney.net/2016/01/18/cgo-is-not-go
ENV CGO_ENABLED=0
ENV GOOS=$GOOS
ENV GOARCH=$GOARCH

WORKDIR /go_app
COPY ./go_app .

RUN go install github.com/cosmtrek/air@latest

RUN go mod download \
    && go mod verify \
    && go mod tidy


