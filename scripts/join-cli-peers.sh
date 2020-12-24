#!/bin/bash

sleep 15

# Now we sleep to ensure the cli peers are not running at the exact same time
sleep $((1 + RANDOM % 5));
sleep .$((1 + RANDOM % 5));

# Get the channels this peer should be a part of from the ENV
CHANNEL_LIST=(${CHANNELS//;/ })

# Loop over channel list and join if not a part of channel
for CHANNEL in "${CHANNEL_LIST[@]}"
do
    sleep 1
    echo "Working on channel $CHANNEL"
    TRANSACTION_FILE=channels/${CHANNEL}.tx

    if test -f "$TRANSACTION_FILE"; then

        # This will fail if one of the other peers have already done it but that's not really a big deal
        peer channel create \
                -c ${CHANNEL} \
                -f ${TRANSACTION_FILE} \
                -o orderer0:7050 \
                --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem
        
        # If it got created then make it available for other peers
        if test -f "${CHANNEL}.block"; then
            cp ${CHANNEL}.block channels/${CHANNEL}.block
        fi

        sleep 5

        # The file containing the current channels this peer is a part of
        CURRENT_CHANNEL_FILE=current_channel_list.txt

        # Get the list of channels this peer is a part of
        peer channel list > ${CURRENT_CHANNEL_FILE}

        # Checking if peer is a part of this channel
        if grep -q $CHANNEL "$CURRENT_CHANNEL_FILE"; then
            echo "Peer has already joined $CHANNEL so just continuing on"
        else
            peer channel join -b channels/${CHANNEL}.block
        fi

    else
        echo "Transaction file $TRANSACTION_FILE does not exist so ignoring channel"
    fi

    # Now we handle the anchors. It's okay if it fails.. no harm no foul
    # TODO: Check if update is necessary
    if [[ "$CORE_PEER_ID" == *"peer0"* ]]; then
        echo "This is an anchor peer so updating the channel -c ${CHANNEL} -f channels/${CORE_PEER_LOCALMSPID}-anchors.tx"
        sleep 5
        peer channel update \
            -o orderer0:7050 --tls \
            --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem \
            -c ${CHANNEL} -f channels/${CORE_PEER_LOCALMSPID}-anchors.tx
    fi

done

echo "Done joining"
