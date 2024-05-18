#!/bin/bash

# get environmental variables from .env file
source ./.env
export LINE_CHANNEL_SECRET=$ENV_LINE_CHANNEL_SECRET
export LINE_CHANNEL_TOKEN=$ENV_LINE_CHANNEL_TOKEN
export DRIVE_FOLDER_ID=$ENV_FOLDER_ID
export GPT_API_URL=$ENV_GPT_API_URL
export GPT_API_KEY=$ENV_GPT_API_KEY

go run ./cmd/api/main.go