package db

import (
	"database/sql"
	"fmt"
	config "userService/config"
	constants "userService/constants"

	"go.uber.org/zap"

	"github.com/go-sql-driver/mysql"
)

type DbManager struct {
	db     *sql.DB
	config *config.DbConfig
	logger *zap.Logger
}

func InitDatabase(dbConfig *config.DbConfig, logger *zap.Logger) (*DbManager, error) {
	cfg := mysql.Config{
		User:                 dbConfig.User,
		Passwd:               dbConfig.Password,
		Net:                  dbConfig.Net,
		Addr:                 fmt.Sprintf("%s:%s", dbConfig.Host, dbConfig.Port),
		DBName:               dbConfig.DbName,
		AllowNativePasswords: true,
	}

	// get database handle
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		logger.Fatal(
			constants.ERROR_DATABASE_CONNECTION_MSG,
			zap.Int32("errorCode", constants.ERROR_DATABASE_CONNECTION),
			zap.Error(err))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal(
			constants.ERROR_DATABASE_CONNECTION_MSG,
			zap.Int32("errorCode", constants.ERROR_DATABASE_CONNECTION),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info(constants.INFO_DATABASE_CONNECT_SUCCESS)

	dbManager := DbManager{
		db:     db,
		config: dbConfig,
		logger: logger,
	}

	return &dbManager, nil
}

func (dm *DbManager) QueryOne(query string) *sql.Row {
	res := dm.db.QueryRow(query)
	dm.logger.Info(
		constants.INFO_DATABASE_QUERY,
		zap.String("query", query),
	)
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

	dm.logger.Info(
		constants.INFO_DATABASE_INSERT,
		zap.String("query", query),
		zap.Any("id", id),
	)

	return id, nil
}

func (dm *DbManager) QueryMany(query string) (*sql.Rows, error) {
	res, err := dm.db.Query(query)
	dm.logger.Info(
		constants.INFO_DATABASE_QUERY,
		zap.String("query", query),
	)
	return res, err
}
