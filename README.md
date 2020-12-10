K8s Hyperledger 2.2 Network
======================================

## Getting Started
- [My bash](https://github.com/ohmyzsh/ohmyzsh)
- [Install Hyperledger Deps](https://hyperledger-fabric.readthedocs.io/en/release-2.2/install.html)
- [Create aws account](aws.amazon.com)
- [Create docker hub account](https://hub.docker.com/)
- curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.2.1 1.4.9
- [Install Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Install Minikube](https://minikube.sigs.k8s.io/docs/start/)

You want to copy over all the bin files into a bin directory where you can easily access them. I usually just like to make a bin folder in my home directory and set the path to include that
```bash
mkdir -p ~/bin
cp fabric-samples/bin/* ~/bin
```

You then need to add to your bash the path and then source it. I use ZSH so this is how I do it. Copy this into your bash file
```bash
export PATH=$PATH:~/bin
```

Then source your bash so they are available. Check the execs and make sure they work
```bash
source ~/.zshrc
```

This is what you should see
```bash
☁  k8s-hyperledger-fabric-2.2 [master] ⚡  configtxgen --version
configtxgen:
 Version: 2.2.1
 Commit SHA: 344fda602
 Go version: go1.14.4
 OS/Arch: linux/amd64
☁  k8s-hyperledger-fabric-2.2 [master] ⚡  
```

You then need to install Go and Nodejs... I would also suggest installing vue cli but I've included all the files you need for the front end.
- [Golang](https://golang.org/dl/)
- [Nodejs](https://nodejs.org/en/)

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
configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channels/org1-anchors.tx -channelID mainchannel -asOrg org1
configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channels/org2-anchors.tx -channelID mainchannel -asOrg org2
```

Now it's time to start the network
```bash
docker-compose -f network/docker/docker-compose.yaml up
```

Lets setup the artifacts
```bash
docker exec -it cli-peer0-org1 bash -c 'peer channel create -c mainchannel -f ./channels/mainchannel.tx -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'

docker exec -it cli-peer0-org1 bash -c 'cp mainchannel.block ./channels/'
docker exec -it cli-peer0-org1 bash -c 'peer channel join -b channels/mainchannel.block'
docker exec -it cli-peer1-org1 bash -c 'peer channel join -b channels/mainchannel.block'
docker exec -it cli-peer0-org2 bash -c 'peer channel join -b channels/mainchannel.block'
docker exec -it cli-peer1-org2 bash -c 'peer channel join -b channels/mainchannel.block'

sleep 5

docker exec -it cli-peer0-org1 bash -c 'peer channel update -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem -c mainchannel -f channels/org1-anchors.tx'
docker exec -it cli-peer0-org2 bash -c 'peer channel update -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem -c mainchannel -f channels/org2-anchors.tx'
```

Now we are going to install the chaincode - NOTE: Make sure you go mod vendor in each chaincode folder... might need to remove the go.sum depending
```bash
docker exec -it cli-peer0-org1 bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1'
docker exec -it cli-peer1-org1 bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1'
docker exec -it cli-peer0-org2 bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1'
docker exec -it cli-peer1-org2 bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1'

docker exec -it cli-peer0-org1 bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt'
docker exec -it cli-peer1-org1 bash -c 'peer lifecycle chaincode install resource_types.tar.gz'
docker exec -it cli-peer0-org2 bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt'
docker exec -it cli-peer1-org2 bash -c 'peer lifecycle chaincode install resource_types.tar.gz'

docker exec -it cli-peer0-org1 bash -c 'peer lifecycle chaincode approveformyorg -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem --channelID mainchannel --name resource_types --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'
docker exec -it cli-peer0-org2 bash -c 'peer lifecycle chaincode approveformyorg -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem --channelID mainchannel --name resource_types --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'

docker exec -it cli-peer0-org1 bash -c 'peer lifecycle chaincode commit -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem --channelID mainchannel --name resource_types --version 1.0 --sequence 1'
```

Lets go ahead and test this chaincode
```bash
docker exec -it cli-peer0-org1 bash -c 'peer chaincode invoke -C mainchannel -n resource_types -c '\''{"Args":["Create","1","Raw Resource"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
sleep 5
docker exec -it cli-peer0-org1 bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
```

Lets try the other chaincode
```bash
docker exec -it cli-peer0-org1 bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1'
docker exec -it cli-peer1-org1 bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1'
docker exec -it cli-peer0-org2 bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1'
docker exec -it cli-peer1-org2 bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1'


docker exec -it cli-peer0-org1 bash -c 'peer lifecycle chaincode install resources.tar.gz &> pkg.txt'
docker exec -it cli-peer1-org1 bash -c 'peer lifecycle chaincode install resources.tar.gz'
docker exec -it cli-peer0-org2 bash -c 'peer lifecycle chaincode install resources.tar.gz &> pkg.txt'
docker exec -it cli-peer1-org2 bash -c 'peer lifecycle chaincode install resources.tar.gz'

docker exec -it cli-peer0-org1 bash -c 'peer lifecycle chaincode approveformyorg -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem --channelID mainchannel --name resources --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'
docker exec -it cli-peer0-org2 bash -c 'peer lifecycle chaincode approveformyorg -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem --channelID mainchannel --name resources --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'

docker exec -it cli-peer0-org1 bash -c 'peer lifecycle chaincode commit -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem --channelID mainchannel --name resources --version 1.0 --sequence 1'

sleep 5

docker exec -it cli-peer0-org1 bash -c 'peer chaincode invoke -C mainchannel -n resources -c '\''{"Args":["Create","1","Iron Ore","1"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-org1 bash -c 'peer chaincode invoke -C mainchannel -n resources -c '\''{"Args":["Create","2","Copper Ore","1"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
sleep 5
docker exec -it cli-peer0-org1 bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer1-org1 bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-org2 bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer1-org2 bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
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
[{"id":"1","name":"Iron Ore","resource_type_id":"1","active":true},{"id":"2","name":"Copper Ore","resource_type_id":"1","active":true}]%                                                                                           ☁  k8s-hyperledger-fabric-2.2 [master] ⚡  curl localhost:4001/resources
[{"id":"1","name":"Iron Ore","resource_type_id":"1","active":true},{"id":"2","name":"Copper Ore","resource_type_id":"1","active":true}]%                                                                                           ☁  k8s-hyperledger-fabric-2.2 [master] ⚡  
```

Lets run the front end and test it
```bash
cd frontend
npm run serve
```

Everything should work!

## Kubernetes - Minikube (Local)

Okay, now that we've successfully ran the network locally, let's do this on a local kubernetes installation.

## Commands
Clean all docker things
```bash
docker system prune 
docker container prune
docker volume prune 
docker rmi $(docker images -q) --force
```

## Common issues

- [Having docker registry issues?](https://github.com/moby/moby/issues/22635)
- There is an issue with CA 1.4.9... use 1.4.7
- [Really good Golang tutorial](https://chainhero.io/2018/06/tutorial-build-blockchain-app-v1-1-0/)
- Issue with watchers
```bash
Error: ENOSPC: System limit for number of file watchers reached, watch '/home/nick/Projects/k8s-hyperledger-fabric-2.2/frontend/public/index.html'
```
Run this command
```bash
echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf && sudo sysctl -p
```
