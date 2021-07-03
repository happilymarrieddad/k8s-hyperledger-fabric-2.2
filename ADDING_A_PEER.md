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
# Use these to change the command to sleep so you can bash in and examine what's going to happen before running
# the script
# sed -i -e 's/bash/sleep/g' ${FOLDER_PATH}/cas/${EXISTING_ORG_NAME}-ca-client-deployment.yaml
# sed -i -e 's/\/scripts\/start-org-client.sh/infinity/g' ${FOLDER_PATH}/cas/${EXISTING_ORG_NAME}-ca-client-deployment.yaml
sed -i -e 's/"2"/"1"/g' ${FOLDER_PATH}/cas/${EXISTING_ORG_NAME}-ca-client-deployment.yaml
sed -i -e 's/"0"/"2"/g' ${FOLDER_PATH}/cas/${EXISTING_ORG_NAME}-ca-client-deployment.yaml
```

Now, time to start ca
```bash
kubectl apply -f ${FOLDER_PATH}/cas
```

Wait for the ca to come online and generate the certs
```bash
sleep 60

```

Let's create the yaml files for the new peer
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
  name: peer2-${EXISTING_ORG_NAME}-service
  labels: {
    component: peer2,
    type: peer,
    org: ${EXISTING_ORG_NAME}
  }
spec:
  type: ClusterIP
  selector:
    component: peer2
    type: peer
    org: ${EXISTING_ORG_NAME}
  ports:
    - port: 7051
      targetPort: 7051
EOT

cat <<EOT > $FOLDER_PATH/peers/peer2-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: peer2-ibm-deployment
  labels: {
    component: peer2,
    type: peer,
    org: ibm
  }
spec:
  replicas: 1
  selector:
    matchLabels:
      component: peer2
      type: peer
      org: ibm
  template:
    metadata:
      labels:
        component: peer2
        type: peer
        org: ibm
    spec:
      volumes:
        - name: my-pv-storage
          persistentVolumeClaim:
            claimName: my-pv-claim
        - name: host
          hostPath:
            path: /var/run
      containers:
        - name: peer2-ibm
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
              value: peer2-ibm-service
            - name: CORE_PEER_LISTENADDRESS
              value: 0.0.0.0:7051
            - name: CORE_PEER_GOSSIP_BOOTSTRAP
              value: peer0-ibm-service:7051
            - name: CORE_PEER_GOSSIP_EXTERNALENDPOINT
              value: peer2-ibm-service:7051
            - name: CORE_PEER_GOSSIP_ENDPOINT
              value: peer2-ibm-service:7051
            - name: CORE_PEER_CHAINCODELISTENADDRESS
              value: 0.0.0.0:7052
            - name: CORE_PEER_LOCALMSPID
              value: ibm
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
              value: peer2-ibm-couchdb:5984
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
              subPath: files/crypto-config/peerOrganizations/ibm/peers/peer2-ibm/msp
            - mountPath: /etc/hyperledger/fabric/tls
              name: my-pv-storage
              subPath: files/crypto-config/peerOrganizations/ibm/peers/peer2-ibm/tls
            - mountPath: /scripts
              name: my-pv-storage
              subPath: files/scripts
            - mountPath: /etc/hyperledger/orderers
              name: my-pv-storage
              subPath: files/crypto-config/ordererOrganizations/orderer

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

Need to get the sequence number for the current chaincode. In my case it was 1. VERY IMPORTANT that you don't mess this up or you're going to have to install and instantiate for all the orgs over again.
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer2-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_1'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer2-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz'
```

Lets test this chaincode
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer2-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
```

Need to get the sequence number for the current chaincode. In my case it was 1. VERY IMPORTANT that you don't mess this up or you're going to have to install and instantiate for all the orgs over again.
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer2-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resources.tar.gz --path /opt/gopath/src/resources --lang golang --label resources_1'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer2-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resources.tar.gz'
```

Testing resource chaincode
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer2-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
```
