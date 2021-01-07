Adding an org
======================

Create folder to store data in
```bash
export ENVIRONMENT=production
export NEW_ORG_NAME=hp
export FOLDER_PATH=configs/${ENVIRONMENT}/${NEW_ORG_NAME}
rm -rf $FOLDER_PATH
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

Next, we need to update the private data
```bash
cat <<EOT > chaincode/resource_types/collections-config.json
[
    {
        "name": "ibmResourceTypesPrivateData",
        "policy": "OR('ibm.member')",
        "requiredPeerCount": 1,
        "maxPeerCount": 3,
        "blockToLive": 0
    },
    {
        "name": "oracleResourceTypesPrivateData",
        "policy": "OR('oracle.member')",
        "requiredPeerCount": 1,
        "maxPeerCount": 3,
        "blockToLive": 0
    },
    {
        "name": "${NEW_ORG_NAME}ResourceTypesPrivateData",
        "policy": "OR('${NEW_ORG_NAME}.member')",
        "requiredPeerCount": 1,
        "maxPeerCount": 3,
        "blockToLive": 0
    }
]
EOT

cat <<EOT > chaincode/resources/collections-config.json
[
    {
        "name": "ibmResourcesPrivateData",
        "policy": "OR('ibm.member')",
        "requiredPeerCount": 1,
        "maxPeerCount": 3,
        "blockToLive": 0
    },
    {
        "name": "oracleResourcesPrivateData",
        "policy": "OR('oracle.member')",
        "requiredPeerCount": 1,
        "maxPeerCount": 3,
        "blockToLive": 0
    },
    {
        "name": "${NEW_ORG_NAME}ResourcesPrivateData",
        "policy": "OR('${NEW_ORG_NAME}.member')",
        "requiredPeerCount": 1,
        "maxPeerCount": 3,
        "blockToLive": 0
    }
]
EOT
```

Start hp ca and create certs
```bash
kubectl apply -f ${FOLDER_PATH}/cas
sleep 20
```

Time to copy over changed files
```bash
kubectl cp ${FOLDER_PATH}/configtx.yaml $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files
kubectl cp ./chaincode/resources $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files/chaincode
kubectl cp ./chaincode/resource_types $(kubectl get pods -o=name | grep example1 | sed "s/^.\{4\}//"):/host/files/chaincode
```

Get the config from the configtx (NOTE: There isn't really a way to pass in ENV into a kubectl command easily. Just do it manually) (In this example, I just hardcoded the orgName in the configtxgen command for "hp")
TODO: Create a script for this with the org name and then copy it and run it in the container
```bash
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

Time to install resource_types chaincode for the 3rd org: TODO: figure out how to get the current sequence number and use it for this new org
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_2'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_2'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_2'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_2'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_2'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode package resource_types.tar.gz --path /opt/gopath/src/resource_types --lang golang --label resource_types_2'


kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz &> pkg.txt'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer1-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode install resource_types.tar.gz'


kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 2.0 --sequence 2 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 2.0 --sequence 2 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 2.0 --sequence 2 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'

kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode commit -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resource_types/collections-config.json --channelID mainchannel --name resource_types --version 2.0 --sequence 2'
```

Lets test this chaincode
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resource_types -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
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


kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resources/collections-config.json --channelID mainchannel --name resources --version 2.0 --sequence 2 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resources/collections-config.json --channelID mainchannel --name resources --version 2.0 --sequence 2 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode approveformyorg -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resources/collections-config.json --channelID mainchannel --name resources --version 2.0 --sequence 2 --package-id $(tail -n 1 pkg.txt | awk '\''NF>1{print $NF}'\'')'



kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer lifecycle chaincode commit -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem --collections-config /opt/gopath/src/resources/collections-config.json --channelID mainchannel --name resources --version 2.0 --sequence 2'
```

Lets test this chaincode
```bash
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode invoke -C mainchannel -n resources -c '\''{"Args":["Create","PCs","1"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode invoke -C mainchannel -n resources -c '\''{"Args":["Create","Printers","1"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
sleep 5
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-ibm-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-oracle-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
kubectl exec -it $(kubectl get pods -o=name | grep cli-peer0-hp-deployment | sed "s/^.\{4\}//") -- bash -c 'peer chaincode query -C mainchannel -n resources -c '\''{"Args":["Index"]}'\'' -o orderer0-service:7050 --tls --cafile=/etc/hyperledger/orderers/msp/tlscacerts/orderers-ca-service-7054.pem'
```

