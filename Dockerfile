FROM golang:1.21

# Active le comportement de module indépendant
ENV GO111MODULE=on

# Je vais faire une build en 2 étapes
# https://dave.cheney.net/2016/01/18/cgo-is-not-go
ENV CGO_ENABLED=0
ENV GOOS=$GOOS
ENV GOARCH=$GOARCH

WORKDIR /gorillachat
COPY ./gorillachat .

#RUN go install github.com/cosmtrek/air@latest
#RUN apt-get update
#RUN apt-get install nano -y

#ENV MYSQL_ADDON_HOST=bsgmzsx3etgjzeywrgbb-mysql.services.clever-cloud.com:3306
#ENV MYSQL_ADDON_DB=bsgmzsx3etgjzeywrgbb
#ENV MYSQL_ADDON_USER=udds4bjysjqdatqk
#ENV MYSQL_ADDON_PORT=3306
#ENV MYSQL_ADDON_PASSWORD=n4dmfb2mx7iaZlQsidHe
#ENV MYSQL_ADDON_URI=mysql://udds4bjysjqdatqk:n4dmfb2mx7iaZlQsidHe@bsgmzsx3etgjzeywrgbb-mysql.services.clever-cloud.com:3306/bsgmzsx3etgjzeywrgbb

RUN go mod download \
    && go mod verify \
    && go mod tidy \
    && go build -o gorillachat

