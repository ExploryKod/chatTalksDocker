#!/bin/bash

if [ "$1" == "start" ]; then
  # If the argument is "start", bring the containers down and then up
  docker-compose down
  docker-compose up -d --build
fi

docker exec -it go-api sh -c "go run main.go"