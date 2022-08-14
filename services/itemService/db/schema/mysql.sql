CREATE DATABASE IF NOT EXISTS itemservicedb;

CREATE USER IF NOT EXISTS 'entrytask'@'localhost' IDENTIFIED BY 'password';

GRANT CREATE, ALTER, DROP, INSERT, UPDATE, DELETE, SELECT, REFERENCES, RELOAD ON itemservicedb.* TO 'entrytask'@'localhost';


CREATE USER IF NOT EXISTS 'exporter'@'%' IDENTIFIED BY 'exporterpassword' WITH MAX_USER_CONNECTIONS 3;
GRANT PROCESS, REPLICATION CLIENT, SELECT ON *.* TO 'exporter'@'%';

USE itemservicedb;
DROP TABLE IF EXISTS Favourites;
CREATE TABLE Favourites (
    id bigint unsigned AUTO_INCREMENT PRIMARY KEY,
    userID bigint unsigned,
    itemID bigint NOT NULL,
    shopID bigint NOT NULL,
    timeAdded TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(userID, shopID, itemID)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
 
CREATE INDEX userId_timeAdded_idx ON Favourites(userID, timeAdded);
CREATE INDEX userId_item_idx ON Favourites(userID, itemID, shopID);