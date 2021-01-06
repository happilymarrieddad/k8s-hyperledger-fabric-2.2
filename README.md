K8s Hyperledger 2.2 Network
======================================

## Important things to read
- [Orderering Service](https://hyperledger-fabric.readthedocs.io/en/release-2.2/orderer/ordering_service.html)
- [Peers](https://hyperledger-fabric.readthedocs.io/en/release-2.2/peers/peers.html)

## Getting Started
- [My bash](https://github.com/ohmyzsh/ohmyzsh)
- [Install Hyperledger Deps](https://hyperledger-fabric.readthedocs.io/en/release-2.2/install.html)
- curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.2.1 1.4.7
- [Create aws account](aws.amazon.com)
- [Create docker hub account](https://hub.docker.com/)
- [Install Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Install Minikube](https://minikube.sigs.k8s.io/docs/start/)
- [Install Docker-compose](https://docs.docker.com/engine/install/ubuntu/)

You want to copy over all the bin files into a bin directory where you can easily access them. I usually just like to make a bin folder in my home directory and set the path to include that
```bash
mkdir -p ~/bin
cp fabric-samples/bin/* ~/bin
```
=kj
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
configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channels/ibm-anchors.tx -channelID mainchannel -asOrg ibm
configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channels/oracle-anchors.tx -channelID mainchannel -asOrg oracle
```

Now it's time to start the network
```bash
docker-compose -f network/docker/docker-compose.yaml up
```

Lets setup the artifacts
```bash
docker exec -it cli-peer0-ibm bash -c 'peer channel create -c mainchannel -f ./channels/mainchannel.tx -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer1-ibm bash -c 'peer channel create -c mainchannel -f ./channels/mainchannel.tx -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-oracle bash -c 'peer channel create -c mainchannel -f ./channels/mainchannel.tx -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer1-oracle bash -c 'peer channel create -c mainchannel -f ./channels/mainchannel.tx -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'

docker exec -it cli-peer0-ibm bash -c 'cp mainchannel.block ./channels/'
docker exec -it cli-peer0-ibm bash -c 'peer channel join -b channels/mainchannel.block'
docker exec -it cli-peer1-ibm bash -c 'peer channel join -b channels/mainchannel.block'
docker exec -it cli-peer0-oracle bash -c 'peer channel join -b channels/mainchannel.block'
docker exec -it cli-peer1-oracle bash -c 'peer channel join -b channels/mainchannel.block'

sleep 5

docker exec -it cli-peer0-ibm bash -c 'peer channel update -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem -c mainchannel -f channels/ibm-anchors.tx'
docker exec -it cli-peer0-oracle bash -c 'peer channel update -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem -c mainchannel -f channels/oracle-anchors.tx'
```

Now we are going to install the chaincode - NOTE: Make sure you go mod vendor in each chaincode folder... might need to remove the go.sum depending
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
```bash
docker exec -it cli-peer0-ibm bash -c 'peer chaincode invoke -C mainchannel -n resource_types -c '\''{"Args":["Create", "1","Parts"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
sleep 5
docker exec -it cli-peer0-ibm bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-ibm bash -c 'peer chaincode invoke -C mainchannel -n resource_types -c '\''{"Args":["Update", "1", "Parts 2"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-ibm bash -c 'peer chaincode invoke -C mainchannel -n resource_types -c '\''{"Args":["Update", "1", "Parts"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-ibm bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Transactions","1"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
```

Lets try the other chaincode
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

docker exec -it cli-peer0-ibm bash -c 'peer chaincode invoke -C mainchannel -n resources -c '\''{"Args":["Create","CPUs","1"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-oracle bash -c 'peer chaincode invoke -C mainchannel -n resources -c '\''{"Args":["Create","Database Servers","1"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
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

## Kubernetes - Minikube (Local)
[Install Kubernetes and Minikube](https://kubernetes.io/docs/tasks/tools/)
[If OSX here is virtual box](https://www.virtualbox.org/wiki/Mac%20OS%20X%20build%20instructions)
[Kubernetes book](hhttps://www.amazon.com/Devops-2-3-Toolkit-Viktor-Farcic/dp/1789135508/ref=tmm_pap_swatch_0?_encoding=UTF8&sr=8-2)
[K8s Persistent Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)

Okay, now that we've successfully ran the network locally, let's do this on a local kubernetes installation.
```bash
minikube start
sleep 5
kubectl apply -f network/minikube/storage/pvc.yaml
sleep 10
kubectl apply -f network/minikube/storage/tests
```

Now we have storage and we're going to test it. You can do a kubectl get pods to see what pods are up. Here's how I can connect to my containers. You should split your terminal and connect to both.
```bash
☁  k8s-hyperledger-fabric-2.2 [master] ⚡  kubectl get pods
NAME                        READY   STATUS    RESTARTS   AGE
example1-6858b4f776-5pgls   1/1     Running   0          17s
example1-6858b4f776-q92vv   1/1     Running   0          17s
example2-55fcbb9cbd-drzwn   1/1     Running   0          17s
example2-55fcbb9cbd-sv4c8   1/1     Running   0          17s
☁  k8s-hyperledger-fabric-2.2 [master] ⚡ 
```

We'll use one of these to setup the files for the network
```bash
kubectl exec -it $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//") -- mkdir -p /host/files/scripts
kubectl exec -it $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//") -- mkdir -p /host/files/chaincode

kubectl cp ./scripts $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files
kubectl cp ./network/minikube/configtx.yaml $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files
kubectl cp ./network/minikube/config.yaml $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files
kubectl cp ./chaincode/resources $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files/chaincode
kubectl cp ./chaincode/resource_types $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files/chaincode
kubectl cp ./fabric-samples/bin $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files
```

Let's bash into the container and make sure everything copied over properly
```bash
kubectl exec -it $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//") bash
```

Finally ready to start the ca containers
```bash
kubectl apply -f network/minikube/cas
```

Your containers should be up and running. You can check the logs like so and it should look liek this.
```bash
☁  k8s-hyperledger-fabric-2.2 [master] ⚡  kubectl logs -f orderers-ca-d69cbc664-dzk4f
2020/12/11 04:12:37 [INFO] Created default configuration file at /etc/hyperledger/fabric-ca-server/fabric-ca-server-config.yaml
2020/12/11 04:12:37 [INFO] Starting server in home directory: /etc/hyperledger/fabric-ca-server
...
2020/12/11 04:12:38 [INFO] generating key: &{A:ecdsa S:256}
2020/12/11 04:12:38 [INFO] encoded CSR
2020/12/11 04:12:38 [INFO] signed certificate with serial number 307836600921505839273746385963411812465330101584
2020/12/11 04:12:38 [INFO] Listening on https://0.0.0.0:7054
```

This should generate the crypto-config files necessary for the network. You can check on those files in any of the containers.
```bash
root@example1-6858b4f776-wmlth:/host# cd files
root@example1-6858b4f776-wmlth:/host/files# ls
bin  chaincode	config.yaml  configtx.yaml  crypto-config  scripts
root@example1-6858b4f776-wmlth:/host/files# cd crypto-config/
root@example1-6858b4f776-wmlth:/host/files/crypto-config# ls
ordererOrganizations  peerOrganizations
root@example1-6858b4f776-wmlth:/host/files/crypto-config# cd peerOrganizations/
root@example1-6858b4f776-wmlth:/host/files/crypto-config/peerOrganizations# ls
ibm  oracle
root@example1-6858b4f776-wmlth:/host/files/crypto-config/peerOrganizations# cd ibm/
root@example1-6858b4f776-wmlth:/host/files/crypto-config/peerOrganizations/ibm# ls
msp  peers  users
root@example1-6858b4f776-wmlth:/host/files/crypto-config/peerOrganizations/ibm# cd msp/
root@example1-6858b4f776-wmlth:/host/files/crypto-config/peerOrganizations/ibm/msp# ls
IssuerPublicKey  IssuerRevocationPublicKey  admincerts	cacerts  keystore  signcerts  tlscacerts  user
root@example1-6858b4f776-wmlth:/host/files/crypto-config/peerOrganizations/ibm/msp# cd tlscacerts/
```

Time to generate the artifacts inside one of the containers and in the files folder - NOTE: if you are on OSX you might have to load the proper libs `curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.2.1 1.4.7`
```bash
kubectl exec -it $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//") bash
...
cd /host/files

rm -rf orderer channels
mkdir -p orderer channels
bin/configtxgen -profile OrdererGenesis -channelID syschannel -outputBlock ./orderer/genesis.block
bin/configtxgen -profile MainChannel -outputCreateChannelTx ./channels/mainchannel.tx -channelID mainchannel
bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channels/ibm-anchors.tx -channelID mainchannel -asOrg ibm
bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channels/oracle-anchors.tx -channelID mainchannel -asOrg oracle
```

Let's try to start up the orderers
```bash
kubectl apply -f network/minikube/orderers
```

Go ahead and check the logs and see that the orderers have selected a leader like so
```bash
 1 became follower at term 2 channel=syschannel node=1
2020-12-11 05:20:15.616 UTC [orderer.consensus.etcdraft] Step -> INFO 029 1 [logterm: 1, index: 3, vote: 0] cast MsgVote for 2 [logterm: 1, index: 3] at term 2 channel=syschannel node=1
2020-12-11 05:20:15.634 UTC [orderer.consensus.etcdraft] run -> INFO 02a raft.node: 1 elected leader 2 at term 2 channel=syschannel node=1
2020-12-11 05:20:15.639 UTC [orderer.consensus.etcdraft] run -> INFO 02b Raft leader changed: 0 -> 2 channel=syschannel node=1
```

We should be able to start the peers now
```bash
kubectl apply -f network/minikube/orgs/ibm/couchdb 
kubectl apply -f network/minikube/orgs/oracle/couchdb

kubectl apply -f network/minikube/orgs/ibm/
kubectl apply -f network/minikube/orgs/oracle/

kubectl apply -f network/minikube/orgs/ibm/cli
kubectl apply -f network/minikube/orgs/oracle/cli
```


Time to actually test the network
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel create -c mainchannel -f ./channels/mainchannel.tx -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'cp mainchannel.block ./channels/'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel join -b channels/mainchannel.block'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel join -b channels/mainchannel.block'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel join -b channels/mainchannel.block'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel join -b channels/mainchannel.block'

sleep 5

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel update -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem -c mainchannel -f channels/ibm-anchors.tx'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel update -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem -c mainchannel -f channels/oracle-anchors.tx'
```

Now we are going to install the chaincode - NOTE: Make sure you go mod vendor in each chaincode folder... might need to remove the go.sum depending
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --channelID mainchannel --collections-config /opt/gopath/src/resource_types/collections-config.json --name resource_types --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode commit -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --channelID mainchannel --collections-config /opt/gopath/src/resource_types/collections-config.json --name resource_types --version 1.0 --sequence 1'
```

Lets go ahead and test this chaincode
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode invoke -C mainchannel -n resource_types -c '\''{"Args":["Create", "1","Parts"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
sleep 5
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
```

Lets try the other chaincode
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1'


kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz &> pkg.txt'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz &> pkg.txt'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --channelID mainchannel --collections-config /opt/gopath/src/resources/collections-config.json --name resources --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --channelID mainchannel --collections-config /opt/gopath/src/resources/collections-config.json --name resources --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode commit -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --channelID mainchannel --collections-config /opt/gopath/src/resources/collections-config.json --name resources --version 1.0 --sequence 1'

sleep 5

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode invoke -C mainchannel -n resources -c '\''{"Args":["Create","CPUs","1"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode invoke -C mainchannel -n resources -c '\''{"Args":["Create","Database Servers","1"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode invoke -C mainchannel -n resources -c '\''{"Args":["Create","Mainframe Boards","1"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
sleep 5
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
```

Start the API
```bash
kubectl apply -f network/minikube/backend
```

Get the address for the nodeport
```bash
minikube service api-service-nodeport --url
```

## Kubernetes - Minikube (Production)

- [AWS Install - depends on your OS](https://linuxhint.com/install_aws_cli_ubuntu/)
- Now, need to configure AWS. Need to get API keys for aws.
- login to AWS and create a user under IAM. You need to click the box for API access
- you will get an Access Key ID and a Secret Access Key
- run the aws configure command and enter those credentials. I also use us-west-1 and json for the other two options

```bash
☁  k8s-hyperledger-fabric-2.2 [master] ⚡  aws configure
AWS Access Key ID [None]: <ID>
AWS Secret Access Key [None]: <secret>
Default region name [None]: us-west-1
Default output format [None]: json
☁  k8s-hyperledger-fabric-2.2 [master] ⚡
```

[Installing KOPS](https://kops.sigs.k8s.io/getting_started/install/)
[READ THIS!!!](https://github.com/kubernetes/kops/blob/master/docs/getting_started/aws.md#setup-iam-user)

- You need to setup an s3 bucket in the same zone as your AWS config and set the appropriate kops bash config
```bash
export KOPS_CLUSTER_NAME=hyperledger.k8s.local
export KOPS_STATE_STORE=s3://<name of your bucket>
export AWS_ACCESS_KEY_ID=$(aws configure get aws_access_key_id)
export AWS_SECRET_ACCESS_KEY=$(aws configure get aws_secret_access_key)
```

Time to create the cluster
```bash
ssh-keygen
(No passphrase)

kops create cluster \
    --zones us-west-1b,us-west-1c \
    --node-count 3 \
    --master-zones us-west-1b,us-west-1c \
    --master-count 3 \
    --authorization AlwaysAllow --yes \
    --master-volume-size 40 \
    --node-volume-size 20
```

To delete the cluster
```bash
kops delete cluster hyperledger.k8s.local --yes
```

You can edit the nodes
```bash
kops edit ig nodes
```

Or the masters
```bash
kops edit ig masters
```

After the cluster is created, need to add some secrets to the network
```bash
kubectl create secret generic couchdb --from-literal username=nick --from-literal password=1234

kubectl create secret docker-registry regcred \
    --docker-server=https://index.docker.io/v1/ \
    --docker-username=happilymarrieddadudemy \
    --docker-password=<password> \
    --docker-email=<email>
```

Adding nginx to our network
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/kops/master/addons/ingress-nginx/v1.6.0.yaml

kubectl create secret tls udemy-hyperledger.com --key ~/udemy-hyperledger.com/privkey.pem --cert ~/udemy-hyperledger.com/cert.pem

kubectl apply -f network/production/ingress-nginx.yaml
```

Now, lets add the NFS file system. Go ahead and login to your AWS account and go to EFS. Create a NFS file system in the same REGION as the cluster and make sure to SET THE VPC the same as the network. VERY IMPORTANT!!!! Also, create mount points and set them to include ALL of the permissions for the network (should be for of them). Now, we can create the storage by using the PV and PVC yaml files. We're going to use multiple PVC's just to show how to do that.
```bash
kubectl apply -f network/production/storage/pv.yaml 
kubectl apply -f network/production/storage/pvc.yaml
kubectl apply -f network/minikube/storage/tests 
```

Bash into the containers, create a file and make sure it's available in the other containers. Make sure you do it in the /host folder because that's the folder that's mounted.
```bash
☁  k8s-hyperledger-fabric-2.2 [master] ⚡  kubectl get pods
NAME                        READY   STATUS              RESTARTS   AGE
example1-657d584cc7-qdgzx   0/1     ContainerCreating   0          24s
example2-fdcd6dfc5-v7p28    0/1     ContainerCreating   0          25s
example2-fdcd6dfc5-x8snv    0/1     ContainerCreating   0          25s
☁  k8s-hyperledger-fabric-2.2 [master] ⚡  kubectl exec -it example1-657d584cc7-qdgzx bash
kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl exec [POD] -- [COMMAND] instead.
root@example1-657d584cc7-qdgzx:/# cd host/file
bash: cd: host/file: No such file or directory
root@example1-657d584cc7-qdgzx:/# cd host/files
bash: cd: host/files: No such file or directory
root@example1-657d584cc7-qdgzx:/# cd host
root@example1-657d584cc7-qdgzx:/host# mkdir files
root@example1-657d584cc7-qdgzx:/host# ls
files
root@example1-657d584cc7-qdgzx:/host# echo "Hello World" >> test.txt
root@example1-657d584cc7-qdgzx:/host# 
```

Other terminal
```bash
☁  k8s-hyperledger-fabric-2.2 [master] ⚡  kubectl get pods
NAME                        READY   STATUS    RESTARTS   AGE
example1-657d584cc7-qdgzx   1/1     Running   0          40s
example2-fdcd6dfc5-v7p28    1/1     Running   0          41s
example2-fdcd6dfc5-x8snv    1/1     Running   0          41s
☁  k8s-hyperledger-fabric-2.2 [master] ⚡  kubectl exec -it example2-fdcd6dfc5-v7p28 -- bash
root@example2-fdcd6dfc5-v7p28:/# cd host
root@example2-fdcd6dfc5-v7p28:/host# cat test.txt 
Hello World
root@example2-fdcd6dfc5-v7p28:/host# 
```

Okay, now we can just create the network the same as we would in minikube.
```bash
☁  k8s-hyperledger-fabric-2.2 [master] ⚡  kubectl get pods
NAME                        READY   STATUS    RESTARTS   AGE
example1-6858b4f776-5pgls   1/1     Running   0          17s
example1-6858b4f776-q92vv   1/1     Running   0          17s
example2-55fcbb9cbd-drzwn   1/1     Running   0          17s
example2-55fcbb9cbd-sv4c8   1/1     Running   0          17s
☁  k8s-hyperledger-fabric-2.2 [master] ⚡ 
```

We'll use one of these to setup the files for the network
```bash
kubectl exec -it $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//") -- mkdir -p /host/files/scripts
kubectl exec -it $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//") -- mkdir -p /host/files/chaincode

kubectl cp ./scripts $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files
kubectl cp ./network/production/configtx.yaml $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files
kubectl cp ./network/production/config.yaml $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files
kubectl cp ./chaincode/resources $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files/chaincode
kubectl cp ./chaincode/resource_types $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files/chaincode
kubectl cp ~/bin $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files
```


Let's bash into the container and make sure everything copied over properly
```bash
kubectl exec -it $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//") bash
```

Finally ready to start the ca containers
```bash
kubectl apply -f network/minikube/cas
```

Your containers should be up and running. You can check the logs like so and it should look liek this.
```bash
☁  k8s-hyperledger-fabric-2.2 [master] ⚡  kubectl logs -f orderers-ca-d69cbc664-dzk4f
2020/12/11 04:12:37 [INFO] Created default configuration file at /etc/hyperledger/fabric-ca-server/fabric-ca-server-config.yaml
2020/12/11 04:12:37 [INFO] Starting server in home directory: /etc/hyperledger/fabric-ca-server
...
2020/12/11 04:12:38 [INFO] generating key: &{A:ecdsa S:256}
2020/12/11 04:12:38 [INFO] encoded CSR
2020/12/11 04:12:38 [INFO] signed certificate with serial number 307836600921505839273746385963411812465330101584
2020/12/11 04:12:38 [INFO] Listening on https://0.0.0.0:7054
```

This should generate the crypto-config files necessary for the network. You can check on those files in any of the containers.
```bash
root@example1-6858b4f776-wmlth:/host# cd files
root@example1-6858b4f776-wmlth:/host/files# ls
bin  chaincode	config.yaml  configtx.yaml  crypto-config  scripts
root@example1-6858b4f776-wmlth:/host/files# cd crypto-config/
root@example1-6858b4f776-wmlth:/host/files/crypto-config# ls
ordererOrganizations  peerOrganizations
root@example1-6858b4f776-wmlth:/host/files/crypto-config# cd peerOrganizations/
root@example1-6858b4f776-wmlth:/host/files/crypto-config/peerOrganizations# ls
ibm  oracle
root@example1-6858b4f776-wmlth:/host/files/crypto-config/peerOrganizations# cd ibm/
root@example1-6858b4f776-wmlth:/host/files/crypto-config/peerOrganizations/ibm# ls
msp  peers  users
root@example1-6858b4f776-wmlth:/host/files/crypto-config/peerOrganizations/ibm# cd msp/
root@example1-6858b4f776-wmlth:/host/files/crypto-config/peerOrganizations/ibm/msp# ls
IssuerPublicKey  IssuerRevocationPublicKey  admincerts	cacerts  keystore  signcerts  tlscacerts  user
root@example1-6858b4f776-wmlth:/host/files/crypto-config/peerOrganizations/ibm/msp# cd tlscacerts/
```

Time to generate the artifacts inside one of the containers and in the files folder
```bash
rm -rf orderer channels
mkdir -p orderer channels
bin/configtxgen -profile OrdererGenesis -channelID syschannel -outputBlock ./orderer/genesis.block
bin/configtxgen -profile MainChannel -outputCreateChannelTx ./channels/mainchannel.tx -channelID mainchannel
bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channels/ibm-anchors.tx -channelID mainchannel -asOrg ibm
bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channels/oracle-anchors.tx -channelID mainchannel -asOrg oracle
```

Let's try to start up the orderers
```bash
kubectl apply -f network/minikube/orderers
```

Go ahead and check the logs and see that the orderers have selected a leader like so
```bash
 1 became follower at term 2 channel=syschannel node=1
2020-12-11 05:20:15.616 UTC [orderer.consensus.etcdraft] Step -> INFO 029 1 [logterm: 1, index: 3, vote: 0] cast MsgVote for 2 [logterm: 1, index: 3] at term 2 channel=syschannel node=1
2020-12-11 05:20:15.634 UTC [orderer.consensus.etcdraft] run -> INFO 02a raft.node: 1 elected leader 2 at term 2 channel=syschannel node=1
2020-12-11 05:20:15.639 UTC [orderer.consensus.etcdraft] run -> INFO 02b Raft leader changed: 0 -> 2 channel=syschannel node=1
```

We should be able to start the peers now
```bash
kubectl apply -f network/minikube/orgs/ibm/couchdb 
kubectl apply -f network/minikube/orgs/oracle/couchdb

kubectl apply -f network/minikube/orgs/ibm/
kubectl apply -f network/minikube/orgs/oracle/

kubectl apply -f network/minikube/orgs/ibm/cli
kubectl apply -f network/minikube/orgs/oracle/cli
```


Time to actually test the network
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel create -c mainchannel -f ./channels/mainchannel.tx -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'cp mainchannel.block ./channels/'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel join -b channels/mainchannel.block'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel join -b channels/mainchannel.block'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel join -b channels/mainchannel.block'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel join -b channels/mainchannel.block'

sleep 5

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel update -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem -c mainchannel -f channels/ibm-anchors.tx'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel update -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem -c mainchannel -f channels/oracle-anchors.tx'
```

Now we are going to install the chaincode - NOTE: Make sure you go mod vendor in each chaincode folder... might need to remove the go.sum depending
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode commit -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 1.0 --sequence 1'
```

Lets go ahead and test this chaincode
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode invoke -C mainchannel -n resource_types -c '\''{"Args":["Create", "1","Parts"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
sleep 5
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
```

Lets try the other chaincode
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1'


kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz &> pkg.txt'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz &> pkg.txt'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resources/collections-config.json --channelID mainchannel --name resources --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resources/collections-config.json --channelID mainchannel --name resources --version 1.0 --sequence 1 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode commit -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resources/collections-config.json --channelID mainchannel --name resources --version 1.0 --sequence 1'

sleep 5

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode invoke -C mainchannel -n resources -c '\''{"Args":["Create","CPUs","1"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode invoke -C mainchannel -n resources -c '\''{"Args":["Create","Database Servers","1"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode invoke -C mainchannel -n resources -c '\''{"Args":["Create","Mainframes","1"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
sleep 5
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
```

Startup the api and the web app
```bash
kubectl apply -f network/production/backend/
kubectl apply -f network/production/frontend
```


## Some other ways to run the network

Of if you want to use the auto-join yaml: NOTE - this is experimental and it may not work for you
- uses some unique scripts to auto join the peers and install chaincodes
```bash
docker-compose -f network/docker/docker-compose-auto-join.yaml up
```


Lets go ahead and test this chaincode
```bash
docker exec -it cli-peer0-ibm bash -c 'peer chaincode invoke -C mainchannel -n resource_types -c '\''{"Args":["SetPrivateData", "1", "IBM Private Name"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-oracle bash -c 'peer chaincode invoke -C mainchannel -n resource_types -c '\''{"Args":["SetPrivateData", "1", "ORACLE Private Name"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'

sleep 5

docker exec -it cli-peer0-ibm bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-oracle bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'

docker exec -it cli-peer0-oracle bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Transactions","1"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
docker exec -it cli-peer0-ibm bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Transactions","1"]}'\'' -o orderer0:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-7054.pem'
```


## Commands
Clean all docker things
```bash
docker system prune 
docker container prune
docker volume prune 
docker rmi $(docker images -q) --force
```

Kill minikube env
```bash
minikube delete
```

You can purge everything with
```bash
minikube delete --purge
```

You can find and replace file names like this
```bash
sudo find . -depth -name '*org1*' -execdir bash -c 'mv -i "$1" "${1//org1/ibm}"' bash {} \;
sudo find . -depth -name '*org2*' -execdir bash -c 'mv -i "$1" "${1//org2/oracle}"' bash {} \;
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


## Updating chaincode

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

## Encryption Methods
+3

Hello World
KhoorZruog
HelloWorld


Private Key
Public Key

Hello World
fdsgdhfdhjyt

## ADDING OTHER ORGS

Create folder to store data in
```bash
export NEW_ORG_NAME=hp
export FOLDER_PATH=configs/minikube/${NEW_ORG_NAME}
mkdir -p $FOLDER_PATH/cas
mkdir -p $FOLDER_PATH/cli
mkdir -p $FOLDER_PATH/couchdb
mkdir -p $FOLDER_PATH/peers
```

Create necessary files
```bash
cat <<EOT > ${FOLDER_PATH}/configtx.yaml
Organizations:
    - &orderer
        Name: orderer
        ID: orderer
        MSPDir: /host/files/crypto-config/ordererOrganizations/orderer/msp
        Policies:
            Readers:
                Type: Signature
                Rule: "OR('orderer.member')"
            Writers:
                Type: Signature
                Rule: "OR('orderer.member')"
            Admins:
                Type: Signature
                Rule: "OR('orderer.admin')"
    - &ibm
        Name: ibm
        ID: ibm
        MSPDir: /host/files/crypto-config/peerOrganizations/ibm/msp
        AnchorPeers:
            - Host: peer0-ibm-service
              Port: 7051
        Policies:
            Readers:
                Type: Signature
                Rule: "OR('ibm.member')"
            Writers:
                Type: Signature
                Rule: "OR('ibm.member')"
            Admins:
                Type: Signature
                Rule: "OR('ibm.admin')"
            Endorsement:
                Type: Signature
                Rule: "OR('ibm.member')"
    - &oracle
        Name: oracle
        ID: oracle
        MSPDir: /host/files/crypto-config/peerOrganizations/oracle/msp
        AnchorPeers:
            - Host: peer0-oracle-service
              Port: 7051
        Policies:
            Readers:
                Type: Signature
                Rule: "OR('oracle.member')"
            Writers:
                Type: Signature
                Rule: "OR('oracle.member')"
            Admins:
                Type: Signature
                Rule: "OR('oracle.admin')"
            Endorsement:
                Type: Signature
                Rule: "OR('oracle.member')"
    - &${NEW_ORG_NAME}
        Name: ${NEW_ORG_NAME}
        ID: ${NEW_ORG_NAME}
        MSPDir: /host/files/crypto-config/peerOrganizations/${NEW_ORG_NAME}/msp
        AnchorPeers:
            - Host: peer0-${NEW_ORG_NAME}-service
              Port: 7051
        Policies:
            Readers:
                Type: Signature
                Rule: "OR('${NEW_ORG_NAME}.member')"
            Writers:
                Type: Signature
                Rule: "OR('${NEW_ORG_NAME}.member')"
            Admins:
                Type: Signature
                Rule: "OR('${NEW_ORG_NAME}.admin')"
            Endorsement:
                Type: Signature
                Rule: "OR('${NEW_ORG_NAME}.member')"

Capabilities:
    Global: &ChannelCapabilities
        V2_0: true
    Orderer: &OrdererCapabilities
        V2_0: true
    Application: &ApplicationCapabilities
        V2_0: true

Application: &ApplicationDefaults
    Organizations:
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: "ANY Readers"
        Writers:
            Type: ImplicitMeta
            Rule: "ANY Writers"
        Admins:
            Type: ImplicitMeta
            Rule: "MAJORITY Admins"
        LifecycleEndorsement:
            Type: ImplicitMeta
            Rule: "ANY Endorsement"
        Endorsement:
            Type: ImplicitMeta
            Rule: "ANY Endorsement"
    Capabilities:
        <<: *ApplicationCapabilities

Orderer: &OrdererDefaults
    OrdererType: etcdraft
    EtcdRaft:
        Consenters:
            - Host: orderer0-service
              Port: 7050
              ClientTLSCert: /host/files/crypto-config/ordererOrganizations/orderer/orderers/orderer0/tls/server.crt
              ServerTLSCert: /host/files/crypto-config/ordererOrganizations/orderer/orderers/orderer0/tls/server.crt
            - Host: orderer1-service
              Port: 7050
              ClientTLSCert: /host/files/crypto-config/ordererOrganizations/orderer/orderers/orderer1/tls/server.crt
              ServerTLSCert: /host/files/crypto-config/ordererOrganizations/orderer/orderers/orderer1/tls/server.crt
            - Host: orderer2-service
              Port: 7050
              ClientTLSCert: /host/files/crypto-config/ordererOrganizations/orderer/orderers/orderer2/tls/server.crt
              ServerTLSCert: /host/files/crypto-config/ordererOrganizations/orderer/orderers/orderer2/tls/server.crt
    Addresses:
        - orderer0:7050-service
        - orderer1:7050-service
        - orderer2:7050-service
    BatchTimeout: 2s
    BatchSize:
        MaxMessageCount: 10
        AbsoluteMaxBytes: 99 MB
        PreferredMaxBytes: 512 KB
    Kafka:
        Brokers:
            - 127.0.0.1:9092
    Organizations:
        - *orderer
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: "ANY Readers"
        Writers:
            Type: ImplicitMeta
            Rule: "ANY Writers"
        Admins:
            Type: ImplicitMeta
            Rule: "MAJORITY Admins"
        BlockValidation:
            Type: ImplicitMeta
            Rule: "ANY Writers"

Channel: &ChannelDefaults
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: "ANY Readers"
        Writers:
            Type: ImplicitMeta
            Rule: "ANY Writers"
        Admins:
            Type: ImplicitMeta
            Rule: "MAJORITY Admins"
    Capabilities:
        <<: *ChannelCapabilities

Profiles:
    MainChannel:
        Consortium: MAIN
        <<: *ChannelDefaults
        Application:
            <<: *ApplicationDefaults
            Organizations:
                - *ibm
                - *oracle
                - *${NEW_ORG_NAME}
            Capabilities:
                <<: *ApplicationCapabilities

EOT

cat <<EOT > ${FOLDER_PATH}/cas/${NEW_ORG_NAME}-ca-client-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${NEW_ORG_NAME}-ca-client
  labels: {
    component: ${NEW_ORG_NAME}-ca-client,
    type: ca
  }
spec:
  replicas: 1
  selector:
    matchLabels:
      component: ${NEW_ORG_NAME}-ca-client
      type: ca
  template:
    metadata:
      labels:
        component: ${NEW_ORG_NAME}-ca-client
        type: ca
    spec:
      volumes:
        - name: my-pv-storage
          persistentVolumeClaim:
            claimName: my-pv-claim
      containers:
        - name: ${NEW_ORG_NAME}-ca-client
          image: hyperledger/fabric-ca:1.4.7
          command: ["bash"]
          args: ["/scripts/start-org-client.sh"]
          env:
            - name: FABRIC_CA_SERVER_HOME
              value: /etc/hyperledger/fabric-ca-client
            - name: ORG_NAME
              value: ${NEW_ORG_NAME}
            - name: CA_SCHEME
              value: https
            - name: CA_URL
              value: "${NEW_ORG_NAME}-ca-service:7054"
            - name: CA_USERNAME
              value: admin
            - name: CA_PASSWORD
              value: adminpw
            - name: CA_CERT_PATH
              value: /etc/hyperledger/fabric-ca-server/tls-cert.pem
          volumeMounts:
            - mountPath: /scripts
              name: my-pv-storage
              subPath: files/scripts
            - mountPath: /state
              name: my-pv-storage
              subPath: state
            - mountPath: /files
              name: my-pv-storage
              subPath: files
            - mountPath: /etc/hyperledger/fabric-ca-server
              name: my-pv-storage
              subPath: state/${NEW_ORG_NAME}-ca
            - mountPath: /etc/hyperledger/fabric-ca-client
              name: my-pv-storage
              subPath: state/${NEW_ORG_NAME}-ca-client
            - mountPath: /etc/hyperledger/fabric-ca/crypto-config
              name: my-pv-storage
              subPath: files/crypto-config
EOT

cat <<EOT > ${FOLDER_PATH}/cas/${NEW_ORG_NAME}-ca-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${NEW_ORG_NAME}-ca
  labels: {
    component: ${NEW_ORG_NAME},
    type: ca
  }
spec:
  replicas: 1
  selector:
    matchLabels:
      component: ${NEW_ORG_NAME}
      type: ca
  template:
    metadata:
      labels:
        component: ${NEW_ORG_NAME}
        type: ca
    spec:
      volumes:
        - name: my-pv-storage
          persistentVolumeClaim:
            claimName: my-pv-claim
      containers:
        - name: ${NEW_ORG_NAME}-ca
          image: hyperledger/fabric-ca:1.4.7
          command: ["sh"]
          args: ["/scripts/start-root-ca.sh"]
          ports:
            - containerPort: 7054
          env:
            - name: FABRIC_CA_HOME
              value: /etc/hyperledger/fabric-ca-server
            - name: USERNAME
              value: admin
            - name: PASSWORD
              value: adminpw
            - name: CSR_HOSTS
              value: ${NEW_ORG_NAME}-ca
          volumeMounts:
            - mountPath: /scripts
              name: my-pv-storage
              subPath: files/scripts
            - mountPath: /etc/hyperledger/fabric-ca-server
              name: my-pv-storage
              subPath: state/${NEW_ORG_NAME}-ca

EOT

cat <<EOT > ${FOLDER_PATH}/cas/${NEW_ORG_NAME}-ca-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: ${NEW_ORG_NAME}-ca-service
  labels: {
    component: ${NEW_ORG_NAME},
    type: ca
  }
spec:
  type: ClusterIP
  selector:
    component: ${NEW_ORG_NAME}
    type: ca
  ports:
    - port: 7054
      targetPort: 7054

EOT
```

Start hp ca and create certs
```bash
kubectl apply -f ${FOLDER_PATH}/cas
sleep 20
```

NOTE: THIS IS NOT NECESSARY (DO NOT DO THIS!!!)
Get the config from the configtx (NOTE: There isn't really a way to pass in ENV into a kubectl command easily. Just do it manually) (In this example, I just hardcoded the orgName in the configtxgen command for "hp")
TODO: Create a script for this with the org name and then copy it and run it in the container
```bash
kubectl cp ${FOLDER_PATH}/configtx.yaml $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files
kubectl exec -it $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//") bash
...
cd /host/files

bin/configtxgen -printOrg hp > ./channels/hp.json
```

Use an existing org to get the current config
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'apk update'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'apk add jq'
sleep 1
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel fetch config config_block.pb -c mainchannel -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
sleep 1
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config > config.json'

sleep 1
cat <<EOT > scripts/org_modified_config.sh
#!/bin/bash

jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"${NEW_ORG_NAME}":.[1]}}}}}' config.json ./channels/${NEW_ORG_NAME}.json > modified_config.json
EOT
chmod +x scripts/org_modified_config.sh
kubectl cp ./scripts $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//"):/opt/gopath/src/github.com/hyperledger/fabric/peer

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c './scripts/org_modified_config.sh'
sleep 1
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'configtxlator proto_encode --input config.json --type common.Config --output config.pb'
sleep 1
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb'
sleep 1
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c '\
	configtxlator compute_update \
	--channel_id mainchannel \
	--original config.pb \
	--updated modified_config.pb \
	--output org_update.pb'
sleep 1
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c '\
	configtxlator proto_decode \
	--input org_update.pb \
	--type common.ConfigUpdate | jq . > org_update.json \
	'
sleep 1
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c './scripts/create-org-envelope.sh'
sleep 1
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c '\
	configtxlator proto_encode --input org_update_in_envelope.json --type common.Envelope --output org_update_in_envelope.pb \
	'
sleep 1
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c '\
	peer channel signconfigtx -f org_update_in_envelope.pb \
	'
sleep 1
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c '\
	cp org_update_in_envelope.pb channels/org_update_in_envelope.pb \
	'
sleep 1
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c '\
	peer channel update -f channels/org_update_in_envelope.pb -c mainchannel -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem \
	'
```

Okay, so now we've updated the existing config with the new org. Time to startup the new org
```bash
cat <<EOT > $FOLDER_PATH/couchdb/services.yaml
apiVersion: v1
kind: Service
metadata:
  name: peer0-${NEW_ORG_NAME}-couchdb
  labels: {
    component: peer0,
    type: couchdb,
    org: ${NEW_ORG_NAME}
  }
spec:
  type: ClusterIP
  selector:
    component: peer0
    type: couchdb
    org: ${NEW_ORG_NAME}
  ports:
    - port: 5984
      targetPort: 5984
---
apiVersion: v1
kind: Service
metadata:
  name: peer1-${NEW_ORG_NAME}-couchdb
  labels: {
    component: peer1,
    type: couchdb,
    org: ${NEW_ORG_NAME}
  }
spec:
  type: ClusterIP
  selector:
    component: peer1
    type: couchdb
    org: ${NEW_ORG_NAME}
  ports:
    - port: 5984
      targetPort: 5984

EOT

cat <<EOT > $FOLDER_PATH/couchdb/peer0-couchdb-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: peer0-${NEW_ORG_NAME}-couchdb-deployment
  labels: {
    component: peer0,
    type: couchdb,
    org: ${NEW_ORG_NAME}
  }
spec:
  replicas: 1
  selector:
    matchLabels:
      component: peer0
      type: couchdb
      org: ${NEW_ORG_NAME}
  template:
    metadata:
      labels:
        component: peer0
        type: couchdb
        org: ${NEW_ORG_NAME}
    spec:
      volumes:
        - name: my-pv-storage
          persistentVolumeClaim:
            claimName: my-pv-claim
      containers:
        - name: peer0-${NEW_ORG_NAME}-couchdb
          image: couchdb:3.1.1
          env:
            - name: COUCHDB_USER
              value: nick
            - name: COUCHDB_PASSWORD
              value: "1234"
          volumeMounts:
            - mountPath: /opt/couchdb/data
              name: my-pv-storage
              subPath: state/${NEW_ORG_NAME}/peers/peer0-couchdb

EOT

cat <<EOT > $FOLDER_PATH/couchdb/peer1-couchdb-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: peer1-${NEW_ORG_NAME}-couchdb-deployment
  labels: {
    component: peer1,
    type: couchdb,
    org: ${NEW_ORG_NAME}
  }
spec:
  replicas: 1
  selector:
    matchLabels:
      component: peer1
      type: couchdb
      org: ${NEW_ORG_NAME}
  template:
    metadata:
      labels:
        component: peer1
        type: couchdb
        org: ${NEW_ORG_NAME}
    spec:
      volumes:
        - name: my-pv-storage
          persistentVolumeClaim:
            claimName: my-pv-claim
      containers:
        - name: peer1-${NEW_ORG_NAME}-couchdb
          image: couchdb:3.1.1
          env:
            - name: COUCHDB_USER
              value: nick
            - name: COUCHDB_PASSWORD
              value: "1234"
          volumeMounts:
            - mountPath: /opt/couchdb/data
              name: my-pv-storage
              subPath: state/${NEW_ORG_NAME}/peers/peer1-couchdb

EOT

cat <<EOT > $FOLDER_PATH/cli/cli-peer0-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cli-peer0-${NEW_ORG_NAME}-deployment
  labels: {
    component: peer0,
    type: cli,
    org: ${NEW_ORG_NAME}
  }
spec:
  replicas: 1
  selector:
    matchLabels:
      component: peer0
      type: cli
      org: ${NEW_ORG_NAME}
  template:
    metadata:
      labels:
        component: peer0
        type: cli
        org: ${NEW_ORG_NAME}
    spec:
      volumes:
        - name: my-pv-storage
          persistentVolumeClaim:
            claimName: my-pv-claim
      containers:
        - name: peer0-${NEW_ORG_NAME}
          image: hyperledger/fabric-tools:2.2.1
          workingDir: /opt/gopath/src/github.com/hyperledger/fabric/peer
          command: ["sleep"]
          args: ["infinity"]
          env:
            - name: GOPATH
              value: /opt/gopath
            - name: CORE_PEER_ADDRESSAUTODETECT
              value: "true"
            - name: CORE_PEER_ID
              value: cli-peer0-${NEW_ORG_NAME}
            - name: CORE_PEER_ADDRESS
              value: peer0-${NEW_ORG_NAME}-service:7051
            - name: CORE_PEER_LOCALMSPID
              value: ${NEW_ORG_NAME}
            - name: CORE_PEER_MSPCONFIGPATH
              value: /etc/hyperledger/fabric/msp/users/Admin@${NEW_ORG_NAME}/msp
            - name: CORE_PEER_TLS_ENABLED
              value: "true"
            - name: CORE_PEER_TLS_CERT_FILE
              value: /etc/hyperledger/fabric/tls/server.crt
            - name: CORE_PEER_TLS_KEY_FILE
              value: /etc/hyperledger/fabric/tls/server.key
            - name: CORE_PEER_TLS_ROOTCERT_FILE
              value: /etc/hyperledger/fabric/tls/ca.crt
          volumeMounts:
            - mountPath: /opt/gopath/src/github.com/hyperledger/fabric/peer/orderer
              name: my-pv-storage
              subPath: files/orderer
            - mountPath: /opt/gopath/src/resources
              name: my-pv-storage
              subPath: files/chaincode/resources
            - mountPath: /opt/gopath/src/resource_types
              name: my-pv-storage
              subPath: files/chaincode/resource_types
            - mountPath: /opt/gopath/src/github.com/hyperledger/fabric/peer/channels
              name: my-pv-storage
              subPath: files/channels
            - mountPath: /etc/hyperledger/fabric/msp
              name: my-pv-storage
              subPath: files/crypto-config/peerOrganizations/${NEW_ORG_NAME}
            - mountPath: /etc/hyperledger/fabric/tls
              name: my-pv-storage
              subPath: files/crypto-config/peerOrganizations/${NEW_ORG_NAME}/peers/peer0-${NEW_ORG_NAME}/tls
            - mountPath: /etc/hyperledger/orderers
              name: my-pv-storage
              subPath: files/crypto-config/ordererOrganizations/orderer

EOT

cat <<EOT > $FOLDER_PATH/cli/cli-peer1-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cli-peer1-${NEW_ORG_NAME}-deployment
  labels: {
    component: peer1,
    type: cli,
    org: ${NEW_ORG_NAME}
  }
spec:
  replicas: 1
  selector:
    matchLabels:
      component: peer1
      type: cli
      org: ${NEW_ORG_NAME}
  template:
    metadata:
      labels:
        component: peer1
        type: cli
        org: ${NEW_ORG_NAME}
    spec:
      volumes:
        - name: my-pv-storage
          persistentVolumeClaim:
            claimName: my-pv-claim
      containers:
        - name: peer1-${NEW_ORG_NAME}
          image: hyperledger/fabric-tools:2.2.1
          workingDir: /opt/gopath/src/github.com/hyperledger/fabric/peer
          command: ["sleep"]
          args: ["infinity"]
          env:
            - name: GOPATH
              value: /opt/gopath
            - name: CORE_PEER_ADDRESSAUTODETECT
              value: "true"
            - name: CORE_PEER_ID
              value: cli-peer1-${NEW_ORG_NAME}
            - name: CORE_PEER_ADDRESS
              value: peer1-${NEW_ORG_NAME}-service:7051
            - name: CORE_PEER_LOCALMSPID
              value: ${NEW_ORG_NAME}
            - name: CORE_PEER_MSPCONFIGPATH
              value: /etc/hyperledger/fabric/msp/users/Admin@${NEW_ORG_NAME}/msp
            - name: CORE_PEER_TLS_ENABLED
              value: "true"
            - name: CORE_PEER_TLS_CERT_FILE
              value: /etc/hyperledger/fabric/tls/server.crt
            - name: CORE_PEER_TLS_KEY_FILE
              value: /etc/hyperledger/fabric/tls/server.key
            - name: CORE_PEER_TLS_ROOTCERT_FILE
              value: /etc/hyperledger/fabric/tls/ca.crt
          volumeMounts:
            - mountPath: /opt/gopath/src/github.com/hyperledger/fabric/peer/orderer
              name: my-pv-storage
              subPath: files/orderer
            - mountPath: /opt/gopath/src/resources
              name: my-pv-storage
              subPath: files/chaincode/resources
            - mountPath: /opt/gopath/src/resource_types
              name: my-pv-storage
              subPath: files/chaincode/resource_types
            - mountPath: /opt/gopath/src/github.com/hyperledger/fabric/peer/channels
              name: my-pv-storage
              subPath: files/channels
            - mountPath: /etc/hyperledger/fabric/msp
              name: my-pv-storage
              subPath: files/crypto-config/peerOrganizations/${NEW_ORG_NAME}
            - mountPath: /etc/hyperledger/fabric/tls
              name: my-pv-storage
              subPath: files/crypto-config/peerOrganizations/${NEW_ORG_NAME}/peers/peer1-${NEW_ORG_NAME}/tls
            - mountPath: /etc/hyperledger/orderers
              name: my-pv-storage
              subPath: files/crypto-config/ordererOrganizations/orderer

EOT

cat <<EOT > $FOLDER_PATH/peers/services.yaml
apiVersion: v1
kind: Service
metadata:
  name: peer0-${NEW_ORG_NAME}-service
  labels: {
    component: peer0,
    type: peer,
    org: ${NEW_ORG_NAME}
  }
spec:
  type: ClusterIP
  selector:
    component: peer0
    type: peer
    org: ${NEW_ORG_NAME}
  ports:
    - port: 7051
      targetPort: 7051
---
apiVersion: v1
kind: Service
metadata:
  name: peer1-${NEW_ORG_NAME}-service
  labels: {
    component: peer1,
    type: peer,
    org: ${NEW_ORG_NAME}
  }
spec:
  type: ClusterIP
  selector:
    component: peer1
    type: peer
    org: ${NEW_ORG_NAME}
  ports:
    - port: 7051
      targetPort: 7051
EOT

cat <<EOT > $FOLDER_PATH/peers/peer0-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: peer0-${NEW_ORG_NAME}-deployment
  labels: {
    component: peer0,
    type: peer,
    org: ${NEW_ORG_NAME}
  }
spec:
  replicas: 1
  selector:
    matchLabels:
      component: peer0
      type: peer
      org: ${NEW_ORG_NAME}
  template:
    metadata:
      labels:
        component: peer0
        type: peer
        org: ${NEW_ORG_NAME}
    spec:
      volumes:
        - name: my-pv-storage
          persistentVolumeClaim:
            claimName: my-pv-claim
        - name: host
          hostPath:
            path: /var/run
      containers:
        - name: peer0-${NEW_ORG_NAME}
          image: hyperledger/fabric-peer:2.2.1
          workingDir: /opt/gopath/src/github.com/hyperledger/fabric/peer
          command: ["peer"]
          args: ["node","start"]
          env:
            # - name: FABRIC_LOGGING_SPEC
            #   value: DEBUG
            - name: CORE_VM_ENDPOINT
              value: unix:///var/run/docker.sock
            - name: CORE_PEER_ADDRESSAUTODETECT
              value: "true"
            - name: CORE_VM_DOCKER_ATTACHOUT
              value: "true"
            - name: CORE_PEER_ID
              value: peer0-${NEW_ORG_NAME}-service
            - name: CORE_PEER_LISTENADDRESS
              value: 0.0.0.0:7051
            - name: CORE_PEER_GOSSIP_BOOTSTRAP
              value: peer1-${NEW_ORG_NAME}-service:7051
            - name: CORE_PEER_GOSSIP_EXTERNALENDPOINT
              value: peer0-${NEW_ORG_NAME}-service:7051
            - name: CORE_PEER_GOSSIP_ENDPOINT
              value: peer0-${NEW_ORG_NAME}-service:7051
            - name: CORE_PEER_CHAINCODELISTENADDRESS
              value: 0.0.0.0:7052
            - name: CORE_PEER_LOCALMSPID
              value: ${NEW_ORG_NAME}
            - name: CORE_PEER_ENDORSER_ENABLED
              value: "true"
            # - name: CORE_PEER_GOSSIP_USELEADERELECTION
            #   value: "true"
            - name: CORE_PEER_TLS_ENABLED
              value: "true"
            - name: CORE_PEER_TLS_CERT_FILE
              value: /etc/hyperledger/fabric/tls/server.crt
            - name: CORE_PEER_TLS_KEY_FILE
              value: /etc/hyperledger/fabric/tls/server.key
            - name: CORE_PEER_TLS_ROOTCERT_FILE
              value: /etc/hyperledger/fabric/tls/ca.crt
            - name: CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS
              value: peer0-${NEW_ORG_NAME}-couchdb:5984
            - name: CORE_LEDGER_STATE_STATEDATABASE
              value: CouchDB
            - name: CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME
              value: nick
            - name: CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD
              value: "1234"
          volumeMounts:
            - mountPath: /var/run
              name: host
            - mountPath: /etc/hyperledger/fabric/msp
              name: my-pv-storage
              subPath: files/crypto-config/peerOrganizations/${NEW_ORG_NAME}/peers/peer0-${NEW_ORG_NAME}/msp
            - mountPath: /etc/hyperledger/fabric/tls
              name: my-pv-storage
              subPath: files/crypto-config/peerOrganizations/${NEW_ORG_NAME}/peers/peer0-${NEW_ORG_NAME}/tls
            - mountPath: /scripts
              name: my-pv-storage
              subPath: files/scripts
            - mountPath: /etc/hyperledger/orderers
              name: my-pv-storage
              subPath: files/crypto-config/ordererOrganizations/orderer

EOT

cat <<EOT > $FOLDER_PATH/peers/peer1-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: peer1-${NEW_ORG_NAME}-deployment
  labels: {
    component: peer1,
    type: peer,
    org: ${NEW_ORG_NAME}
  }
spec:
  replicas: 1
  selector:
    matchLabels:
      component: peer1
      type: peer
      org: ${NEW_ORG_NAME}
  template:
    metadata:
      labels:
        component: peer1
        type: peer
        org: ${NEW_ORG_NAME}
    spec:
      volumes:
        - name: my-pv-storage
          persistentVolumeClaim:
            claimName: my-pv-claim
        - name: host
          hostPath:
            path: /var/run
      containers:
        - name: peer1-${NEW_ORG_NAME}
          image: hyperledger/fabric-peer:2.2.1
          workingDir: /opt/gopath/src/github.com/hyperledger/fabric/peer
          command: ["peer"]
          args: ["node","start"]
          env:
            - name: CORE_VM_ENDPOINT
              value: unix:///var/run/docker.sock
            - name: CORE_PEER_ADDRESSAUTODETECT
              value: "true"
            - name: CORE_VM_DOCKER_ATTACHOUT
              value: "true"
            - name: CORE_PEER_ID
              value: peer1-${NEW_ORG_NAME}-service
            - name: CORE_PEER_LISTENADDRESS
              value: 0.0.0.0:7051
            - name: CORE_PEER_GOSSIP_BOOTSTRAP
              value: peer0-${NEW_ORG_NAME}-service:7051
            - name: CORE_PEER_GOSSIP_EXTERNALENDPOINT
              value: peer1-${NEW_ORG_NAME}-service:7051
            - name: CORE_PEER_GOSSIP_ENDPOINT
              value: peer1-${NEW_ORG_NAME}-service:7051
            - name: CORE_PEER_CHAINCODELISTENADDRESS
              value: 0.0.0.0:7052
            - name: CORE_PEER_LOCALMSPID
              value: ${NEW_ORG_NAME}
            - name: CORE_PEER_ENDORSER_ENABLED
              value: "true"
            # - name: CORE_PEER_GOSSIP_USELEADERELECTION
            #   value: "true"
            - name: CORE_PEER_TLS_ENABLED
              value: "true"
            - name: CORE_PEER_TLS_CERT_FILE
              value: /etc/hyperledger/fabric/tls/server.crt
            - name: CORE_PEER_TLS_KEY_FILE
              value: /etc/hyperledger/fabric/tls/server.key
            - name: CORE_PEER_TLS_ROOTCERT_FILE
              value: /etc/hyperledger/fabric/tls/ca.crt
            - name: CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS
              value: peer1-${NEW_ORG_NAME}-couchdb:5984
            - name: CORE_LEDGER_STATE_STATEDATABASE
              value: CouchDB
            - name: CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME
              value: nick
            - name: CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD
              value: "1234"
          volumeMounts:
            - mountPath: /var/run
              name: host
            - mountPath: /etc/hyperledger/fabric/msp
              name: my-pv-storage
              subPath: files/crypto-config/peerOrganizations/${NEW_ORG_NAME}/peers/peer1-${NEW_ORG_NAME}/msp
            - mountPath: /etc/hyperledger/fabric/tls
              name: my-pv-storage
              subPath: files/crypto-config/peerOrganizations/${NEW_ORG_NAME}/peers/peer1-${NEW_ORG_NAME}/tls
            - mountPath: /scripts
              name: my-pv-storage
              subPath: files/scripts
            - mountPath: /etc/hyperledger/orderers
              name: my-pv-storage
              subPath: files/crypto-config/ordererOrganizations/orderer

EOT
```

Let's bring up the third org (NOTE: Ensure each group is up before adding the next)
```bash
kubectl apply -f $FOLDER_PATH/couchdb
sleep 30
kubectl apply -f $FOLDER_PATH/peers
sleep 30
kubectl apply -f $FOLDER_PATH/cli
```

Time to join the peers to the network
```bash

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel join -b channels/mainchannel.block'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel join -b channels/mainchannel.block'
```

Time to install resource_types chaincode for the 3rd org (Need to use the NEXT seq number)
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_2' &
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_2' &
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_2' &
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_2' &
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_2' &
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_2'


kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt' &
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz' &
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt' &
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz' &
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt' &
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz'


kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --channelID mainchannel --collections-config /opt/gopath/src/resource_types/collections-config.json --name resource_types --version 2.0 --sequence 2 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')' &
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 2.0 --sequence 2 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')' &
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 2.0 --sequence 2 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode commit -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --channelID mainchannel --collections-config /opt/gopath/src/resource_types/collections-config.json --name resource_types --version 2.0 --sequence 2'
```

Time to install resources chaincode for the 3rd org (Need to use the NEXT seq number)
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_2'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_2'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_2'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_2'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_2'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_2'


kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz &> pkg.txt'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz &> pkg.txt'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz &> pkg.txt'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz'


kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --channelID mainchannel --collections-config /opt/gopath/src/resources/collections-config.json --name resources --version 2.0 --sequence 2 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resources/collections-config.json --channelID mainchannel --name resources --version 2.0 --sequence 2 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resources/collections-config.json --channelID mainchannel --name resources --version 2.0 --sequence 2 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode commit -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --channelID mainchannel --collections-config /opt/gopath/src/resources/collections-config.json --name resources --version 2.0 --sequence 2'
```


## UNNECESSARY NOTES!!!
NOTE: This is not necessary - just here for a note (DO NOT DO THIS!!)
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel update -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem -c mainchannel -f channels/hp-anchors.tx'
```