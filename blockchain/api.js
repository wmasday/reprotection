const express = require('express');
const { ethers } = require('ethers');
const cors = require('cors');
require('dotenv').config();

const app = express();
const PORT = process.env.BLOCKCHAIN_API_PORT || 3001;

app.use(cors());
app.use(express.json());

// Load contract ABI and address
const fs = require('fs');
const path = require('path');

let contractAddress;
try {
    const contractInfo = JSON.parse(fs.readFileSync('./contract-address.json', 'utf8'));
    contractAddress = contractInfo.address;
    console.log('Contract address loaded:', contractAddress);
} catch (error) {
    console.error('Error loading contract address:', error);
    console.log('Please deploy the contract first using: npx hardhat run scripts/deploy.js --network localhost');
    process.exit(1);
}

// Load contract ABI
const contractArtifact = JSON.parse(fs.readFileSync('./artifacts/contracts/BlacklistManager.sol/BlacklistManager.json', 'utf8'));
const contractABI = contractArtifact.abi;

// Connect to local blockchain
const provider = new ethers.JsonRpcProvider(process.env.BLOCKCHAIN_RPC_URL || 'http://127.0.0.1:8545');
const wallet = new ethers.Wallet(process.env.PRIVATE_KEY || '0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80', provider);
const contract = new ethers.Contract(contractAddress, contractABI, wallet);

console.log('Connected to blockchain at:', process.env.BLOCKCHAIN_RPC_URL || 'http://127.0.0.1:8545');
console.log('Using wallet address:', wallet.address);

// Get all keywords from blockchain (active and inactive)
app.get('/blacklist', async (req, res) => {
    try {
        const keywords = await contract.getAllKeywords();
        const entries = [];

        for (let i = 0; i < keywords.length; i++) {
            const keyword = keywords[i];
            entries.push({
                index: i,
                keyword: keyword.keyword,
                createdBy: keyword.creator,
                username: keyword.username,
                timestamp: keyword.timestamp.toString(),
                isActive: keyword.isActive
            });
        }

        res.json(entries);
    } catch (err) {
        console.error('Error getting keywords:', err);
        res.status(500).json({ error: err.message });
    }
});

// Get specific keyword by index
app.get('/blacklist/:index', async (req, res) => {
    try {
        const index = parseInt(req.params.index);
        const keywordData = await contract.getKeyword(index);

        res.json({
            index: index,
            keyword: keywordData.keyword,
            createdBy: keywordData.creator,
            username: keywordData.username,
            timestamp: keywordData.timestamp.toString(),
            isActive: keywordData.isActive
        });
    } catch (err) {
        console.error('Error getting keyword:', err);
        res.status(500).json({ error: err.message });
    }
});

// Get specific keyword
app.get('/blacklist/:keyword', async (req, res) => {
    try {
        const keyword = req.params.keyword;
        const keywordData = await contract.getKeyword(keyword);

        if (keywordData.creator === ethers.ZeroAddress) {
            return res.status(404).json({ error: 'Keyword not found' });
        }

        res.json({
            keyword: keywordData.keyword,
            createdBy: keywordData.creator,
            username: keywordData.username,
            timestamp: keywordData.timestamp.toString(),
            isActive: keywordData.isActive
        });
    } catch (err) {
        console.error('Error getting keyword:', err);
        res.status(500).json({ error: err.message });
    }
});

// Add keyword to blockchain
app.post('/blacklist/add', async (req, res) => {
    try {
        const { keyword, username } = req.body;

        if (!keyword || !username) {
            return res.status(400).json({ error: 'Keyword and username are required' });
        }

        if (keyword.length > 100) {
            return res.status(400).json({ error: 'Keyword too long (max 100 characters)' });
        }

        const tx = await contract.addKeyword(keyword, username);
        await tx.wait();

        res.json({
            success: true,
            message: 'Keyword added to blockchain successfully',
            transactionHash: tx.hash
        });
    } catch (err) {
        console.error('Error adding keyword:', err);
        res.status(500).json({ error: err.message });
    }
});

// Remove keyword from blockchain
app.post('/blacklist/remove', async (req, res) => {
    try {
        const { index, username } = req.body;

        if (index === undefined || !username) {
            return res.status(400).json({ error: 'Index and username are required' });
        }

        const tx = await contract.removeKeyword(index);
        await tx.wait();

        res.json({
            success: true,
            message: 'Keyword deactivated successfully',
            transactionHash: tx.hash
        });
    } catch (err) {
        console.error('Error removing keyword:', err);
        res.status(500).json({ error: err.message });
    }
});

// Toggle keyword active/inactive status
app.post('/blacklist/toggle', async (req, res) => {
    try {
        const { index, username } = req.body;

        if (index === undefined || !username) {
            return res.status(400).json({ error: 'Index and username are required' });
        }

        const tx = await contract.toggleKeyword(index);
        await tx.wait();

        // Get the new status
        const keywordData = await contract.getKeyword(index);
        const action = keywordData.isActive ? 'activated' : 'deactivated';

        res.json({
            success: true,
            message: `Keyword ${action} successfully`,
            isActive: keywordData.isActive,
            transactionHash: tx.hash
        });
    } catch (err) {
        console.error('Error toggling keyword:', err);
        res.status(500).json({ error: err.message });
    }
});

// Check if keyword is active
app.get('/blacklist/:index/isActive', async (req, res) => {
    try {
        const index = parseInt(req.params.index);
        const isActive = await contract.isKeywordActive(index);

        res.json({
            index: index,
            isActive: isActive
        });
    } catch (err) {
        console.error('Error checking keyword status:', err);
        res.status(500).json({ error: err.message });
    }
});

// Get keyword creator
app.get('/blacklist/:index/creator', async (req, res) => {
    try {
        const index = parseInt(req.params.index);
        const keywordData = await contract.getKeyword(index);

        res.json({ creator: keywordData.username });
    } catch (err) {
        console.error('Error getting keyword creator:', err);
        res.status(500).json({ error: err.message });
    }
});

// Get statistics
app.get('/blacklist/statistics', async (req, res) => {
    try {
        const totalKeywords = await contract.getKeywordCount();
        const activeKeywords = await contract.getActiveKeywordCount();

        res.json({
            totalBlacklistEntries: totalKeywords.toString(),
            activeKeywords: activeKeywords.toString(),
            totalApply: "0", // This would be tracked separately if needed
            totalVerify: "0"  // This would be tracked separately if needed
        });
    } catch (err) {
        console.error('Error getting statistics:', err);
        res.status(500).json({ error: err.message });
    }
});

// Health check
app.get('/health', (req, res) => {
    res.json({
        status: 'healthy',
        contractAddress: contractAddress,
        walletAddress: wallet.address,
        network: provider.network
    });
});

app.listen(PORT, () => {
    console.log(`Blockchain API server running on port ${PORT}`);
    console.log(`Contract address: ${contractAddress}`);
    console.log(`Wallet address: ${wallet.address}`);
}); 