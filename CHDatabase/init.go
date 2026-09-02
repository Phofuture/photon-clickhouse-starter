package chdatabase

import (
	"context"
	"time"

	"github.com/Phofuture/photon-core-starter/configuration"
	"github.com/Phofuture/photon-core-starter/log"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	gormClickHouse "gorm.io/driver/clickhouse"

	"gorm.io/gorm"
)

// 連線池預設值（driver 預設僅 MaxOpenConns=10/MaxIdleConns=5，重查詢併發下會 acquire conn timeout）
const (
	defaultMaxOpenConns = 50
	defaultMaxIdleConns = 10
)

type DbAction func(ctx context.Context, db *gorm.DB) (err error)

var customAction = []DbAction{}

func RegisterDbCustomize(action DbAction) {
	customAction = append(customAction, action)
}

func Start(ctx context.Context) (err error) {
	log.Logger().Info(ctx, "init clickhouse database")
	config, err = configuration.Get[Config](ctx)
	if err != nil {
		log.Logger().Error(ctx, "failed to get clickhouse database config", "error", err)
		return
	}

	if masterDb, err = connectDB(ctx, config.ClickHouse.Master); err != nil {
		log.Logger().Error(ctx, "fail to connect master clickhouse database", "error", err, "config", config)
		return
	}

	if slaveDb, err = connectDB(ctx, config.ClickHouse.Slave); err != nil {
		log.Logger().Error(ctx, "fail to connect slave clickhouse database", "error", err, "config", config)
		return
	}

	for _, action := range customAction {
		if err = action(ctx, masterDb); err != nil {
			log.Logger().Error(ctx, "failed to customize master clickhouse database", "error", err)
			return
		}
		if err = action(ctx, slaveDb); err != nil {
			log.Logger().Error(ctx, "failed to customize slave clickhouse database", "error", err)
			return
		}
	}

	if clickConn, err = connect(ctx, config.ClickHouse.Master); err != nil {
		log.Logger().Error(ctx, "fail to connect slave clickhouse database", "error", err, "config", config)
		return
	}

	return
}

// 組裝連線選項（不含連線池；OpenDB 不允許在 Options 設 pool，須由 sql.DB 設定）
func buildOptions(connectData ConnectData) *clickhouse.Options {
	opt := &clickhouse.Options{
		Addr: connectData.Hosts,
		Auth: clickhouse.Auth{
			Database: connectData.Auth.Database,
			Username: connectData.Auth.Username,
			Password: connectData.Auth.Password,
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{
					Name:    connectData.ClientInfo.Name,
					Version: connectData.ClientInfo.Version,
				},
			},
		},
	}

	if connectData.Pool.DialTimeoutSeconds > 0 {
		opt.DialTimeout = time.Duration(connectData.Pool.DialTimeoutSeconds) * time.Second
	}

	return opt
}

// 解析連線池大小，未設定則套用預設值
func poolSize(connectData ConnectData) (maxOpen, maxIdle int) {
	maxOpen = defaultMaxOpenConns
	if connectData.Pool.MaxOpenConns > 0 {
		maxOpen = connectData.Pool.MaxOpenConns
	}
	maxIdle = defaultMaxIdleConns
	if connectData.Pool.MaxIdleConns > 0 {
		maxIdle = connectData.Pool.MaxIdleConns
	}
	return
}

// 連線資料庫（native，pool 由 Options 設定）
func connect(ctx context.Context, connectData ConnectData) (clickConn clickhouse.Conn, err error) {
	opt := buildOptions(connectData)
	opt.MaxOpenConns, opt.MaxIdleConns = poolSize(connectData)

	clickConn, err = clickhouse.Open(opt)

	if err != nil {
		log.Logger().Error(ctx, "open click house conn error", "error", err)
		return nil, err
	}

	return clickConn, nil
}

// 連線資料庫 DB（OpenDB，pool 須由 sql.DB 設定，不可放進 Options）
func connectDB(ctx context.Context, connectData ConnectData) (db *gorm.DB, err error) {
	sqlDB := clickhouse.OpenDB(buildOptions(connectData))
	maxOpen, maxIdle := poolSize(connectData)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)

	err = sqlDB.Ping()
	if err != nil {
		log.Logger().Error(ctx, "ping click house error", "error", err)
		return nil, err
	}

	db, err = gorm.Open(gormClickHouse.New(gormClickHouse.Config{
		Conn: sqlDB,
	}))
	if err != nil {
		log.Logger().Error(ctx, "gorm open click house error", "error", err)
		return nil, err
	}
	return db, nil
}
