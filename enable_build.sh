#!/bin/bash

if [ "$1" == "start" ]; then
  # If the argument is "start", bring the containers down and then up
  docker-compose down --remove-orphans
  docker-compose up -d --build
fi

docker exec -it go-api sh -c "go build -o go-api . && ./go-api"
