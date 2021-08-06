Docker Compose portion of the Readme
=====================================

## Docker (Local)

Okay, lets generate the certs for the network
```bash
docker-compose -f network/docker/docker-compose-ca.yaml up
```

After it's done close down the network. Now it's time to generate the network artifacts
```bash
sudo chmod 777 -R crypto-config
sudo chown $USER:$USER -R crypto-config

configtxgen -profile OrdererGenesis -channelID syschannel -outputBlock ./orderer/genesis.block
configtxgen -profile MainChannel -outputCreateChannelTx ./channels/mainchannel.tx -channelID mainchannel
configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channels/ibm-anchors.tx -channelID mainchannel -asOrg ibm
configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channels/oracle-anchors.tx -channelID mainchannel -asOrg oracle
```

Now it's time to start the network
- If OSX - need to jump to minikube for now. I can't find a way to make it work with Catalina. I'll keep playing around with it. It's a problem with the docker.sock file and not being able to mount it to the peer. It's needed to spin up containers for the chaincode.
```bash
docker-compose -f network/docker/docker-compose.yaml up
```

Lets setup the artifacts
```bash
docker exec -it cli-peer0-ibm bash -c 'peer channel create -c mainchannel -f ./channels/mainchannel.tx -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'

docker exec -it cli-peer0-ibm bash -c 'cp mainchannel.block ./channels/'
docker exec -it cli-peer0-ibm bash -c 'peer channel join -b channels/mainchannel.block'
docker exec -it cli-peer1-ibm bash -c 'peer channel join -b channels/mainchannel.block'
docker exec -it cli-peer0-oracle bash -c 'peer channel join -b channels/mainchannel.block'
docker exec -it cli-peer1-oracle bash -c 'peer channel join -b channels/mainchannel.block'

sleep 5

docker exec -it cli-peer0-ibm bash -c 'peer channel update -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem -c mainchannel -f channels/ibm-anchors.tx'
docker exec -it cli-peer0-oracle bash -c 'peer channel update -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem -c mainchannel -f channels/oracle-anchors.tx'
```

Now we are going to install the chaincode
- Make sure you go mod vendor in each chaincode folder... might need to remove the go.sum depending
```bash
docker exec -it cli-peer0-ibm bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1' &
docker exec -it cli-peer1-ibm bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1' &
docker exec -it cli-peer0-oracle bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1' &
docker exec -it cli-peer1-oracle bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1'

docker exec -it cli-peer0-ibm bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt' &
docker exec -it cli-peer1-ibm bash -c 'peer lifecycle chaincode install resource_types.tar.gz' &
docker exec -it cli-peer0-oracle bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt' &
docker exec -it cli-peer1-oracle bash -c 'peer lifecycle chaincode install resource_types.tar.gz'

docker exec -it cli-peer0-ibm bash -c 'peer lifecycle chaincode approveformyorg -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')' &
docker exec -it cli-peer0-oracle bash -c 'peer lifecycle chaincode approveformyorg -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'

docker exec -it cli-peer0-ibm bash -c 'peer lifecycle chaincode commit -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 1.0 --sequence 1'
```

Lets go ahead and test this chaincode
- If OSX
1. eval $(docker-machine env)
```bash
docker exec -it cli-peer0-ibm bash -c 'peer chaincode invoke -C mainchannel -n resource_types -c '\''{"Args":["Create", "1","Parts"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
sleep 5
docker exec -it cli-peer0-ibm bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-ibm bash -c 'peer chaincode invoke -C mainchannel -n resource_types -c '\''{"Args":["Update", "1", "Parts 2"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-ibm bash -c 'peer chaincode invoke -C mainchannel -n resource_types -c '\''{"Args":["Update", "1", "Parts"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-ibm bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Transactions","1"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
```

Lets try the other chaincode
- If OSX
1. eval $(docker-machine env)
```bash
docker exec -it cli-peer0-ibm bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1' &
docker exec -it cli-peer1-ibm bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1' &
docker exec -it cli-peer0-oracle bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1' &
docker exec -it cli-peer1-oracle bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1'


docker exec -it cli-peer0-ibm bash -c 'peer lifecycle chaincode install resources.tar.gz &> pkg.txt' &
docker exec -it cli-peer1-ibm bash -c 'peer lifecycle chaincode install resources.tar.gz' &
docker exec -it cli-peer0-oracle bash -c 'peer lifecycle chaincode install resources.tar.gz &> pkg.txt' &
docker exec -it cli-peer1-oracle bash -c 'peer lifecycle chaincode install resources.tar.gz'

docker exec -it cli-peer0-ibm bash -c 'peer lifecycle chaincode approveformyorg -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem --collections-config /opt/gopath/src/resources/collections-config.json --channelID mainchannel --name resources --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')' &
docker exec -it cli-peer0-oracle bash -c 'peer lifecycle chaincode approveformyorg -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem --collections-config /opt/gopath/src/resources/collections-config.json --channelID mainchannel --name resources --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'

docker exec -it cli-peer0-ibm bash -c 'peer lifecycle chaincode commit -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem --collections-config /opt/gopath/src/resources/collections-config.json --channelID mainchannel --name resources --version 1.0 --sequence 1'

sleep 5

docker exec -it cli-peer0-ibm bash -c 'peer chaincode invoke -C mainchannel -n resources -c '\''{"Args":["Create","1","CPUs","1"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-oracle bash -c 'peer chaincode invoke -C mainchannel -n resources -c '\''{"Args":["Create","2","Database Servers","1"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
sleep 5
docker exec -it cli-peer0-ibm bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer1-ibm bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-oracle bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer1-oracle bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
```

Okay, now everything should be working as normal. Lets test the apis and make sure they are connected properly.
```bash
cd node-api
node index.js
```

Start the GO api in a different terminal
```bash
cd go-api
go run main.go
``` 

In a third terminal test the apis
```bash
curl localhost:3000/v1/resources
curl localhost:4001/resources
```

You should see this
```bash
☁  k8s-hyperledger-fabric-2.2 [master] ⚡  curl localhost:3000/v1/resources
[{"id":"1","name":"CPUs","resource_type_id":"1","active":true},{"id":"2","name":"Database Servers","resource_type_id":"1","active":true}]%                                                                                           ☁  k8s-hyperledger-fabric-2.2 [master] ⚡  curl localhost:4001/resources
[{"id":"1","name":"CPUs","resource_type_id":"1","active":true},{"id":"2","name":"Database Servers","resource_type_id":"1","active":true}]%                                                                                           ☁  k8s-hyperledger-fabric-2.2 [master] ⚡  
```

Lets run the front end and test it
```bash
cd frontend
npm run serve
```

Everything should work!

