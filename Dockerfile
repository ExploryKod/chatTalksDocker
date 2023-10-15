# Stage 1: Build the Go application
FROM golang:1.18 as BUILDER

# Active le comportement de module indépendant
ENV GO111MODULE=on

# Je vais faire une build en 2 étapes
# https://dave.cheney.net/2016/01/18/cgo-is-not-go
ENV CGO_ENABLED=0
ENV GOOS=darwin
ENV GOARCH=arm64

WORKDIR /go_app
COPY ./go_app .
RUN go mod download \
    && go mod verify \
    && go build -o /build/buildedApp main/main.go

# Stage 2: Create the final image
FROM mariadb:latest

# Définir les variables d'environnement pour la base de données
ENV MYSQL_ROOT_PASSWORD=password
ENV MYSQL_DATABASE=chatbdd

# Copier le binaire de l'application Go depuis le stage précédent
COPY --from=BUILDER /build/buildedApp /app/buildedApp

# Exposer le port de l'application Go
EXPOSE 8000

# Commande d'entrée pour exécuter l'application Go
CMD ["/app/buildedApp"]



