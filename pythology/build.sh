#!/bin/bash

BUILDER_NAME="pythology-builder"
REGISTRY="ghcr.io/jordanharrington/playground"
PLATFORMS="linux/arm64,linux/amd64"

init_buildx() {
  if ! docker buildx inspect "$BUILDER_NAME" > /dev/null 2>&1; then
    echo "Creating new buildx builder: $BUILDER_NAME"
    docker buildx create --name "$BUILDER_NAME" --use
  else
    echo "Using existing buildx builder: $BUILDER_NAME"
    docker buildx use "$BUILDER_NAME"
  fi
}

build_extractor() {
  init_buildx
  source ./pipeline/extract/version.sh
  docker buildx build \
  --progress=plain \
  --platform "$PLATFORMS" \
  --build-arg PYTHON="$PYTHON_VERSION" \
  --build-arg API_VERSION="$API_VERSION" \
  -f pipeline/extract/Dockerfile \
  -t "$REGISTRY"/pythology-extractor:dev-"$SEMVER" \
  -t "$REGISTRY"/pythology-extractor:dev-latest \
  . --push || true
}