require('dotenv').config();
const hre = require("hardhat");

async function main() {
    const contractAddress = process.env.BLACKLIST_MANAGER_ADDRESS;
    if (!contractAddress) {
        throw new Error("BLACKLIST_MANAGER_ADDRESS is not set in .env");
    }

    const [deployer] = await hre.ethers.getSigners();
    const BlacklistManager = await hre.ethers.getContractFactory("BlacklistManager");
    const contract = BlacklistManager.attach(contractAddress);

    const dummyData = [
        { keyword: "phishing", username: "alice" },
        { keyword: "scam", username: "bob" },
        { keyword: "malware", username: "charlie" },
        { keyword: "ransomware", username: "dave" },
        { keyword: "fraud", username: "eve" }
    ];

    for (const entry of dummyData) {
        try {
            const tx = await contract.addBlacklistKeyword(entry.keyword, entry.username);
            await tx.wait();
            console.log(`Added: ${entry.keyword} by ${entry.username}`);
        } catch (err) {
            console.error(`Failed to add ${entry.keyword}:`, err.message);
        }
    }
}

main()
    .then(() => process.exit(0))
    .catch((err) => {
        console.error(err);
        process.exit(1);
    }); 