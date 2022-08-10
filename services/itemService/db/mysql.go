package db

import (
	"database/sql"
	"fmt"

	"itemService/config"
	"itemService/constants"

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
		ParseTime:            true,
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

func (dm *DbManager) QueryRows(query string) (*sql.Rows, error) {
	rows, err := dm.db.Query(query)
	dm.logger.Info(
		constants.INFO_DATABASE_QUERY_ROWS,
		zap.String("query", query),
		zap.Error(err),
	)
	return rows, err
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

func (dm *DbManager) DeleteOne(query string) (int64, error) {
	res, err := dm.db.Exec(query)
	dm.logger.Info(
		constants.INFO_DATABASE_DELETE,
		zap.String("query", query),
		zap.Any("res", res),
	)

	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
