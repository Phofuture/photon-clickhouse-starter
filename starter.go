package clickhouseStarter

import (
	"photon-core-starter/core"
	"photon-clickhouse-starter/CHDatabase"
)

func init() {
	core.RegisterAddModule(CHDatabase.Start)
}
