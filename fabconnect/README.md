Fabconnect
=============

## Expose Ports
kubectl expose deployment orderers-ca --port=7054 --target-port=7054
kubectl expose deployment ibm-ca --port=7054 --target-port=7055
kubectl expose deployment oracle-ca --port=7054 --target-port=7056

kubectl expose deployment orderer0-deployment --port=7050 --target-port=7050
kubectl expose deployment orderer1-deployment --port=7050 --target-port=7051
kubectl expose deployment orderer2-deployment --port=7050 --target-port=7052

kubectl expose deployment peer0-ibm-deployment --port=7051 --target-port=7060
kubectl expose deployment peer1-ibm-deployment --port=7051 --target-port=7061
kubectl expose deployment peer0-oracle-deployment --port=7051 --target-port=7062
kubectl expose deployment peer1-oracle-deployment --port=7051 --target-port=7063

## Curl command
