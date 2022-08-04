package db

import (
	"database/sql"
	"fmt"
	"log"
	"userService/config"

	"github.com/go-sql-driver/mysql"
)

type DbManager struct {
	db     *sql.DB
	config *config.DbConfig
}

func Init(dbConfig *config.DbConfig) (*DbManager, error) {
	cfg := mysql.Config{
		User:   dbConfig.User,
		Passwd: dbConfig.Password,
		Net:    dbConfig.Net,
		Addr:   fmt.Sprintf("%s:%s", dbConfig.Host, dbConfig.Port),
		DBName: dbConfig.DbName,
	}

	// get database handle
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalln("Error occured when opening database: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Println("Connected to mysql database!")

	dbManager := DbManager{
		db:     db,
		config: dbConfig,
	}

	return &dbManager, nil
}

func (dm *DbManager) QueryOne(query string) *sql.Row {
	res := dm.db.QueryRow(query)
	return res
}

func (dm *DbManager) InsertRow(query string) (int64, error) {
	res, err := dm.db.Exec(query)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
