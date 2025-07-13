// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract BlacklistManager {
    struct BlacklistEntry {
        string keyword;
        address createdBy;
        string username;
        uint256 timestamp;
        bool isActive;
    }
    
    struct Statistics {
        uint256 totalApply;
        uint256 totalVerify;
        uint256 totalBlacklistEntries;
    }
    
    mapping(string => BlacklistEntry) public blacklistEntries;
    mapping(address => string) public userUsernames;
    mapping(address => uint256) public userCreatedEntries;
    
    Statistics public stats;
    
    event BlacklistKeywordAdded(string keyword, address createdBy, string username, uint256 timestamp);
    event BlacklistKeywordRemoved(string keyword, address removedBy, uint256 timestamp);
    event UsernameSet(address user, string username);
    event ApplyIncremented(address user, uint256 newTotal);
    event VerifyIncremented(address user, uint256 newTotal);
    
    modifier onlyValidKeyword(string memory keyword) {
        require(bytes(keyword).length > 0, "Keyword cannot be empty");
        require(bytes(keyword).length <= 100, "Keyword too long");
        _;
    }
    
    modifier onlyValidUsername(string memory username) {
        require(bytes(username).length > 0, "Username cannot be empty");
        require(bytes(username).length <= 50, "Username too long");
        _;
    }
    
    function addBlacklistKeyword(string memory keyword, string memory username) 
        public 
        onlyValidKeyword(keyword) 
        onlyValidUsername(username) 
    {
        require(!blacklistEntries[keyword].isActive, "Keyword already in blacklist");
        
        // Set username for the user if not already set
        if (bytes(userUsernames[msg.sender]).length == 0) {
            userUsernames[msg.sender] = username;
            emit UsernameSet(msg.sender, username);
        }
        
        blacklistEntries[keyword] = BlacklistEntry({
            keyword: keyword,
            createdBy: msg.sender,
            username: userUsernames[msg.sender],
            timestamp: block.timestamp,
            isActive: true
        });
        
        userCreatedEntries[msg.sender]++;
        stats.totalBlacklistEntries++;
        
        emit BlacklistKeywordAdded(keyword, msg.sender, userUsernames[msg.sender], block.timestamp);
    }
    
    function removeBlacklistKeyword(string memory keyword) public {
        require(blacklistEntries[keyword].isActive, "Keyword not in blacklist");
        require(blacklistEntries[keyword].createdBy == msg.sender, "Only creator can remove keyword");
        
        blacklistEntries[keyword].isActive = false;
        stats.totalBlacklistEntries--;
        
        emit BlacklistKeywordRemoved(keyword, msg.sender, block.timestamp);
    }
    
    function incrementApply() public {
        stats.totalApply++;
        emit ApplyIncremented(msg.sender, stats.totalApply);
    }
    
    function incrementVerify() public {
        stats.totalVerify++;
        emit VerifyIncremented(msg.sender, stats.totalVerify);
    }
    
    function setUsername(string memory username) public onlyValidUsername(username) {
        userUsernames[msg.sender] = username;
        emit UsernameSet(msg.sender, username);
    }
    
    function getBlacklistEntry(string memory keyword) public view returns (
        string memory keywordName,
        address createdBy,
        string memory username,
        uint256 timestamp,
        bool isActive
    ) {
        BlacklistEntry memory entry = blacklistEntries[keyword];
        return (
            entry.keyword,
            entry.createdBy,
            entry.username,
            entry.timestamp,
            entry.isActive
        );
    }
    
    function isKeywordBlacklisted(string memory keyword) public view returns (bool) {
        return blacklistEntries[keyword].isActive;
    }
    
    function getStatistics() public view returns (
        uint256 totalApply,
        uint256 totalVerify,
        uint256 totalBlacklistEntries
    ) {
        return (
            stats.totalApply,
            stats.totalVerify,
            stats.totalBlacklistEntries
        );
    }
    
    function getUserInfo(address user) public view returns (
        string memory username,
        uint256 createdEntries
    ) {
        return (
            userUsernames[user],
            userCreatedEntries[user]
        );
    }
    
    function getUsername(address user) public view returns (string memory) {
        return userUsernames[user];
    }
} 