package utils

import (
	"GAIOpsAcLog/models"
	"encoding/xml"
	"reflect"
	"regexp"
	"strings"
)

func Handle(acLog models.AcLog, cf string) map[string]map[string][]byte{
	var result = make(map[string]map[string][]byte,0)
	var ac=make(map[string][]byte)
	var t xml.Token
	var err error
	var netLog models.NetLog
	var mac, resource string
	var urls []string
	input := strings.NewReader(acLog.Result)
	decoder := xml.NewDecoder(input)
	for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.CharData:
			content := string([]byte(token))
			// 匹配 mac 地址
			reg1 := regexp.MustCompile("[0-9a-fA-F]{2}(-[0-9a-fA-F]{2}){5}")
			// 匹配 网页地址
			reg2 := regexp.MustCompile(`[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?$`)
			if reg1.MatchString(content) {
				mac = content
			} else if reg2.MatchString(content) {
				urls = append(urls, content)
			}
		default:

		}
	}

	urls = Distinct(urls)
	if len(urls) > 1 {
		for j := 0; j < len(urls)-1; j++ {
			if len(urls[j]) > len(urls[j+1]) {

				resource = urls[j]
			} else {

				resource = urls[j+1]
			}
		}
	} else if len(urls) == 0 {
		resource = ""
	} else {
		resource = urls[0]
	}

	netLog.RecordId = acLog.RecordId
	netLog.User = acLog.User
	netLog.RecordTime = acLog.RecordTime
	netLog.Mac = mac
	netLog.Url = resource
	netLog.Serv = acLog.Server
	netLog.App = acLog.App
	netLog.DestIp = acLog.DstIp
	netLog.ServPort = acLog.ServPort
	//result=acLog.User+","+acLog.HostIp+","+acLog.DstIp+","+acLog.Server+","+acLog.App+","+mac+
	//	","+strconv.Itoa(acLog.ServPort)+","+acLog.RecordTime+","+resource

	//ac=append(ac,result)
	// 结构体转 map
	key := reflect.TypeOf(netLog)
	value := reflect.ValueOf(netLog)
	for i := 0; i < key.NumField(); i++ {
		ac[key.Field(i).Name] = []byte(value.Field(i).String())
		result[cf]=ac
	}
	//fmt.Println("数据大小: ", len(result))
	return result
}
