#!/bin/bash

echo "🚀 Setting up Blockchain Blacklist Manager..."

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js first."
    exit 1
fi

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed. Please install npm first."
    exit 1
fi

# Navigate to blockchain directory
cd blockchain

# Install dependencies
echo "📦 Installing dependencies..."
npm install

# Check if .env file exists, if not create one
if [ ! -f .env ]; then
    echo "📝 Creating .env file..."
    cat > .env << EOF
PORT=3001
BLACKLIST_MANAGER_ADDRESS=
EOF
    echo "✅ Created .env file. Please update BLACKLIST_MANAGER_ADDRESS with your deployed contract address."
fi

# Check if Hardhat is installed globally
if ! command -v npx hardhat &> /dev/null; then
    echo "📦 Installing Hardhat..."
    npm install -g hardhat
fi

echo "✅ Blockchain setup complete!"
echo ""
echo "To start the blockchain API server:"
echo "1. cd blockchain"
echo "2. npm start"
echo ""
echo "To deploy the smart contract:"
echo "1. cd blockchain"
echo "2. npx hardhat node (in one terminal)"
echo "3. npx hardhat run scripts/deploy.js --network localhost (in another terminal)"
echo ""
echo "Then update the BLACKLIST_MANAGER_ADDRESS in .env with the deployed contract address." 