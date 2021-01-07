Updating Chaincode
=============================

Lets go ahead and update this chaincode
```bash
docker exec -it cli-peer0-ibm bash -c 'peer chaincode invoke -C mainchannel -n resource_types -c '\''{"Args":["Create", "1","Parts"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
sleep 5
docker exec -it cli-peer0-ibm bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
```

Lets try the other chaincode
```bash
docker exec -it cli-peer0-ibm bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_3' &
docker exec -it cli-peer1-ibm bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_3' &
docker exec -it cli-peer0-oracle bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_3' &
docker exec -it cli-peer1-oracle bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_3'


docker exec -it cli-peer0-ibm bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt' &
docker exec -it cli-peer1-ibm bash -c 'peer lifecycle chaincode install resource_types.tar.gz' &
docker exec -it cli-peer0-oracle bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt' &
docker exec -it cli-peer1-oracle bash -c 'peer lifecycle chaincode install resource_types.tar.gz'


docker exec -it cli-peer0-ibm bash -c 'peer lifecycle chaincode approveformyorg -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 3.0 --sequence 3 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')' &
docker exec -it cli-peer0-oracle bash -c 'peer lifecycle chaincode approveformyorg -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 3.0 --sequence 3 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'


docker exec -it cli-peer0-ibm bash -c 'peer lifecycle chaincode commit -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 3.0 --sequence 3'

sleep 5
docker exec -it cli-peer0-ibm bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-oracle bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'

sleep 5
docker exec -it cli-peer0-ibm bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Update", "e7f48784-0702-4ded-801a-46bb4d727957", "CPUs", "1", "oracle"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'

sleep 5
docker exec -it cli-peer0-ibm bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-oracle bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
```