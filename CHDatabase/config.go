package chdatabase

type ConnectData struct {
	Hosts []string `yaml:"host"`
	Auth  struct {
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"auth"`
	ClientInfo struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"clientInfo"`
}

type Config struct {
	ClickHouse struct {
		Master ConnectData `yaml:"master"`
		Slave  ConnectData `yaml:"slave"`
	} `yaml:"clickhouse"`
}

var config Config
