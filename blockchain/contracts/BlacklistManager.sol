// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract BlacklistManager {
    struct Keyword {
        string keyword;
        address creator;
        string username;
        uint256 timestamp;
        bool isActive;
    }
    
    Keyword[] public keywords;
    
    event KeywordAdded(string keyword, address creator, string username, uint256 timestamp);
    event KeywordRemoved(string keyword, address creator, uint256 timestamp);
    
    modifier keywordExists(uint256 _index) {
        require(_index < keywords.length, "Keyword does not exist");
        _;
    }
    
    modifier onlyKeywordCreator(uint256 _index) {
        require(keywords[_index].creator == msg.sender, "Only creator can modify keyword");
        _;
    }
    
    function addKeyword(string memory _keyword, string memory _username) public {
        require(bytes(_keyword).length > 0, "Keyword cannot be empty");
        require(bytes(_keyword).length <= 100, "Keyword too long");
        
        keywords.push(Keyword({
            keyword: _keyword,
            creator: msg.sender,
            username: _username,
            timestamp: block.timestamp,
            isActive: true
        }));
        
        emit KeywordAdded(_keyword, msg.sender, _username, block.timestamp);
    }
    
    function removeKeyword(uint256 _index) public keywordExists(_index) onlyKeywordCreator(_index) {
        require(keywords[_index].isActive, "Keyword is already inactive");
        
        keywords[_index].isActive = false;
        
        emit KeywordRemoved(keywords[_index].keyword, msg.sender, block.timestamp);
    }
    
    function toggleKeyword(uint256 _index) public keywordExists(_index) onlyKeywordCreator(_index) {
        keywords[_index].isActive = !keywords[_index].isActive;
        
        if (keywords[_index].isActive) {
            emit KeywordAdded(keywords[_index].keyword, msg.sender, keywords[_index].username, block.timestamp);
        } else {
            emit KeywordRemoved(keywords[_index].keyword, msg.sender, block.timestamp);
        }
    }
    
    function getKeyword(uint256 _index) public view returns (
        string memory keyword,
        address creator,
        string memory username,
        uint256 timestamp,
        bool isActive
    ) {
        require(_index < keywords.length, "Keyword does not exist");
        Keyword memory k = keywords[_index];
        return (k.keyword, k.creator, k.username, k.timestamp, k.isActive);
    }
    
    function isKeywordActive(uint256 _index) public view returns (bool) {
        require(_index < keywords.length, "Keyword does not exist");
        return keywords[_index].isActive;
    }
    
    function getAllKeywords() public view returns (Keyword[] memory) {
        return keywords;
    }
    
    function getKeywordCount() public view returns (uint256) {
        return keywords.length;
    }
    
    function getActiveKeywordCount() public view returns (uint256) {
        uint256 count = 0;
        for (uint256 i = 0; i < keywords.length; i++) {
            if (keywords[i].isActive) {
                count++;
            }
        }
        return count;
    }
} 