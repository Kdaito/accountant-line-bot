#!/bin/bash

# get environmental variables from .env file
source ./.env
env_line_channel_secret=$ENV_LINE_CHANNEL_SECRET
env_line_channel_token=$ENV_LINE_CHANNEL_TOKEN
env_drive_folder_id=$ENV_FOLDER_ID

# get current environmental variables in the shell
line_channel_secret=$LINE_CHANNEL_SECRET
line_channel_token=$LINE_CHANNEL_TOKEN
drive_folder_id=$DRIVE_FOLDER_ID

# If environment variables are not set in the shell, set variables in the .env file
if [ -z "${line_channel_secret}" ]; then
  export LINE_CHANNEL_SECRET=$env_line_channel_secret
  echo "LINE_CHANNEL_SECRET has been set!"
fi

if [ -z "${line_channel_token}" ]; then
  export LINE_CHANNEL_TOKEN=$env_line_channel_token
  echo "LINE_CHANNEL_TOKEN has been set!"
fi

if [ -z "${drive_folder_id}" ]; then
  export DRIVE_FOLDER_ID=$env_drive_folder_id
  echo "DRIVE_FOLDER_ID has been set!"
fi

go run ./cmd/api/main.go