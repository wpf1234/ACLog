package db

import (
	"GAIOpsAcLog/models"
	"GAIOpsAcLog/utils"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"strings"
)

var ESClient *elastic.Client
var key = []byte("aiops*2019*1801.")

func InitES() {
	var err error
	host := strings.Split(EC.Host, ",")
	str, err := base64.StdEncoding.DecodeString(EC.Password)
	if err != nil {
		log.Error("Base64 decode failed: ", err)
		return
	}
	pwd, err := utils.AesDecrypt(str, key)
	if err != nil {
		log.Error("Decrypt failed: ", err)
		return
	}
	password := string(pwd)
	ESClient, err = elastic.NewClient(
		elastic.SetBasicAuth(EC.User, password),
		elastic.SetURL(host...),
		elastic.SetSniff(false),

	)
	if err != nil {
		log.Error("创建 ES Client 失败: ", err)
		return
	}

	info, code, err := ESClient.Ping(host[0]).Do(context.Background())
	if err != nil {
		log.Error("连接 ES 失败: ", err)
		return
	}
	logs.Info("Elasticsearch returned with code %d and version %s\n",
		code, info.Version.Number)
}

func CopyDataToEs(index string, mp map[string]map[string]string) {
	var netLog models.NetLog
	var id string
	ctx := context.TODO()
	bulkRequest := ESClient.Bulk()
	for k, v := range mp {
		id = k
		js, err := json.Marshal(v)
		if err != nil {
			log.Error("数据转换失败: ", err)
			return
		}
		json.Unmarshal(js, &netLog)
		str := strings.Split(k, "_")
		netLog.RecordTime = str[0]
		netLog.RecordId = str[1]
		doc := elastic.NewBulkIndexRequest().
			Index(index).
			Type(index).
			Id(id).
			Doc(netLog)
		bulkRequest = bulkRequest.Add(doc)
	}
	res, err := bulkRequest.Do(ctx)
	if err != nil {
		log.Error("批量插入失败: ", err)
		return
	}
	failed := res.Failed()
	it := len(failed)
	log.WithFields(log.Fields{
		"Index":     index,
		"Is error":  res.Errors,
		"Error num": it,
	}).Info("Result")
}
