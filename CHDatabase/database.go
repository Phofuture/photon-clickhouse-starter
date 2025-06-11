package chdatabase

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
	"gorm.io/gorm"
)

var masterDb, slaveDb *gorm.DB
var clickConn clickhouse.Conn

// 取得 master database
func Master(ctx context.Context) *gorm.DB {
	return masterDb.WithContext(ctx)
}

// 取得 slave database
func Slave(ctx context.Context) *gorm.DB {
	return slaveDb.WithContext(ctx)
}

func Conn() clickhouse.Conn {
	return clickConn
}
