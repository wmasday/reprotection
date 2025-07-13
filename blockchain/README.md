# Blockchain Blacklist Keyword System

This is a Solidity smart contract system for managing blacklist keywords on the blockchain.

## Features

- Add keywords to blockchain
- Remove keywords (only by creator)
- View all keywords
- Check keyword status
- Get statistics

## Setup

### 1. Install Dependencies
```bash
npm install
```

### 2. Compile Contracts
```bash
npm run compile
```

### 3. Start Local Blockchain
```bash
npm run node
```

### 4. Deploy Contract
In a new terminal:
```bash
npm run deploy
```

### 5. Start API Server
```bash
npm run api
```

## Environment Variables

Create a `.env` file in the blockchain directory:

```env
# Blockchain Configuration
BLOCKCHAIN_RPC_URL=http://127.0.0.1:8545
BLOCKCHAIN_API_PORT=3001

# Private key for the wallet (default Hardhat account)
PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

# Network Configuration
NETWORK_NAME=localhost
CHAIN_ID=1337
```

## API Endpoints

- `GET /blacklist` - Get all keywords
- `GET /blacklist/:keyword` - Get specific keyword
- `POST /blacklist/add` - Add keyword
- `POST /blacklist/remove` - Remove keyword
- `GET /blacklist/:keyword/isActive` - Check if keyword is active
- `GET /blacklist/:keyword/creator` - Get keyword creator
- `GET /blacklist/statistics` - Get statistics
- `GET /health` - Health check

## Smart Contract Functions

- `addKeyword(string keyword, string username)` - Add keyword
- `removeKeyword(string keyword)` - Remove keyword
- `getKeyword(string keyword)` - Get keyword data
- `isKeywordActive(string keyword)` - Check if active
- `getAllKeywords()` - Get all keywords
- `getKeywordCount()` - Get total count
- `getActiveKeywordCount()` - Get active count

## Integration with Go App

The Go application connects to this blockchain API via the `BLOCKCHAIN_API_URL` environment variable. 