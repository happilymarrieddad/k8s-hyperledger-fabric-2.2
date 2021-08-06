const express = require('express');
const router = express.Router();
const uuid = require('uuid');
const network = require('../system/network');

// Querying
router.get('/', async (req, res) => {
    mainNetwork = await network.setup();

    const contract = mainNetwork.getContract('resources');

    // Submit transactions for the smart contract
    const submitResult = await contract.evaluateTransaction('index').catch(err => {
        res.status(400).send(err)
    });

    // Error handled in catch
    if (!submitResult) {
        return
    }

    if (submitResult.statusCode && submitResult.statusCode === 400) {
        return res.status(400).send(submitResult.statusMessage)
    }

    // Remove the unnecessary quotes
    res.json(JSON.parse(submitResult.toString()));
})

router.get('/:id/transactions', async (req, res) => {
    if (!req || !req.params) {
        res.status(400).send('resource id required in url parameters');
        return;
    }

    mainNetwork = await network.setup();

    const contract = mainNetwork.getContract('resources');

    // Submit transactions for the smart contract
    const submitResult = await contract.evaluateTransaction('transactions', req.params.id).catch(err => res.status(400).send(err));

    // Remove the unnecessary quotes
    res.json(JSON.parse(submitResult.toString()));
})

router.get('/:id', async (req, res) => {
    if (!req || !req.params) {
        res.status(400).send('resource id required in url parameters');
        return;
    }

    mainNetwork = await network.setup();

    const contract = mainNetwork.getContract('resources');

    // Submit transactions for the smart contract
    const submitResult = await contract.evaluateTransaction('read', [req.params.id]).catch(err => res.status(400).send(err));

    // Remove the unnecessary quotes
    res.json(JSON.parse(submitResult.toString()));
})

// Commiting Transactions
router.post('/', async (req, res) => {
    if (!req || !req.body) {
        res.status(400).send('resource required in body');
        return;
    }

    mainNetwork = await network.setup();

    const contract = mainNetwork.getContract('resources');

    const newObj = {
        name: req.body.name,
        resource_type_id: req.body.resource_type_id
    }

    await contract.submitTransaction('create', '', newObj.name, newObj.resource_type_id)

    res.json(newObj)
})

router.put('/:id', async (req, res) => {
    if (!req || !req.params) {
        res.status(400).send('resource id required in url parameters');
        return;
    }

    if (!req || !req.body) {
        res.status(400).send('resource required in body');
        return;
    }

    const orgName = 'ibm';

    const newResource = {
        owner: orgName,
    }

    if (req.body.name) {
        newResource.name = req.body.name
    }
    
    if (req.body.resource_type_id) {
        newResource.resource_type_id = req.body.resource_type_id
    }

    if (newResource.name.length == 0 || newResource.resource_type_id < 1) {
        res.status(400).send('resource requires name and resource_type_id');
        return;
    }

    mainNetwork = await network.setup();

    // You can get MSPIds from the getway... it depends on your network
    /* 
    <ref *1> Channel {
  type: 'Channel',
  name: 'mainchannel',
  client: Client {
    type: 'Client',
    name: 'gateway client',
    mspid: null,
    _tls_mutual: {
      selfGenerated: true,
      clientKey: '-----BEGIN PRIVATE KEY-----\r\n' +
        'MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgT0PSO0ISr6fzX/EX\r\n' +
        'dc40BjUmLYfS/0DekNrHRlCIQjmhRANCAARfp6Uc1C8cRMHMVq4eONP0aQ4v9zag\r\n' +
        '5gJa+QUNNUNL/kGeCwTVhPg3aCX4zW6Nbrd1Cw7qWHcMUTKGHagMhcZP\r\n' +
        '-----END PRIVATE KEY-----\r\n',
      clientCert: '-----BEGIN CERTIFICATE-----\r\n' +
        'MIIBVDCB+6ADAgECAgEEMAoGCCqGSM49BAMCMBgxFjAUBgNVBAMMDWZhYnJpYy1j\r\n' +
        'b21tb24wIhgPMjAyMDEyMTQwMzM4MDBaGA8yMDIwMTIxNDIxNDEyMFowGDEWMBQG\r\n' +
        'A1UEAwwNZmFicmljLWNvbW1vbjBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABF+n\r\n' +
        'pRzULxxEwcxWrh440/RpDi/3NqDmAlr5BQ01Q0v+QZ4LBNWE+DdoJfjNbo1ut3UL\r\n' +
        'DupYdwxRMoYdqAyFxk+jMjAwMAwGA1UdEwEB/wQCMAAwCwYDVR0PBAQDAgbAMBMG\r\n' +
        'A1UdJQQMMAoGCCsGAQUFBwMCMAoGCCqGSM49BAMCA0gAMEUCIBupsTCPfzPeooEP\r\n' +
        'apVPpMRpb07umWRaI9ae2WerqqMEAiEAzURvWTGRRMXZp8BCVum354F4G5Yb3Nx/\r\n' +
        'N5lU/Mw1Idw=\r\n' +
        '-----END CERTIFICATE-----\r\n',
      clientCertHash: <Buffer 78 c2 90 c2 92 9f 62 ec 69 da d6 f1 33 86 e4 6c 6e 87 43 aa 75 ae 63 17 ff 9d cc 1a 04 9a 9c cc>
    },
    endorsers: Map(1) { 'peer0-ibm' => [Endorser] },
    committers: Map(1) { 'orderer0' => [Committer] },
    channels: Map(1) { 'mainchannel' => [Circular *1] },
    centralizedOptions: null
  },
  endorsers: Map(1) {
    'peer0-ibm' => Endorser {
      name: 'peer0-ibm',
      mspid: 'ibm',
      client: [Client],
      connected: true,
      connectAttempted: true,
      endpoint: [Endpoint],
      service: [ServiceClientImpl],
      serviceClass: [Function],
      type: 'Endorser',
      options: [Object],
      chaincodes: [Array],
      discovered: false
    }
  },
  committers: Map(1) {
    'orderer0' => Committer {
      name: 'orderer0',
      mspid: undefined,
      client: [Client],
      connected: true,
      connectAttempted: true,
      endpoint: [Endpoint],
      service: [ServiceClientImpl],
      serviceClass: [Function],
      type: 'Committer',
      options: [Object]
    }
  },
  msps: Map(0) {}
}
    */
    console.log(mainNetwork.getChannel())
    console.log(mainNetwork.getGateway()())

    const contract = mainNetwork.getContract('resources');

    try {
        await contract.submitTransaction('update', req.params.id, `${newResource.name}`, `${newResource.resource_type_id}`, `${newResource.owner}`)

        res.json(newResource);
    } catch(err) {
        res.status(400).send(err)
    }
})

router.delete('/:id', async (req, res) => {
    if (!req || !req.params) {
        res.status(400).send('resource id required in url parameters');
        return;
    }

    mainNetwork = await network.setup();

    const contract = mainNetwork.getContract('resources');

    // Submit transactions for the smart contract
    const submitResult = await contract.evaluateTransaction('delete', [req.params.id]).catch(err => res.status(400).send(err));

    // Remove the unnecessary quotes
    res.send("success");
})

module.exports = router;
