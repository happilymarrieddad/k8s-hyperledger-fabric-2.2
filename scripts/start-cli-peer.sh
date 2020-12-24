#!/bin/bash

function handle() {
    echo "Joining peers to channels"
    bash /scripts/join-cli-peers.sh

    echo "Sleeping for 30 seconds to ensure that everything is good to go"
    sleep 30

    echo "Attaching chaincodes"
    bash /scripts/attach-cli-peer-chaincodes.sh

    echo "Completed"
}

handle &

sleep infinity
