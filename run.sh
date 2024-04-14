#!/bin/bash

# get environmental variables from .env file
source ./.env
env_line_channel_secret=$ENV_LINE_CHANNEL_SECRET
env_line_channel_token=$ENV_LINE_CHANNEL_TOKEN

# get current environmental variables in the shell
line_channel_secret=$LINE_CHANNEL_SECRET
line_channel_token=$LINE_CHANNEL_TOKEN

# If environment variables are not set in the shell, set variables in the .env file
if [ -z "${line_channel_secret}" ]; then
  export LINE_CHANNEL_SECRET=env_line_channel_secret
  echo "LINE_CHANNEL_SECRET has been set!"
fi

if [ -z "${line_channel_token}" ]; then
  export LINE_CHANNEL_TOKEN=env_line_channel_token
  echo "LINE_CHANNEL_TOKEN has been set!"
fi

# go run ./cmd/api/main.go