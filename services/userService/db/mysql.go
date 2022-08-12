package db

import (
	"database/sql"
	"fmt"
	config "userService/config"
	constants "userService/constants"
	metrics "userService/metrics"

	"github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const (
	queryTypeInsert = "INSERT"
	queryTypeSelect = "SELECT"
	queryTypeDelete = "DELETE"
	falseStr        = "false"
	trueStr         = "true"
)

// DatabaseManager is a database manager struct containing a reference to the database connection, zap logger, and the database config
type DatabaseManager struct {
	db     *sql.DB
	config *config.DbConfig
	logger *zap.Logger
}

// InitDatabase opens the database connection. It returns an error if the database fails to respond when pinged.
func InitDatabase(dbConfig *config.DbConfig, logger *zap.Logger) (*DatabaseManager, error) {
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

	logger.Info(constants.INFO_DATABASE_CONNECT_SUCCESS)

	dbManager := DatabaseManager{
		db:     db,
		config: dbConfig,
		logger: logger,
	}

	return &dbManager, nil
}

// QueryOne will return a single row
// func (dm *DatabaseManager) QueryOne(query string) *sql.Row {
// 	res := dm.db.QueryRow(query)
// 	dm.logger.Info(
// 		constants.INFO_DATABASE_QUERY,
// 		zap.String("query", query),
// 	)
// 	return res
// }
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
		constants.INFO_DATABASE_QUERY,
		zap.String("query", query),
	)
	return err
}

// InsertRow will insert a single row and return its ID
// func (dm *DatabaseManager) InsertRow(query string) (int64, error) {
// 	res, err := dm.db.Exec(query)

// 	if err != nil {
// 		return 0, err
// 	}

// 	id, err := res.LastInsertId()
// 	if err != nil {
// 		return 0, err
// 	}

// 	dm.logger.Info(
// 		constants.INFO_DATABASE_INSERT,
// 		zap.String("query", query),
// 		zap.Any("id", id),
// 	)

// 	return id, nil
// }
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
		constants.INFO_DATABASE_INSERT,
		zap.String("query", query),
		zap.Any("id", id),
	)

	return id, err
}

// QueryMany will query multiple rows
func (dm *DatabaseManager) QueryMany(query string) (*sql.Rows, error) {
	res, err := dm.db.Query(query)
	dm.logger.Info(
		constants.INFO_DATABASE_QUERY,
		zap.String("query", query),
	)
	return res, err
}
