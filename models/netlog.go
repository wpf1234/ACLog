package models

//type NetLog struct {
//	RecordId int `json:"record_id"`
//	User string `json:"user"`
//	HostIp string `json:"host_ip"`
//	DstIp string `json:"dst_ip"`
//	Serv string `json:"serv"`
//	App string `json:"app"`
//	SrcPort int `json:"src_port"`
//	ServPort int `json:"serv_port"`
//	RecordTime string `json:"record_time"`
//	Mac string `json:"mac"`
//	Url string `json:"url"`
//}

type NetLog struct {
	RecordId   string `json:"record_id"`
	App        string `json:"app"`
	DestIp     string `json:"dest_ip"`
	Mac        string `json:"mac"`
	RecordTime string `json:"record_time"`
	Serv       string `json:"serv"`
	Url        string `json:"url"`
	User       string `json:"user"`
	ServPort   string `json:"serv_port"`
}

type Action struct {
	Id           string `json:"id"`
	Date         string `json:"date"`
	User         string `json:"user"`
	TargetIp     string `json:"targetIp"`
	Behavior     string `json:"behavior"`
	BehaviorType string `json:"behaviorType"`
	Url          string `json:"url"`
	Mac          string `json:"mac"`
	Port         string `json:"port"`
}

type AcLog struct {
	RecordId   string `json:"record_id"`
	App        string `json:"app"`
	User       string `json:"user"`
	HostIp     string `json:"host_ip"`
	DstIp      string `json:"dst_ip"`
	Server     string `json:"server"`
	ServPort   string `json:"serv_port"`
	RecordTime string `json:"record_time"`
	Result     string `json:"result"`
}

type LogConf struct {
	LogPath string `json:"log_path"`
	LogFile string `json:"log_file"`
}
