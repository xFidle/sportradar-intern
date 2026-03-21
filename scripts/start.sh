#!/usr/bin/bash

if [[ ! -d "../storage" ]]; then
   sleep 1 # placeholder, will download data later on 
fi

echo "Starting containers..."
docker compose up
