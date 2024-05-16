#!/bin/bash

# get environmental variables from .env file
source ./.env
export LINE_CHANNEL_SECRET=$ENV_LINE_CHANNEL_SECRET
export LINE_CHANNEL_TOKEN=$ENV_LINE_CHANNEL_TOKEN
export DRIVE_FOLDER_ID=$ENV_FOLDER_ID

go run ./cmd/api/main.go