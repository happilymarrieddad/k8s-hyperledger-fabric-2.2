K8s Hyperledger 2.2 Network
======================================

## Important things to read
- [Orderering Service](https://hyperledger-fabric.readthedocs.io/en/release-2.2/orderer/ordering_service.html)
- [Peers](https://hyperledger-fabric.readthedocs.io/en/release-2.2/peers/peers.html)

## Getting Started
- [My bash](https://github.com/ohmyzsh/ohmyzsh)
- [Install Hyperledger Deps](https://hyperledger-fabric.readthedocs.io/en/release-2.2/install.html)
- Install fabric binaries and images
```bash
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.2.1 1.4.7
```
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

[Click Here](https://github.com/happilymarrieddad/k8s-hyperledger-fabric-2.2/blob/master/DOCKER_COMPOSE.md)

## Kubernetes - Minikube (Local)

[Click Here](https://github.com/happilymarrieddad/k8s-hyperledger-fabric-2.2/blob/master/MINIKUBE.md)

## Kubernetes - Minikube (Production)

[Click Here](https://github.com/happilymarrieddad/k8s-hyperledger-fabric-2.2/blob/master/PRODUCTION.md)

## Some other ways to run the network (good place to look for automation scripts in cli-peer start scripts)

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

[Click Here](https://github.com/happilymarrieddad/k8s-hyperledger-fabric-2.2/blob/master/UPDATING_CHAINCODE.md)

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

[Click Here](https://github.com/happilymarrieddad/k8s-hyperledger-fabric-2.2/blob/master/ADDING_AN_ORG.md)
