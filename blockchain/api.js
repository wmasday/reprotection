require('dotenv').config();
const express = require('express');
const { ethers } = require('ethers');
const fs = require('fs');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 3001;

// Update this address after deployment
envBlacklistManagerAddress = process.env.BLACKLIST_MANAGER_ADDRESS || '';
const BLACKLIST_MANAGER_ADDRESS = envBlacklistManagerAddress;

const ABI_PATH = path.join(__dirname, 'artifacts', 'contracts', 'BlacklistManager.sol', 'BlacklistManager.json');
const abi = JSON.parse(fs.readFileSync(ABI_PATH)).abi;

// Connect to local Hardhat node
const provider = new ethers.JsonRpcProvider('http://127.0.0.1:8545');
const contract = new ethers.Contract(BLACKLIST_MANAGER_ADDRESS, abi, provider);

app.use(express.json());

// Dummy list of keywords for demo; in production, use an off-chain indexer or event logs
const knownKeywords = [
    "phishing",
    "scam",
    "malware",
    "ransomware",
    "fraud"
];

// In-memory storage for demo purposes
let blacklistData = new Map();
let statistics = {
    totalApply: 0,
    totalVerify: 0,
    totalBlacklistEntries: 0
};

// Get all blacklist entries (demo: using knownKeywords)
app.get('/blacklist', async (req, res) => {
    try {
        const results = [];
        for (const keyword of knownKeywords) {
            try {
                const entry = await contract.getBlacklistEntry(keyword);
                if (entry.isActive) {
                    results.push({
                        keyword: entry.keywordName,
                        createdBy: entry.createdBy,
                        username: entry.username,
                        timestamp: entry.timestamp.toString(),
                        isActive: entry.isActive
                    });
                }
            } catch (err) {
                // skip if not found
            }
        }

        // Add demo data from in-memory storage
        for (const [keyword, entry] of blacklistData) {
            if (entry.isActive) {
                results.push({
                    keyword: entry.keyword,
                    createdBy: entry.createdBy,
                    username: entry.username,
                    timestamp: entry.timestamp.toString(),
                    isActive: entry.isActive
                });
            }
        }

        res.json(results);
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

// Get a blacklist entry by keyword
app.get('/blacklist/:keyword', async (req, res) => {
    try {
        const keyword = req.params.keyword;
        const entry = await contract.getBlacklistEntry(keyword);
        res.json({
            keyword: entry.keywordName,
            createdBy: entry.createdBy,
            username: entry.username,
            timestamp: entry.timestamp.toString(),
            isActive: entry.isActive
        });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

// Check if a keyword is blacklisted
app.get('/blacklist/:keyword/isBlacklisted', async (req, res) => {
    try {
        const keyword = req.params.keyword;
        const isBlacklisted = await contract.isKeywordBlacklisted(keyword);
        res.json({ keyword, isBlacklisted });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

// Get statistics
app.get('/blacklist/statistics', async (req, res) => {
    try {
        const stats = await contract.getStatistics();
        res.json({
            totalApply: (parseInt(stats.totalApply.toString()) + statistics.totalApply).toString(),
            totalVerify: (parseInt(stats.totalVerify.toString()) + statistics.totalVerify).toString(),
            totalBlacklistEntries: (parseInt(stats.totalBlacklistEntries.toString()) + statistics.totalBlacklistEntries).toString()
        });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

// Add blacklist keyword (demo endpoint)
app.post('/blacklist/add', async (req, res) => {
    try {
        const { keyword, username } = req.body;

        if (!keyword || !username) {
            return res.status(400).json({ error: 'Keyword and username are required' });
        }

        if (keyword.length > 100) {
            return res.status(400).json({ error: 'Keyword too long (max 100 characters)' });
        }

        if (blacklistData.has(keyword) && blacklistData.get(keyword).isActive) {
            return res.status(400).json({ error: 'Keyword already in blacklist' });
        }

        const entry = {
            keyword: keyword,
            createdBy: '0x' + Math.random().toString(16).substr(2, 40),
            username: username,
            timestamp: Math.floor(Date.now() / 1000),
            isActive: true
        };

        blacklistData.set(keyword, entry);
        statistics.totalBlacklistEntries++;

        res.json({ success: true, message: 'Keyword added successfully' });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

// Remove blacklist keyword (demo endpoint)
app.post('/blacklist/remove', async (req, res) => {
    try {
        const { keyword, username } = req.body;

        if (!keyword || !username) {
            return res.status(400).json({ error: 'Keyword and username are required' });
        }

        if (!blacklistData.has(keyword) || !blacklistData.get(keyword).isActive) {
            return res.status(400).json({ error: 'Keyword not in blacklist' });
        }

        const entry = blacklistData.get(keyword);
        if (entry.username !== username) {
            return res.status(403).json({ error: 'Only the creator can remove this keyword' });
        }

        blacklistData.get(keyword).isActive = false;
        statistics.totalBlacklistEntries--;

        res.json({ success: true, message: 'Keyword removed successfully' });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

// Increment apply counter
app.post('/blacklist/apply', async (req, res) => {
    try {
        statistics.totalApply++;
        res.json({ success: true, totalApply: statistics.totalApply });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

// Get keyword creator
app.get('/blacklist/:keyword/creator', async (req, res) => {
    try {
        const keyword = req.params.keyword;
        
        if (!blacklistData.has(keyword)) {
            return res.status(404).json({ error: 'Keyword not found' });
        }
        
        const entry = blacklistData.get(keyword);
        res.json({ creator: entry.username });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

// Increment verify counter
app.post('/blacklist/verify', async (req, res) => {
    try {
        statistics.totalVerify++;
        res.json({ success: true, totalVerify: statistics.totalVerify });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

app.listen(PORT, () => {
    console.log(`API server running on port ${PORT}`);
}); 