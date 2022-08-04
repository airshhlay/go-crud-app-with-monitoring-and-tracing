USE userservicedb;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    userId bigint unsigned AUTO_INCREMENT PRIMARY KEY, 
    username varchar(15) UNIQUE NOT NULL,
    password binary(60) NOT NULL
 ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX username_pwd_idx ON users(username, password);
