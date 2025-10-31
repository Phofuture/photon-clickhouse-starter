package clickhouseStarter

import (
	CHDatabase "github.com/Phofuture/photon-clickhouse-starter/CHDatabase"
	"github.com/Phofuture/photon-core-starter/core"
)

func init() {
	core.RegisterAddModule(CHDatabase.Start)
}
