#!/bin/bash

if [ "$1" == "start" ]; then
  # If the argument is "start", bring the containers down and then up
  docker-compose down --remove-orphans
  docker-compose up -d --build
fi

if [ "$1" == "clean" ]; then
  # If the argument is "clean", clean the host docker and build all
  docker system prune
  docker volume prune
  docker image  prune
  docker-compose down --remove-orphans
  docker-compose up -d --build
fi

docker exec -it go-api sh -c "go run main.go"