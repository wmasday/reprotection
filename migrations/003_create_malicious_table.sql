CREATE TABLE IF NOT EXISTS malicious (
    id INT AUTO_INCREMENT PRIMARY KEY,
    item_id INT NOT NULL,
    filepath VARCHAR(255) NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items(id)
); 