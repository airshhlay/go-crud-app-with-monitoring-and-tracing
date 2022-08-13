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

// DatabaseManager is a database manager struct containing a reference to the database connection, zap logger, and the database config
type DatabaseManager struct {
	db     *sql.DB
	config *config.DbConfig
	logger *zap.Logger
}

const (
	queryTypeInsert = "INSERT"
	queryTypeSelect = "SELECT"
	queryTypeDelete = "DELETE"
)

// InitDatabase opens the database connection. It returns an error if the database fails to respond when pinged.
func InitDatabase(dbConfig *config.DbConfig, logger *zap.Logger) (*DatabaseManager, error) {
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
			constants.ErrorDatabaseConnectionMsg,
			zap.Int32("errorCode", constants.ErrorDatabaseConnection),
			zap.Error(err))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal(
			constants.ErrorDatabaseConnectionMsg,
			zap.Int32("errorCode", constants.ErrorDatabaseConnection),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info(constants.InfoDatabaseConnectSuccess)

	dbManager := DatabaseManager{
		db:     db,
		config: dbConfig,
		logger: logger,
	}

	return &dbManager, err
}

// QueryOne will query for a single *sql.Row, and write its contents into destination.
func (dm *DatabaseManager) QueryOne(query string, opName string, destination ...any) error {
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
			constants.ErrorDatabaseQueryMsg,
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
		constants.InfoDatabaseQuery,
		zap.String("query", query),
	)
	return err
}

// QueryRows executes the given query and returns the queried rows.
func (dm *DatabaseManager) QueryRows(query string, opName string) (*sql.Rows, error) {
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
		constants.InfoDatabaseQueryRows,
		zap.String("query", query),
		zap.Error(err),
	)
	return rows, err
}

// InsertRow will insert a single row and return its ID.
func (dm *DatabaseManager) InsertRow(query string, opName string) (int64, error) {
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
			constants.ErrorDatabaseInsertMsg,
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
		constants.InfoDatabaseInsert,
		zap.String("query", query),
		zap.Any("id", id),
	)

	return id, err
}

// DeleteOne deletes a row from the database and returns the number of rows deleted
func (dm *DatabaseManager) DeleteOne(query string, opName string) (int64, error) {
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
		constants.InfoDatabaseDelete,
		zap.String("query", query),
		zap.Any("res", res),
	)

	if err != nil {
		dm.logger.Error(
			constants.ErrorDatabaseDeleteMsg,
			zap.String("query", query),
			zap.String("opName", opName),
			zap.Error(err),
		)
		successStr = falseStr
		return 0, err
	}

	return res.RowsAffected()
}
