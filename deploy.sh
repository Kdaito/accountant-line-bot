#!/bin/bash

source ./.env

IMAGE_NAME="$ENV_AR_REGION-docker.pkg.dev/$ENV_GCP_PROJECT_ID/$ENV_GCP_REPOSITORY_NAME/$ENV_AR_IMAGE:$ENV_AR_IMAGE_TAG"

docker build -t $IMAGE_NAME \
  --build-arg LINE_CHANNEL_SECRET=$ENV_LINE_CHANNEL_SECRET \
  --build-arg LINE_CHANNEL_TOKEN=$ENV_LINE_CHANNEL_TOKEN \
  --build-arg FOLDER_ID=$ENV_FOLDER_ID \
  --build-arg GPT_API_URL=$ENV_GPT_API_URL \
  --build-arg GPT_API_KEY=$ENV_GPT_API_KEY \
  --build-arg PORT=$ENV_PORT \
  .