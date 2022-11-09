#!/bin/bash

cd internal
docker-compose up -d

echo "Wait for 30s pls , to make sure mySQL is up before starting the app"

cd ../cmd/redditclone
sleep 30
go run main.go