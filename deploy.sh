#!/bin/bash

# Using git to version control the images
VERSION=$(git rev-parse HEAD)

# These are optional arguments
# usage - ./deploy.sh <docker-username> <go-api-name> <node-api-name> <web-app-name>
# example - ./deploy.sh someuser go-api node-api web-app
DOCKER_USERNAME=${1-happilymarrieddadudemy}
GO_API_REPO=${2-udemy-go-api}
NODE_API_REPO=${3-udemy-node-api}
WEB_APP_REPO=${4-udemy-web-app}

## Push up and set GO API
docker build -t ${DOCKER_USERNAME}/${GO_API_REPO}:${VERSION} \
    -f ./go-api/Dockerfile \
    ./go-api

docker push ${DOCKER_USERNAME}/${GO_API_REPO}:${VERSION}

docker build -t ${DOCKER_USERNAME}/${GO_API_REPO}:latest \
    -f ./go-api/Dockerfile \
    ./go-api

docker push ${DOCKER_USERNAME}/${GO_API_REPO}:latest

# sleep 10

#kubectl set image deployments/api-deployment api=${DOCKER_USERNAME}/${GO_API_REPO}:${VERSION}


# sleep 5

## Push up and set Node API
docker build -t ${DOCKER_USERNAME}/${NODE_API_REPO}:${VERSION} \
    -f ./node-api/Dockerfile \
    ./node-api

docker push ${DOCKER_USERNAME}/${NODE_API_REPO}:${VERSION}

docker build -t ${DOCKER_USERNAME}/${NODE_API_REPO}:latest \
    -f ./node-api/Dockerfile \
    ./node-api

docker push ${DOCKER_USERNAME}/${NODE_API_REPO}:latest

# sleep 10

#kubectl set image deployments/api-deployment api=${DOCKER_USERNAME}/${NODE_API_REPO}:${VERSION}


# sleep 5


## Push up web app production
docker build -t ${DOCKER_USERNAME}/${WEB_APP_REPO}:${VERSION} \
    -f ./frontend/Dockerfile \
    ./frontend

docker push ${DOCKER_USERNAME}/${WEB_APP_REPO}:${VERSION}

docker build -t ${DOCKER_USERNAME}/${WEB_APP_REPO}:latest \
    -f ./frontend/Dockerfile \
    ./frontend

docker push ${DOCKER_USERNAME}/${WEB_APP_REPO}:latest

## Push up web app local
docker build -t ${DOCKER_USERNAME}/${WEB_APP_REPO}:${VERSION}-local \
    -f ./frontend/Dockerfile.dev \
    ./frontend

docker push ${DOCKER_USERNAME}/${WEB_APP_REPO}:${VERSION}-local

docker build -t ${DOCKER_USERNAME}/${WEB_APP_REPO}:latest-local \
    -f ./frontend/Dockerfile.dev \
    ./frontend

docker push ${DOCKER_USERNAME}/${WEB_APP_REPO}:latest-local

# sleep 10

#kubectl set image deployments/web-deployment web=${DOCKER_USERNAME}/${WEB_APP_REPO}:${VERSION}

sleep 1

echo "Completed!"
