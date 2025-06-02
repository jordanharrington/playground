#!/bin/bash

set -e

MINIO_IMAGE="minio/minio"
MINIO_CONTAINER="objectstore"
MINIO_PORT_API=9000
MINIO_PORT_CONSOLE=9001

MINIO_ACCESS_KEY="admin"
MINIO_SECRET_KEY="secretkey1234"

if [ "$(docker ps -q -f name=$MINIO_CONTAINER)" ]; then
  echo "MinIO container is already running. Skipping startup."
else
  echo "Starting MinIO container..."
  docker run -d \
    --name $MINIO_CONTAINER \
    -p ${MINIO_PORT_API}:9000 \
    -p ${MINIO_PORT_CONSOLE}:9001 \
    -e MINIO_ROOT_USER="$MINIO_ACCESS_KEY" \
    -e MINIO_ROOT_PASSWORD="$MINIO_SECRET_KEY" \
    $MINIO_IMAGE server /data --console-address ":${MINIO_PORT_CONSOLE}"
fi

echo "Applying TF..."
cd terraform/local/
terraform init
terraform apply -auto-approve -var="minio_access_key=$MINIO_ACCESS_KEY" -var="minio_secret_key=$MINIO_SECRET_KEY"