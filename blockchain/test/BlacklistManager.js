const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("BlacklistManager", function () {
    let BlacklistManager, blacklistManager, owner, addr1;

    beforeEach(async function () {
        [owner, addr1] = await ethers.getSigners();
        BlacklistManager = await ethers.getContractFactory("BlacklistManager");
        blacklistManager = await BlacklistManager.deploy();
        await blacklistManager.deployed();
    });

    it("Should deploy and have zero stats initially", async function () {
        const stats = await blacklistManager.getStatistics();
        expect(stats.totalApply).to.equal(0);
        expect(stats.totalVerify).to.equal(0);
        expect(stats.totalBlacklistEntries).to.equal(0);
    });

    it("Should add a blacklist keyword and update stats", async function () {
        await blacklistManager.addBlacklistKeyword("phishing", "alice");
        const entry = await blacklistManager.getBlacklistEntry("phishing");
        expect(entry.keywordName).to.equal("phishing");
        expect(entry.isActive).to.equal(true);
        const stats = await blacklistManager.getStatistics();
        expect(stats.totalBlacklistEntries).to.equal(1);
    });

    it("Should increment apply and verify counts", async function () {
        await blacklistManager.incrementApply();
        await blacklistManager.incrementVerify();
        const stats = await blacklistManager.getStatistics();
        expect(stats.totalApply).to.equal(1);
        expect(stats.totalVerify).to.equal(1);
    });
}); 