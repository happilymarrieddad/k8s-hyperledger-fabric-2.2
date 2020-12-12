#!/bin/bash

# Using git to version control the images
VERSION=$(git rev-parse HEAD)
DOCKER_USERNAME=${1-happilymarrieddadudemy}

## Push up and set GO API
docker build -t ${DOCKER_USERNAME}/udemy-go-api:${VERSION} \
    -f ./go-api/Dockerfile \
    ./go-api

docker push ${DOCKER_USERNAME}/udemy-go-api:${VERSION}

docker build -t ${DOCKER_USERNAME}/udemy-go-api:latest \
    -f ./go-api/Dockerfile \
    ./go-api

docker push ${DOCKER_USERNAME}/udemy-go-api:latest

# sleep 10

#kubectl set image deployments/api-deployment api=${DOCKER_USERNAME}/udemy-go-api:${VERSION}


# sleep 5

## Push up and set Node API
docker build -t ${DOCKER_USERNAME}/udemy-node-api:${VERSION} \
    -f ./node-api/Dockerfile \
    ./node-api

docker push ${DOCKER_USERNAME}/udemy-node-api:${VERSION}

docker build -t ${DOCKER_USERNAME}/udemy-node-api:latest \
    -f ./node-api/Dockerfile \
    ./node-api

docker push ${DOCKER_USERNAME}/udemy-node-api:latest

# sleep 10

#kubectl set image deployments/api-deployment api=${DOCKER_USERNAME}/udemy-node-api:${VERSION}


# sleep 5


## Push up web app production
docker build -t ${DOCKER_USERNAME}/udemy-web-app:${VERSION} \
    -f ./frontend/Dockerfile \
    ./frontend

docker push ${DOCKER_USERNAME}/udemy-web-app:${VERSION}

docker build -t ${DOCKER_USERNAME}/udemy-web-app:latest \
    -f ./frontend/Dockerfile \
    ./frontend

docker push ${DOCKER_USERNAME}/udemy-web-app:latest

## Push up web app local
docker build -t ${DOCKER_USERNAME}/udemy-web-app:${VERSION}-local \
    -f ./frontend/Dockerfile.dev \
    ./frontend

docker push ${DOCKER_USERNAME}/udemy-web-app:${VERSION}-local

docker build -t ${DOCKER_USERNAME}/udemy-web-app:latest-local \
    -f ./frontend/Dockerfile.dev \
    ./frontend

docker push ${DOCKER_USERNAME}/udemy-web-app:latest-local

# sleep 10

#kubectl set image deployments/web-deployment web=${DOCKER_USERNAME}/udemy-web-app:${VERSION}

sleep 1

echo "Completed!"
