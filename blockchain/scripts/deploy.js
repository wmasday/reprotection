const hre = require("hardhat");

async function main() {
    console.log("Deploying BlacklistManager contract...");

    const BlacklistManager = await hre.ethers.getContractFactory("BlacklistManager");
    const blacklistManager = await BlacklistManager.deploy();

    await blacklistManager.waitForDeployment();

    const address = await blacklistManager.getAddress();
    console.log("BlacklistManager deployed to:", address);

    // Save the contract address to a file for the API to use
    const fs = require('fs');
    const contractInfo = {
        address: address,
        network: hre.network.name,
        deployedAt: new Date().toISOString()
    };

    fs.writeFileSync('./contract-address.json', JSON.stringify(contractInfo, null, 2));
    console.log("Contract address saved to contract-address.json");
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    }); 