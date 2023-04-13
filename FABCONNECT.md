Fabconnect
=================================

## Example Post
```bash
curl -X POST -d '{"headers": {"type": "SendTransaction"},"func": "Index","args": []}' -H 'Context-Type: application/json' 'http://localhost:3000/transactions?fly-sync=true&fly-signer=Admin&fly-channel=mainchannel&fly-chaincode=resources'
```

{"headers": {"type": "SendTransaction"},"func": "Create","args": ["1", "CPUs Example", "1"]}

"Create","1","CPUs","1"

## MINIKUBE
Start up the fabconnect container
```bash
kubectl apply -f network/minikube/fabconnect
```

If you shut down the CA before need to start it up again
```bash
kubectl apply -f network/minikube/cas/ibm-ca-deployment.yaml 
kubectl apply -f network/minikube/cas/ibm-ca-service.yaml 
```

time to copy over the configs
```bash
kubectl cp ./network/minikube/fabconnect/config $(kubectl get pods -o=name | grep fabconnect | sed "s/^.\{4\}//"):/fabconnect
```
