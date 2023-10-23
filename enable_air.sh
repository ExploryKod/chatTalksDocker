#!/bin/bash

docker exec -it go-api sh -c "go install github.com/cosmtrek/air@latest"

docker exec -it go-api sh -c "air init"

# Modify .air.toml to set poll = true
docker exec -it go-api sed -i 's/poll = false/poll = true/' .air.toml

docker exec -it go-api sh -c "air -c .air.toml"