#!/bin/bash

cd internal
docker-compose up -d

echo "Wait for 30s pls, just to be sure mySQL is up before starting the app, if mySQL driver error occurs, try running srcipt again =)"

cd ../cmd/redditclone
sleep 30
go run main.go