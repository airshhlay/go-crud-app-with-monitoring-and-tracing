package db

import (
	"context"
	"database/sql"
	"fmt"
	"itemService/tracing"

	"itemService/config"
	"itemService/constants"
	metrics "itemService/metrics"

	"github.com/go-sql-driver/mysql"
	ot "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const (
	mysqlInsertRow = "databaseManager.InsertRow"
	mysqlQueryOne  = "databaseManager.QueryOne"
	mysqlQueryRows = "databaseManager.QueryRows"
	mysqlDeleteOne = "databaseManager.DeleteOne"
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
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	// get database handle
	db, err := sql.Open(constants.MySQL, cfg.FormatDSN())
	if err != nil {
		logger.Fatal(
			constants.ErrorDatabaseConnectionMsg,
			zap.Int32(constants.ErrorCode, constants.ErrorDatabaseConnection),
			zap.Error(err))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal(
			constants.ErrorDatabaseConnectionMsg,
			zap.Int32(constants.ErrorCode, constants.ErrorDatabaseConnection),
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
func (dm *DatabaseManager) QueryOne(ctx context.Context, query string, opName string, destination ...any) error {
	// start tracing span from context
	span, ctx := ot.StartSpanFromContext(ctx, mysqlQueryOne)
	dm.addSpanTags(span, query)
	defer span.Finish()
	successStr := constants.True
	// time database query
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.DatabaseOpDuration.WithLabelValues(dm.config.ServiceLabel, constants.Select, opName, successStr).Observe(v)
	}))
	defer func() {
		// observe duration at the end of this function
		timer.ObserveDuration()
	}()

	res := dm.db.QueryRowContext(ctx, query)
	err := res.Scan(destination...)
	if err != nil {
		dm.logger.Error(
			constants.ErrorDatabaseQueryMsg,
			zap.String(constants.Query, query),
			zap.String(constants.OpName, opName),
			zap.Error(err),
		)
		if err != sql.ErrNoRows {
			// avoid false negatives, a select query can return no rows
			successStr = constants.False
		}
		return err
	}
	dm.logger.Info(
		constants.InfoDatabaseQuery,
		zap.String(constants.Query, query),
	)
	return err
}

// QueryRows executes the given query and returns the queried rows.
func (dm *DatabaseManager) QueryRows(ctx context.Context, query string, opName string) (*sql.Rows, error) {
	// start tracing span from context
	span, ctx := ot.StartSpanFromContext(ctx, mysqlQueryRows)
	dm.addSpanTags(span, query)
	defer span.Finish()
	successStr := constants.True
	// time database query
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.DatabaseOpDuration.WithLabelValues(dm.config.ServiceLabel, constants.Insert, opName, successStr).Observe(v)
	}))
	defer func() {
		// observe duration at the end of this function
		timer.ObserveDuration()
	}()
	rows, err := dm.db.QueryContext(ctx, query)
	dm.logger.Info(
		constants.InfoDatabaseQueryRows,
		zap.String(constants.Query, query),
		zap.Error(err),
	)
	return rows, err
}

// InsertRow will insert a single row and return its ID.
func (dm *DatabaseManager) InsertRow(ctx context.Context, query string, opName string) (int64, error) {
	// start tracing span from context
	span, ctx := ot.StartSpanFromContext(ctx, mysqlInsertRow)
	dm.addSpanTags(span, query)
	defer span.Finish()
	successStr := constants.True
	// time database query
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.DatabaseOpDuration.WithLabelValues(dm.config.ServiceLabel, constants.Insert, opName, successStr).Observe(v)
	}))
	defer func() {
		// observe duration at the end of this function
		timer.ObserveDuration()
	}()

	res, err := dm.db.ExecContext(ctx, query)

	if err != nil {
		dm.logger.Error(
			constants.ErrorDatabaseInsertMsg,
			zap.String(constants.Query, query),
			zap.String(constants.OpName, opName),
			zap.Error(err),
		)
		successStr = constants.False
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		successStr = constants.False
		return 0, err
	}

	dm.logger.Info(
		constants.InfoDatabaseInsert,
		zap.String(constants.Query, query),
		zap.Any(constants.ID, id),
	)

	return id, err
}

// DeleteOne deletes a row from the database and returns the number of rows deleted
func (dm *DatabaseManager) DeleteOne(ctx context.Context, query string, opName string) (int64, error) {
	// start tracing span from context
	span, ctx := ot.StartSpanFromContext(ctx, mysqlDeleteOne)
	dm.addSpanTags(span, query)
	defer span.Finish()
	successStr := constants.True
	// time database query
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.DatabaseOpDuration.WithLabelValues(dm.config.ServiceLabel, constants.Delete, opName, successStr).Observe(v)
	}))
	defer func() {
		// observe duration at the end of this function
		timer.ObserveDuration()
	}()

	res, err := dm.db.ExecContext(ctx, query)
	dm.logger.Info(
		constants.InfoDatabaseDelete,
		zap.String(constants.Query, query),
		zap.Any(constants.Res, res),
	)

	if err != nil {
		dm.logger.Error(
			constants.ErrorDatabaseDeleteMsg,
			zap.String(constants.Query, query),
			zap.String(constants.OpName, opName),
			zap.Error(err),
		)
		successStr = constants.False
		return 0, err
	}

	return res.RowsAffected()
}

func (dm *DatabaseManager) addSpanTags(span ot.Span, statement string) {
	span.SetTag(tracing.DatabaseType, tracing.DatabaseTypeSQL)
	span.SetTag(tracing.DatabaseInstance, dm.config.DbName)
	span.SetTag(tracing.DatabaseUser, dm.config.User)
	span.SetTag(tracing.DatabaseStatement, statement)
	span.SetTag(tracing.Component, tracing.ComponentDB)
}
