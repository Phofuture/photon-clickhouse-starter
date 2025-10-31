package chdatabase

import (
	"context"

	"github.com/Phofuture/photon-core-starter/configuration"
	"github.com/Phofuture/photon-core-starter/log"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	gormClickHouse "gorm.io/driver/clickhouse"

	"gorm.io/gorm"
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

// 連線資料庫
func connect(ctx context.Context, connectData ConnectData) (clickConn clickhouse.Conn, err error) {
	clickConn, err = clickhouse.Open(&clickhouse.Options{
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
	})

	if err != nil {
		log.Logger().Error(ctx, "open click house conn error", "error", err)
		return nil, err
	}

	return clickConn, nil
}

// 連線資料庫 DB
func connectDB(ctx context.Context, connectData ConnectData) (db *gorm.DB, err error) {
	sqlDB := clickhouse.OpenDB(&clickhouse.Options{
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
	})
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
