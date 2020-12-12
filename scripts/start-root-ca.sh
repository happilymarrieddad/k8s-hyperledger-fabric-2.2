#!/bin/bash

fabric-ca-server start \
    -b ${USERNAME}:${PASSWORD} \
    --tls.enabled \
    --csr.hosts ${CSR_HOSTS},${CSR_HOSTS}-service
