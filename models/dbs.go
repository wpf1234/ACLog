package models

// HBase
type HBaseConf struct {
	Host string `json:"host"`
	User string `json:"user"`
	Cf string `json:"cf"`
	InnerTable string `json:"inner_table"`
	NetTable string `json:"net_table"`
}

// ES
type ESConf struct {
	Host string `json:"host"`
	InnerTp string `json:"inner_tp"`
	NetTp string `json:"net_tp"`
	User string `json:"user"`
	Password string `json:"password"`
}

// kafka
type KafkaConf struct {
	Broker string `json:"broker"`
	InnerTopic string `json:"inner_topic"`
	NetTopic string `json:"net_topic"`
}

// mysql
type MysqlConf struct {
	IntranetDsn string`json:"intranet_dsn"`
	NetDsn string `json:"net_dsn"`
	Table string `json:"table"`
}