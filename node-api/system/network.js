const fabric = require('fabric-network');
const fs = require('fs');

const walletDirectoryPath = './system/walletstore'
const connectionProfilePath = `./system/configs/${process.env.ENV == 'dev' ? 'network-local' : 'network'}.json`;
const admin1org1MSPPath = `${process.env.ENV == 'dev' ? '../crypto-config' : '/tmp/crypto'}/peerOrganizations/org1/users/Admin@org1/msp`
const certPath = `${admin1org1MSPPath}/signcerts/Admin@org1-cert.pem`
const pvtKeyPath = `${admin1org1MSPPath}/keystore/pvt-cert.pem`
let mainchannelNetwork = null;

async function setup() {
    if (mainchannelNetwork) {
        return mainchannelNetwork;
    }

    const cert = (await fs.promises.readFile(certPath)).toString();
    const pvkey = (await fs.promises.readFile(pvtKeyPath)).toString();

    // Connect to a gateway peer
    const wallet = await fabric.Wallets.newFileSystemWallet(walletDirectoryPath);
    const identity = {
        credentials: {
            certificate: cert,
            privateKey: pvkey,
        },
        mspId: 'org1',
        type: 'X.509',
    };
    await wallet.put('admin', identity);
    const gatewayOptions = {
        identity: 'admin', // Previously imported identity
        wallet,
        discovery: {
            asLocalhost: true,
            enabled: false
        }
    };
    // read a common connection profile in json format
    const data = fs.readFileSync(connectionProfilePath);
    const connectionProfile = JSON.parse(data);

    // use the loaded connection profile
    const gateway = new fabric.Gateway();
    await gateway.connect(connectionProfile, gatewayOptions);

    // Obtain the smart contract with which our application wants to interact
    mainchannelNetwork = await gateway.getNetwork('mainchannel');

    return mainchannelNetwork;
}

module.exports.setup = setup;