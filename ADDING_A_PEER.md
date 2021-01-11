Adding a peer
========================

Create folder to store data in
```bash
export ENVIRONMENT=production
export EXISTING_ORG_NAME=ibm
export FOLDER_PATH=configs/${ENVIRONMENT}/${EXISTING_ORG_NAME}
rm -rf $FOLDER_PATH
mkdir -p $FOLDER_PATH/cas
mkdir -p $FOLDER_PATH/cli
mkdir -p $FOLDER_PATH/couchdb
mkdir -p $FOLDER_PATH/peers
mkdir -p $FOLDER_PATH/chaincode/resource_types
mkdir -p $FOLDER_PATH/chaincode/resources
```

Copy ca configs, set client to not do anything, then start them up
```bash
cp network/minikube/cas/${EXISTING_ORG_NAME}-ca-client-deployment.yaml ${FOLDER_PATH}/cas
cp network/minikube/cas/${EXISTING_ORG_NAME}-ca-deployment.yaml ${FOLDER_PATH}/cas
cp network/minikube/cas/${EXISTING_ORG_NAME}-ca-service.yaml ${FOLDER_PATH}/cas
```

There's already 2 peers so we need to use the new index to create more. i.e. 1 more node and starting at index 2
```bash
sed -i -e 's/"2"/"1"/g' ${FOLDER_PATH}/cas/${EXISTING_ORG_NAME}-ca-client-deployment.yaml
sed -i -e 's/"0"/"2"/g' ${FOLDER_PATH}/cas/${EXISTING_ORG_NAME}-ca-client-deployment.yaml
```

Now, time to start ca
```bash
kubectl apply -f ${FOLDER_PATH}/cas
```

```bash
cat <<EOT > $FOLDER_PATH/couchdb/services.yaml
apiVersion: v1
kind: Service
metadata:
  name: peer0-${EXISTING_ORG_NAME}-couchdb
  labels: {
    component: peer0,
    type: couchdb,
    org: ${EXISTING_ORG_NAME}
  }
spec:
  type: ClusterIP
  selector:
    component: peer0
    type: couchdb
    org: ${EXISTING_ORG_NAME}
  ports:
    - port: 5984
      targetPort: 5984
---
apiVersion: v1
kind: Service
metadata:
  name: peer1-${EXISTING_ORG_NAME}-couchdb
  labels: {
    component: peer1,
    type: couchdb,
    org: ${EXISTING_ORG_NAME}
  }
spec:
  type: ClusterIP
  selector:
    component: peer1
    type: couchdb
    org: ${EXISTING_ORG_NAME}
  ports:
    - port: 5984
      targetPort: 5984
---
apiVersion: v1
kind: Service
metadata:
  name: peer2-${EXISTING_ORG_NAME}-couchdb
  labels: {
    component: peer2,
    type: couchdb,
    org: ${EXISTING_ORG_NAME}
  }
spec:
  type: ClusterIP
  selector:
    component: peer2
    type: couchdb
    org: ${EXISTING_ORG_NAME}
  ports:
    - port: 5984
      targetPort: 5984

EOT

cp network/minikube/orgs/${EXISTING_ORG_NAME}/couchdb/peer0-couchdb-deployment.yaml $FOLDER_PATH/couchdb/peer2-couchdb-deployment.yaml
cp network/minikube/orgs/${EXISTING_ORG_NAME}/cli/cli-peer0-deployment.yaml $FOLDER_PATH/cli/cli-peer2-deployment.yaml


sed -i -e 's/peer0/peer2/g' $FOLDER_PATH/couchdb/peer2-couchdb-deployment.yaml
sed -i -e 's/peer0/peer2/g' $FOLDER_PATH/cli/cli-peer2-deployment.yaml

cat <<EOT > $FOLDER_PATH/peers/services.yaml
apiVersion: v1
kind: Service
metadata:
  name: peer0-${EXISTING_ORG_NAME}-service
  labels: {
    component: peer0,
    type: peer,
    org: ${EXISTING_ORG_NAME}
  }
spec:
  type: ClusterIP
  selector:
    component: peer0
    type: peer
    org: ${EXISTING_ORG_NAME}
  ports:
    - port: 7051
      targetPort: 7051
---
apiVersion: v1
kind: Service
metadata:
  name: peer1-${EXISTING_ORG_NAME}-service
  labels: {
    component: peer1,
    type: peer,
    org: ${EXISTING_ORG_NAME}
  }
spec:
  type: ClusterIP
  selector:
    component: peer1
    type: peer
    org: ${EXISTING_ORG_NAME}
  ports:
    - port: 7051
      targetPort: 7051
---
apiVersion: v1
kind: Service
metadata:
  name: 2-${EXISTING_ORG_NAME}-service
  labels: {
    component: 2,
    type: peer,
    org: ${EXISTING_ORG_NAME}
  }
spec:
  type: ClusterIP
  selector:
    component: 2
    type: peer
    org: ${EXISTING_ORG_NAME}
  ports:
    - port: 7051
      targetPort: 7051
EOT

cp $FOLDER_PATH/peers/peer2-deployment.yaml
sed -i -e 's/peer0/peer2/g' $FOLDER_PATH/peers/peer2-deployment.yaml

EOT
```

Time for the couchdb and cli to join
```bash
kubectl apply -f $FOLDER_PATH/couchdb
kubectl apply -f $FOLDER_PATH/cli
```

Time to bring up the peers (NEED TO WAIT FOR COUCHDB)
```bash
kubectl apply -f $FOLDER_PATH/peers
```

Time to join the peer(s) to the network
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer2-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer channel join -b channels/mainchannel.block'
```

Need to get the sequence number for the current chaincode. In my case it was 2. VERY IMPORTANT that you don't mess this up or you're going to have to install and instantiate for all the orgs over again.
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer2-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_2'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer2-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz'
```

Lets test this chaincode
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer2-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
```

Need to get the sequence number for the current chaincode. In my case it was 2. VERY IMPORTANT that you don't mess this up or you're going to have to install and instantiate for all the orgs over again.
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer2-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_2'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer2-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz'
```

Testing resource chaincode
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer2-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
```
