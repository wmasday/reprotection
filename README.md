<!-- Centered Hero Icon -->
<p align="center">
  <img src="static/icon.svg" alt="Hero Icon" width="180" />
</p>

# Reprotection Dashboard

> **A Modern, Secure, and Powerful Blacklist Keyword Management Platform with Blockchain Integration**

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)](https://golang.org/) [![Solidity](https://img.shields.io/badge/Solidity-0.8.19-363636?logo=solidity)](https://soliditylang.org/) [![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

---

## 🚀 Overview

**Reprotection** is a comprehensive dashboard for managing and monitoring blacklist keywords with blockchain integration. It combines traditional keyword management with decentralized blockchain storage, providing enhanced security, transparency, and collaboration features. Designed for security teams and administrators, it offers both local keyword management and blockchain-based keyword sharing.

---

## ✨ Features

### 🔗 Blockchain Integration
- **Decentralized Storage**: Keywords stored on Ethereum blockchain via smart contracts
- **Duplicate Keywords**: Allow multiple users to add the same keyword with individual tracking
- **Creator Tracking**: Each keyword entry tracks its creator and timestamp
- **Active/Inactive Toggle**: Users can activate/deactivate their own keywords
- **Apply to Items**: Direct integration between blockchain keywords and local items system

### 🔍 Local Management
- **Search & Filter**: Instantly search and filter blacklist keywords
- **Keyword Analytics**: Visualize keyword distribution with modern charts
- **Malicious File Tracking**: See which files are flagged as malicious per keyword
- **Add/Remove Keywords**: Manage your local blacklist with ease

### 🛡️ Security & Access
- **User Authentication**: Secure access with session management
- **Creator Permissions**: Only keyword creators can modify their entries
- **Project Configuration**: Set and update your working project path
- **Modular Structure**: Clean, maintainable Go backend and Bootstrap frontend

---

## 🏗️ Architecture

### System Components

```text
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Go Web App    │    │  Blockchain API │    │  Smart Contract │
│   (Port 1337)   │◄──►│   (Port 3001)   │◄──►│  (Ethereum)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   SQLite DB     │    │   Node.js/      │    │   Hardhat       │
│   (Local Items) │    │   Express API   │    │   Development   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Smart Contract Features
- **BlacklistManager.sol**: Manages keyword storage on blockchain
- **Add Keywords**: Users can add keywords with username tracking
- **Toggle Status**: Activate/deactivate keywords
- **Creator Control**: Only creators can modify their keywords
- **Duplicate Support**: Multiple users can add same keyword

---

## 🛠️ Tech Stack

### Backend
- **Go (Golang)**: Main web application
- **SQLite**: Local database for items and users
- **GORM**: Database ORM
- **Gin**: Web framework (if used)

### Blockchain
- **Solidity**: Smart contract development
- **Hardhat**: Ethereum development environment
- **Node.js/Express**: Blockchain API server
- **Ethers.js**: Ethereum interaction

### Frontend
- **Bootstrap 5**: UI framework
- **Chart.js**: Data visualization
- **Font Awesome**: Icons
- **JavaScript**: Dynamic interactions

---

## 📦 Project Structure

```text
reprotection/
├── blockchain/              # Blockchain components
│   ├── contracts/          # Solidity smart contracts
│   ├── api.js             # Express API server
│   ├── hardhat.config.js  # Hardhat configuration
│   └── scripts/           # Deployment scripts
├── cmd/                   # CLI tools (user creation, migrations)
├── config/                # Configuration files (DB, session, migration)
├── controllers/           # HTTP controllers (auth, config, item, blockchain)
├── main.go               # Main entry point
├── middleware/            # Authentication middleware
├── migrations/            # SQL migration scripts
├── models/               # Data models (User, Item)
├── static/               # Static assets (CSS, JS)
├── views/                # HTML templates
└── README.md             # This file
```

---

## ⚡ Getting Started

### Prerequisites

- **Go 1.20+**: [Download here](https://golang.org/)
- **Node.js 16+**: [Download here](https://nodejs.org/)
- **Git**: [Download here](https://git-scm.com/)

### 1. Clone the Repository

```bash
git clone https://github.com/wmasday/reprotection.git
cd reprotection
```

### 2. Install Go Dependencies

```bash
go mod tidy
```

### 3. Setup Blockchain Components

```bash
# Navigate to blockchain directory
cd blockchain

# Install Node.js dependencies
npm install

# Compile smart contracts
npx hardhat compile

# Deploy smart contract (requires local Ethereum node or testnet)
npx hardhat run scripts/deploy.js --network localhost

# Start blockchain API server
node api.js
```

### 4. Setup Database

```bash
# Run database migrations
go run cmd/migrate/main.go

# Create admin user
go run cmd/create_user/main.go
```

### 5. Configure Environment

Create a `.env` file in the root directory:

```env
BLOCKCHAIN_API_URL=http://localhost:3001
BLOCKCHAIN_CONTRACT_ADDRESS=0x... # Your deployed contract address
```

### 6. Start the Application

```bash
# Start the Go web application
go run main.go
```

The dashboard will be available at [http://localhost:1337](http://localhost:1337)

---

## 🔗 Blockchain Integration

### Smart Contract Features

The `BlacklistManager.sol` contract provides:

- **Add Keywords**: `addKeyword(string keyword, string username)`
- **Toggle Status**: `toggleKeyword(uint256 index)`
- **Get All Keywords**: `getAllKeywords()`
- **Creator Verification**: Only creators can modify their keywords

### API Endpoints

The blockchain API (`api.js`) provides:

- `GET /blacklist` - Get all keywords
- `POST /blacklist/add` - Add new keyword
- `POST /blacklist/toggle` - Toggle keyword status
- `GET /blacklist/statistics` - Get statistics

### Integration Flow

1. **User adds keyword** → Blockchain API → Smart Contract
2. **User applies keyword** → Go App → Items System
3. **User toggles keyword** → Blockchain API → Smart Contract
4. **Display keywords** → Blockchain API → Go App → UI

---

## 🖼️ Screenshots

### Dashboard
![Dashboard Screenshot](static/hero.png)

### Blockchain Keywords
- **Keyword Management**: Add, toggle, and apply blockchain keywords
- **Creator Tracking**: See who created each keyword
- **Duplicate Support**: Multiple users can add same keyword
- **Real-time Updates**: Dynamic timestamp formatting

---

## 🔧 Configuration

### Environment Variables

```env
# Blockchain Configuration
BLOCKCHAIN_API_URL=http://localhost:3001
BLOCKCHAIN_CONTRACT_ADDRESS=0x...

# Database Configuration
DB_PATH=./database.db

# Server Configuration
PORT=1337
```

### Database Schema

```sql
-- Users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

-- Items table
CREATE TABLE items (
    id INTEGER PRIMARY KEY,
    title VARCHAR(255) NOT NULL
);

-- Malicious files table
CREATE TABLE malicious (
    id INTEGER PRIMARY KEY,
    item_id INTEGER NOT NULL,
    filepath VARCHAR(255) NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items(id)
);
```

---

## 🚀 Deployment

### Local Development

```bash
# Terminal 1: Start blockchain API
cd blockchain && node api.js

# Terminal 2: Start Go application
go run main.go
```

### Production Deployment

1. **Deploy Smart Contract** to mainnet/testnet
2. **Update Contract Address** in environment variables
3. **Deploy Go Application** to your server
4. **Configure Reverse Proxy** (nginx/Apache)
5. **Setup SSL Certificate** for HTTPS

---

## 🤝 Contributing

Contributions are welcome! Please open issues and submit pull requests for new features, bug fixes, or improvements.

### Development Setup

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Make your changes
4. Test thoroughly (both Go app and blockchain components)
5. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
6. Push to the branch (`git push origin feature/AmazingFeature`)
7. Open a Pull Request

### Code Style

- **Go**: Follow standard Go formatting (`gofmt`)
- **JavaScript**: Use ESLint with standard configuration
- **Solidity**: Follow Solidity style guide
- **HTML/CSS**: Use Bootstrap classes and maintain consistency

---

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## 📬 Contact

- **Author:** Relyze Team by Recov
- **Email:** [withmasday@gmail.com](mailto:withmasdayl@gmail.com)
- **GitHub:** [wmasday](https://github.com/wmasday)

---

## 🙏 Acknowledgments

- **Bootstrap**: For the beautiful UI framework
- **Chart.js**: For data visualization
- **Hardhat**: For Ethereum development tools
- **Font Awesome**: For the icon library

---

> **Reprotection** – Secure your projects, empower your team, embrace blockchain innovation.
