#!/bin/bash

. /scripts/start-client-ca.sh

sleep 10

. /scripts/create-org-peer-certs.sh \
    ${ORG_NAME} \
    ${CA_SCHEME} \
    ${CA_USERNAME} \
    ${CA_PASSWORD} \
    ${CA_URL} \
    ${CA_CERT_PATH} \
    2

sleep infinity
