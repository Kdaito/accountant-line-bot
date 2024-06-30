#!/bin/bash
# get environmental variables from .env.gcp file
source ./.env.gcp

IMAGE_TAG=$ENV_LOCATION-docker.pkg.dev/$ENV_PROJECT_ID/$ENV_REPOSITORY/$ENV_IMAGE

# # イメージをビルド
docker build -f Dockerfile.deploy -t $IMAGE_TAG .

# # Artifact Registryにプッシュ
docker push $IMAGE_TAG

# Artifact RegistryからCloud Runにデプロイ
# 環境変数はSecret Managerから取得する
gcloud run deploy $ENV_SERVICE \
  --image $IMAGE_TAG \
  --allow-unauthenticated \
  --region $ENV_LOCATION \
  --set-secrets LINE_CHANNEL_SECRET=$SECRET_LINE_CHANNEL_SECRET,LINE_CHANNEL_TOKEN=$SECRET_LINE_CHANNEL_TOKEN,FOLDER_ID=$SECRET_FOLDER_ID,GPT_API_URL=$SECRET_GPT_API_URL,GPT_API_KEY=$SECRET_GPT_API_KEY,SERVICE_ACCOUNT_JSON=$SECRET_SERVICE_ACCOUNT_JSON