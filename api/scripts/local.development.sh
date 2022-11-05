#!/bin/bash

export COMPOSE_IGNORE_ORPHANS=True

export BACKEND_IMAGE=farmacare-go
export BACKEND_IMAGE_TAG=development
export BACKEND_CONTAINER=farmacare-go-development
export BACKEND_HOST=farmacare-go.service
export BACKEND_STAGE=development

docker build -t "$BACKEND_IMAGE:$BACKEND_IMAGE_TAG" .
docker-compose -f ./manifest/docker-compose.development.yaml up -d --build
