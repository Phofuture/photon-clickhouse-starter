package clickhouseStarter

import (
	CHDatabase "github.com/Phofuture/photon-clickhouse-starter/CHDatabase"
	"github.com/dennesshen/photon-core-starter/core"
)

func init() {
	core.RegisterAddModule(CHDatabase.Start)
}
