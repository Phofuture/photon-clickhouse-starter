package clickhouseStarter

import (
	CHDatabase "git.ggyyonline.com/aggregator/photon-core/photon-clickhouse-starter/CHDatabase"
	"github.com/dennesshen/photon-core-starter/core"
)

func init() {
	core.RegisterAddModule(CHDatabase.Start)
}
