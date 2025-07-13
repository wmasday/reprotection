const hre = require("hardhat");

async function main() {
    const BlacklistManager = await hre.ethers.getContractFactory("BlacklistManager");
    const blacklistManager = await BlacklistManager.deploy();
    await blacklistManager.waitForDeployment();
    console.log("BlacklistManager deployed to:", blacklistManager.target);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    }); 