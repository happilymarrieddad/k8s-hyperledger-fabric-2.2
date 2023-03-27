Fabconnect
=================================

## Example Post
```bash
curl -X POST -d '{"headers": {"type": "SendTransaction"},"func": "Create","args": ["1", "CPUs Example", "1"]}' -H 'Context-Type: application/json' 'http://localhost:3000/transactions?fly-sync=true&fly-signer=Admin&fly-channel=mainchannel&fly-chaincode=resources'
```

{"headers": {"type": "SendTransaction"},"func": "Create","args": ["1", "CPUs Example", "1"]}

"Create","1","CPUs","1"