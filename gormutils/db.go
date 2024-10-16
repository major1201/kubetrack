package gormutils

import (
	"os"
	"time"

	"github.com/major1201/kubetrack/utils"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	EnvDBInit            = "DB_INIT"
	EnvDBDriver          = "DB_DRIVER"
	EnvDBConnection      = "DB_CONNECTION"
	EnvDBConnMaxLifetime = "DB_CONNMAXLIFETIME"
	EnvDBMaxOpenConns    = "DB_MAXOPENCONNS"
	EnvDBMaxIdleConns    = "DB_MAXIDLECONNS"
	EnvDBDebug           = "DB_DEBUG"
)

var (
	_db *gorm.DB
)

func init() {
	dbInit := utils.BoolDefault(os.Getenv(EnvDBInit), false)
	if !dbInit {
		return
	}

	// try db connection
	sqlDB, err := GetDB().DB()
	if err != nil {
		logger.Error(errors.WithStack(err), "get sql db failed")
		os.Exit(1)
	}
	if err = sqlDB.Ping(); err != nil {
		logger.Error(errors.WithStack(err), "db connect failed")
		os.Exit(1)
	}
}

func gormInit() {
	dbDriver := os.Getenv(EnvDBDriver)
	dbConnection := os.Getenv(EnvDBConnection)
	dbConnMaxLifetime := time.Duration(utils.IntDefault(os.Getenv(EnvDBConnMaxLifetime), 30)) * time.Second // default 30s
	dbMaxOpenConns := utils.IntDefault(os.Getenv(EnvDBMaxOpenConns), 100)                                   // default 100
	dbMaxIdleConns := utils.IntDefault(os.Getenv(EnvDBMaxIdleConns), 10)                                    // default 10
	dbDebug := utils.BoolDefault(os.Getenv(EnvDBDebug), false)                                              // default false

	dbLogLevel := gormLogger.Error
	if dbDebug {
		dbLogLevel = gormLogger.Info
	}

	if dbDriver == "" {
		logger.Error(nil, "environment not set", "env", EnvDBDriver)
		os.Exit(1)
	}
	if dbConnection == "" {
		logger.Error(nil, "environment not set", "env", EnvDBConnection)
		os.Exit(1)
	}

	logger.Info("connecting db",
		"driver", dbDriver,
		"connection", dbConnection,
		"conn_max_lifetime", dbConnMaxLifetime,
		"max_open_conns", dbMaxOpenConns,
		"max_idle_conns", dbMaxIdleConns,
	)

	var err error
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger: gormLogger.New(NewLogrAdapter(logger), gormLogger.Config{
			LogLevel: dbLogLevel,
		}),
	}

	var dialector gorm.Dialector
	switch dbDriver {
	case "postgres":
		dialector = postgres.Open(dbConnection)
	case "mysql":
		dialector = mysql.Open(dbConnection)
	default:
		logger.Error(nil, "unknown db driver", "driver", dbDriver)
		os.Exit(1)
	}

	_db, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		logger.Error(errors.WithStack(err), "db init failed")
		os.Exit(1)
	}

	sqlDB, err := _db.DB()
	if err != nil {
		logger.Error(errors.WithStack(err), "get sql db failed")
		os.Exit(1)
	}
	sqlDB.SetConnMaxLifetime(dbConnMaxLifetime)
	sqlDB.SetMaxOpenConns(dbMaxOpenConns)
	sqlDB.SetMaxIdleConns(dbMaxIdleConns)
}

// GetDB get gorm DB client
func GetDB() *gorm.DB {
	if _db == nil {
		gormInit()
	}
	return _db
}
