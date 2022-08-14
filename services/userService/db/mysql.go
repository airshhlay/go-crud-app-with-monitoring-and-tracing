package db

import (
	"context"
	"database/sql"
	"fmt"
	config "userService/config"
	constants "userService/constants"
	metrics "userService/metrics"
	"userService/tracing"

	"github.com/go-sql-driver/mysql"
	ot "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const (
	mysqlInsertRow = "mysql.InsertRow"
	mysqlQueryOne  = "mysql.QueryOne"
)

// DatabaseManager is a struct containing a reference to the database connection, logger, and the database config
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

	return &dbManager, nil
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

func (dm *DatabaseManager) addSpanTags(span ot.Span, statement string) {
	span.SetTag(tracing.DatabaseType, tracing.DatabaseTypeSQL)
	span.SetTag(tracing.DatabaseInstance, dm.config.DbName)
	span.SetTag(tracing.DatabaseUser, dm.config.User)
	span.SetTag(tracing.DatabaseStatement, statement)
}
