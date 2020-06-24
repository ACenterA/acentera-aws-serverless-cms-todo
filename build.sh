#!/bin/bash

BUILDNAME=aws-graphql-backend
export COMPOSE_PROJECT_NAME=acentera-core
docker build -t ${COMPOSE_PROJECT_NAME}_${BUILDNAME}:latest -f ./backend/go/Dockerfile-build .

docker tag ${COMPOSE_PROJECT_NAME}_${BUILDNAME}:latest acentera/dev:${BUILDNAME}-v0.0.1
docker tag ${COMPOSE_PROJECT_NAME}_${BUILDNAME}:latest registry.gitlab.com/acentera/docker-shared/tools/${BUILDNAME}-v0.0.1
#docker push registry.gitlab.com/acentera/docker-shared/tools/${BUILDNAME}-v0.0.1


BUILDNAME=aws-graphql
docker build -t ${COMPOSE_PROJECT_NAME}_${BUILDNAME}:latest -f ./Dockerfile-build-graph .
docker tag ${COMPOSE_PROJECT_NAME}_${BUILDNAME}:latest acentera/dev:${BUILDNAME}-v0.0.1
docker tag ${COMPOSE_PROJECT_NAME}_${BUILDNAME}:latest registry.gitlab.com/acentera/docker-shared/tools/${BUILDNAME}-v0.0.1
#docker push registry.gitlab.com/acentera/docker-shared/tools/${BUILDNAME}-v0.0.1

#      - ./backend/:/usr/app/backend/
#      - ./backend/go/schema.graphql:/usr/app/backend/schema.graphql
#      - ./backend/go/template.yml:/usr/app/backend/template.yml
#      - ./src/gql/:/usr/app/gql

