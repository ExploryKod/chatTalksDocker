#!/bin/bash

if [ "$1" == "start" ]; then
  # If the argument is "start", bring the containers up only with build
  docker compose up -d --build
fi

if [ "$1" == "restart" ]; then
  # If the argument is "start", bring the containers down and then up
  docker compose down --remove-orphans
  docker compose up -d --build
fi

if [ "$1" == "clean" ]; then
  # If the argument is "clean", clean the host docker and build all
  docker system prune
  docker volume prune
  docker image  prune
  docker-compose down --remove-orphans
  docker-compose up -d --build
fi

# start air directly once docker container is set up
if [ "$1" == "air" ]; then
docker exec -it go-api sh -c "go install github.com/cosmtrek/air@latest"

docker exec -it go-api sh -c "air init"

# Modify .air.toml to set poll = true
docker exec -it go-api sed -i 's/poll = false/poll = true/' .air.toml

docker exec -it go-api sh -c "air -c .air.toml"
fi

# start only air init once docker container is set up
if [ "$1" == "air-init" ]; then
docker exec -it go-api sh -c "air init"

# Modify .air.toml to set poll = true
docker exec -it go-api sed -i 's/poll = false/poll = true/' .air.toml

docker exec -it go-api sh -c "air -c .air.toml"
fi

# start air with installing it and build
if [ "$2" == "air" ]; then
docker exec -it go-api sh -c "go install github.com/cosmtrek/air@latest"

docker exec -it go-api sh -c "air init"

# Modify .air.toml to set poll = true
docker exec -it go-api sed -i 's/poll = false/poll = true/' .air.toml

docker exec -it go-api sh -c "air -c .air.toml"
fi

# start only air init
if [ "$2" == "air-init" ]; then
docker exec -it go-api sh -c "air init"

# Modify .air.toml to set poll = true
docker exec -it go-api sed -i 's/poll = false/poll = true/' .air.toml

docker exec -it go-api sh -c "air -c .air.toml"
fi

# Restart air without build directly if a container is already up
if [ "$1" == "restart-air" ]; then
docker exec -it go-api sh -c "air -c .air.toml"
fi

# Restart air without build
if [ "$1" == "restart-air" ]; then
docker exec -it go-api sh -c "air -c .air.toml"
fi