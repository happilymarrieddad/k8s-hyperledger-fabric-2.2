#!/bin/bash

# Ensure Deps
apk update
apk add jq

# Get the channels this peer should be a part of from the ENV
CHANNEL_LIST=(${CHANNELS//;/ })

# Loop over channel list and join if not a part of channel
for CHANNEL in "${CHANNEL_LIST[@]}"
do
    sleep 1
    echo "Working on channel $CHANNEL"

    # Lets go ahead and get what's installed
    ICFILE=installed_chaincodes.txt
    peer lifecycle chaincode queryinstalled -O json > ${ICFILE}

    # Now time to handle chaincodes in this channel
    CHAINCODE_LIST=(${CHAINCODES//;/ })

    for CHAINCODE in "${CHAINCODE_LIST[@]}"
    do
        SEQ=0
        # Get the current version of the chaincode
        VERSION=$(jq ".installed_chaincodes[].references.\"$CHANNEL\".chaincodes[] | select(.name == \"$CHAINCODE\").version" ${ICFILE})
        if [ -z "$VERSION" ]; then
            # Set the version
            SEQ=1
        else
            # Chaincode exists and has a version
            SEQ=${VERSION%.*}
            SEQ=${SEQ:1}
            SEQ=$((SEQ+1))
        fi

        # Package the chaincode
        peer lifecycle chaincode package ${CHAINCODE}.tar.gz \
            --path /opt/gopath/src/${CHAINCODE} \
            --lang golang \
            --label ${CHAINCODE}_${SEQ}

        # Use the package to install the chaincode and save the unique ID in the text file
        peer lifecycle chaincode install ${CHAINCODE}.tar.gz &> pkg.txt

        sleep 5
        if [[ "$CORE_PEER_ID" == *"peer0"* ]]; then
            echo "This is an anchor peer so updating the channel -c ${CHANNEL} -f channels/${CORE_PEER_LOCALMSPID}-anchors.tx"
            sleep 10
            peer lifecycle chaincode approveformyorg \
                -o orderer0:7050 --tls \
                --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem \
                --channelID ${CHANNEL} \
                --name ${CHAINCODE} \
                --collections-config /opt/gopath/src/${CHAINCODE}/collections-config.json \
                --version ${SEQ}.0 \
                --sequence ${SEQ} \
                --package-id $(tail -n 1 pkg.txt | awk 'NF>1{print $NF}')

            # Need a longer delay just in case...
            echo "Waiting to commit chaincode..."
            sleep 10
            # Only 1 peer needs to be commiting the chaincode
            if [[ "$CORE_PEER_LOCALMSPID" == *"ibm"* ]]; then
                peer lifecycle chaincode commit \
                    -o orderer0:7050 --tls \
                    --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem \
                    --collections-config /opt/gopath/src/${CHAINCODE}/collections-config.json \
                    --channelID ${CHANNEL} \
                    --name ${CHAINCODE} \
                    --version ${SEQ}.0 --sequence ${SEQ}
            fi

        fi
    done

done

echo "Done installing chaincodes"
