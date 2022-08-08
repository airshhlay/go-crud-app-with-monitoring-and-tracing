USE itemservicedb;
DROP TABLE IF EXISTS Favourites;
CREATE TABLE Favourites (
    id bigint unsigned AUTO_INCREMENT PRIMARY KEY,
    userId bigint unsigned,
    itemId bigint NOT NULL,
    shopId bigint NOT NULL,
    timeAdded TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(userId, shopId, itemId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
 
CREATE INDEX userId_timeAdded_idx ON Favourites(userId, timeAdded);
CREATE INDEX userId_item_idx ON Favourites(userId, itemId, shopId);