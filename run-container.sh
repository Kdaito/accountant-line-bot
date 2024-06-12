#!/bin/bash
source ./.env

docker build \
  -f Dockerfile.local \
  -t line-accountant-bot-dev-test \
  --build-arg LINE_CHANNEL_SECRET=$ENV_LINE_CHANNEL_SECRET \
  --build-arg LINE_CHANNEL_TOKEN=$ENV_LINE_CHANNEL_TOKEN \
  --build-arg FOLDER_ID=$ENV_FOLDER_ID \
  --build-arg GPT_API_URL=$ENV_GPT_API_URL \
  --build-arg GPT_API_KEY=$ENV_GPT_API_KEY \
  --build-arg PORT=$ENV_PORT \
  --build-arg CREDENTIALS_JSON=$ENV_CREDENTIALS_JSON \
  --build-arg TOKEN_JSON=$ENV_TOKEN_JSON \
  .
