package db

import (
	"database/sql"
	"fmt"

	"itemService/config"
	"itemService/constants"
	metrics "itemService/metrics"

	"github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type DbManager struct {
	db     *sql.DB
	config *config.DbConfig
	logger *zap.Logger
}

const (
	queryTypeInsert = "INSERT"
	queryTypeSelect = "SELECT"
	queryTypeDelete = "DELETE"
)

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

	return &dbManager, err
}

func (dm *DbManager) QueryOne(query string, opName string, destination ...any) error {
	successStr := trueStr
	// time database query
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.DatabaseOpDuration.WithLabelValues(dm.config.ServiceLabel, queryTypeSelect, opName, successStr).Observe(v)
	}))
	defer func() {
		// observe duration at the end of this function
		timer.ObserveDuration()
	}()

	res := dm.db.QueryRow(query)
	err := res.Scan(destination...)
	if err != nil {
		dm.logger.Error(
			constants.ERROR_DATABASE_QUERY_MSG,
			zap.String("query", query),
			zap.String("opName", opName),
			zap.Error(err),
		)
		if err != sql.ErrNoRows {
			// avoid false negatives, a select query can return no rows
			successStr = falseStr
		}
		return err
	}
	dm.logger.Info(
		constants.INFO_DATABASE_QUERY,
		zap.String("query", query),
	)
	return err
}

func (dm *DbManager) QueryRows(query string, opName string) (*sql.Rows, error) {
	successStr := trueStr
	// time database query
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.DatabaseOpDuration.WithLabelValues(dm.config.ServiceLabel, queryTypeInsert, opName, successStr).Observe(v)
	}))
	defer func() {
		// observe duration at the end of this function
		timer.ObserveDuration()
	}()
	rows, err := dm.db.Query(query)
	dm.logger.Info(
		constants.INFO_DATABASE_QUERY_ROWS,
		zap.String("query", query),
		zap.Error(err),
	)
	return rows, err
}

func (dm *DbManager) InsertRow(query string, opName string) (int64, error) {
	successStr := trueStr
	// time database query
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.DatabaseOpDuration.WithLabelValues(dm.config.ServiceLabel, queryTypeInsert, opName, successStr).Observe(v)
	}))
	defer func() {
		// observe duration at the end of this function
		timer.ObserveDuration()
	}()

	res, err := dm.db.Exec(query)

	if err != nil {
		dm.logger.Error(
			constants.ERROR_DATABASE_INSERT_MSG,
			zap.String("query", query),
			zap.String("opName", opName),
			zap.Error(err),
		)
		successStr = falseStr
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		successStr = falseStr
		return 0, err
	}

	dm.logger.Info(
		constants.INFO_DATABASE_INSERT,
		zap.String("query", query),
		zap.Any("id", id),
	)

	return id, err
}

func (dm *DbManager) DeleteOne(query string, opName string) (int64, error) {
	successStr := trueStr
	// time database query
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.DatabaseOpDuration.WithLabelValues(dm.config.ServiceLabel, queryTypeDelete, opName, successStr).Observe(v)
	}))
	defer func() {
		// observe duration at the end of this function
		timer.ObserveDuration()
	}()

	res, err := dm.db.Exec(query)
	dm.logger.Info(
		constants.INFO_DATABASE_DELETE,
		zap.String("query", query),
		zap.Any("res", res),
	)

	if err != nil {
		dm.logger.Error(
			constants.ERROR_DATABASE_DELETE_MSG,
			zap.String("query", query),
			zap.String("opName", opName),
			zap.Error(err),
		)
		successStr = falseStr
		return 0, err
	}

	return res.RowsAffected()
}
