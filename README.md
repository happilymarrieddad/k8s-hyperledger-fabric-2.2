K8s Hyperledger 2.2 Network
======================================

## Getting Started
[My bash](https://github.com/ohmyzsh/ohmyzsh)
[Install Hyperledger Deps](https://hyperledger-fabric.readthedocs.io/en/release-2.2/install.html)
- curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.2.1 1.4.9

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
[Golang](https://golang.org/dl/)
[Nodejs](https://nodejs.org/en/)

## Docker (Local)

## Kubernetes - Minikube (Local)

## Commands
Clean all docker things
```bash
docker system prune 
docker container prune
docker volume prune 
docker rmi $(docker images -q) --force
```