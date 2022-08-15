CREATE DATABASE IF NOT EXISTS userservicedb;

CREATE USER IF NOT EXISTS 'entrytask'@'localhost' IDENTIFIED BY 'password';
GRANT CREATE, ALTER, DROP, INSERT, UPDATE, DELETE, SELECT, REFERENCES, RELOAD on userservicedb.* TO 'entrytask'@'localhost';

CREATE USER IF NOT EXISTS 'exporter'@'%' IDENTIFIED BY 'exporterpassword' WITH MAX_USER_CONNECTIONS 3;
GRANT PROCESS, REPLICATION CLIENT, SELECT ON *.* TO 'exporter'@'%';

USE userservicedb;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    userID bigint unsigned AUTO_INCREMENT PRIMARY KEY, 
    username varchar(15) UNIQUE NOT NULL,
    password binary(60) NOT NULL
 ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX username_pwd_idx ON users(username, password);
