package handler

import (
	"GAIOpsAcLog/db"
	"GAIOpsAcLog/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Gin struct {
}

//如果数据量小（10000条内），或者只关注结果集的TopN数据，可以使用from / size 分页，简单粗暴
//数据量大，深度翻页，后台批处理任务（数据迁移）之类的任务，使用 scroll 方式
//数据量大，深度翻页，用户实时、高并发查询需求，使用 search after 方式

func (g *Gin) GetData(c *gin.Context) {
	var total int64
	var index string
	var list []models.Action
	var response models.Response
	var res *elastic.SearchResult
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error("获取数据失败: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err,
			"data":    "",
		})
		return
	}

	req := models.Request{}
	json.Unmarshal(body, &req)

	//fmt.Println(req.Filter.Date == "")
	multiple := elastic.NewBoolQuery()

	multiple.Must(elastic.NewWildcardQuery("user.keyword", "*"+req.Filter.User+"*"))
	multiple.Must(elastic.NewWildcardQuery("dest_ip.keyword", "*"+req.Filter.TargetIp+"*"))
	multiple.Must(elastic.NewWildcardQuery("app.keyword", "*"+req.Filter.Behavior+"*"))
	multiple.Must(elastic.NewWildcardQuery("serv.keyword", "*"+req.Filter.BehaviorType+"*"))
	multiple.Must(elastic.NewWildcardQuery("url.keyword", "*"+req.Filter.Url+"*"))
	multiple.Must(elastic.NewWildcardQuery("mac.keyword", "*"+req.Filter.Mac+"*"))
	if req.Filter.Port!=""{
		multiple.Must(elastic.NewMatchQuery("serv_port", req.Filter.Port))
	}

	if strings.Contains(req.Filter.Date,",") {
		tm := strings.Split(req.Filter.Date, ",")
		start := tm[0]
		st, _ := time.Parse("2006-01-02 15:04:05", start)
		start = st.Format("20060102150405")
		fmt.Println("开始时间: ", start)
		end := tm[1]
		if end == "" {
			end = time.Now().Format("20060102150405")
		} else {
			e, _ := time.Parse("2006-01-02 15:04:05", end)
			end = e.Format("20060102150405")
		}
		fmt.Println("结束时间: ", end)
		// 按照时间进行范围查询 range query
		multiple.Filter(elastic.NewRangeQuery("record_time").
			Gt(start).Lt(end))

	}
	//if !strings.Contains(req.Filter.Date,",") && req.Filter.Date!=""{
	//	st, _ := time.Parse("2006-01-02 15:04:05", req.Filter.Date)
	//	start:=st.Format("20060102150405")
	//	fmt.Println("开始时间: ", start)
	//	end := time.Now().Format("20060102150405")
	//	fmt.Println("结束时间: ", end)
	//	// 按照时间进行范围查询 range query
	//	multiple.Filter(elastic.NewRangeQuery("record_time").
	//		Gt(start).Lt(end))
	//}

	// 方法一：浅分页
	// 存在的问题，在查询的数量大于from+size默认配置的 1W 时，会报错；或者会在后面产生超时问题
	var data models.NetLog
	if req.Filter.Type == 1 {
		// 内网
		index = db.EC.InnerTp

	} else {
		// 其他情况全部显示二套网
		index = db.EC.NetTp

	}


	total, err = db.ESClient.Count(index).Type(index).
		Query(multiple).Do(context.Background())
	if err != nil {
		log.Error("查询总数失败: ", err)
		c.JSON(200, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    response,
			"message": err,
		})
		return
	}
	/*
	** sort 报错
	** elastic: Error 400 (Bad Request): all shards failed [type=search_phase_execution_exception]
	** 解决方法
	** 用region.keyword进行聚合，排序
	 */
	res, err = db.ESClient.Search(index).Type(index).
		Query(multiple).Sort("record_time.keyword", false).
		From(req.Skip * req.Limit).Size(req.Limit).
		Do(context.Background())
	if err != nil {
		log.Error("查询数据失败: ", err)
		c.JSON(200, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    response,
			"message": err,
		})
		return
	}

	for _, item := range res.Each(reflect.TypeOf(data)) {
		ac := models.Action{}
		t := item.(models.NetLog)
		ac.Id = t.RecordId
		ac.User = t.User
		tm, _ := time.Parse("20060102150405", t.RecordTime)
		ac.Date = tm.Format("2006-01-02 15:04:05")
		ac.TargetIp = t.DestIp
		ac.Behavior = t.App
		ac.BehaviorType = t.Serv
		ac.Url = t.Url
		ac.Mac = t.Mac
		ac.Port = t.ServPort

		list = append(list, ac)
	}

	// 方法二：使用 scroll 滚动查询
	//var index string
	//docs := 0
	//if req.Filter.Type == 1 {
	//	// 内网
	//	index = db.EC.InnerTp
	//} else {
	//	// 其他情况全部显示二套网
	//	index = db.EC.NetTp
	//}
	//
	//total, err = db.ESClient.Count(index).Type(index).
	//	Query(multiple).Do(context.Background())
	//if err != nil {
	//	log.Error("查询总数失败: ", err)
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"code":    http.StatusInternalServerError,
	//		"data":    response,
	//		"message": err,
	//	})
	//	return
	//}
	//
	//svc := db.ESClient.Scroll(index).Query(multiple).
	//	Sort("record_time.keyword", false).Size(req.Limit)
	//for i := 0; i <= req.Skip; i++ {
	//	res, err = svc.Do(context.Background())
	//	if err != nil {
	//		log.Error("查询数据失败: ", err)
	//		return
	//	}
	//	if err == io.EOF {
	//		return
	//	}
	//
	//	for _, hit := range res.Hits.Hits {
	//		if hit.Index != index {
	//			log.WithFields(log.Fields{
	//				"Hit.Index": index,
	//				"Get index": hit.Index,
	//			}).Error("索引值不一致:")
	//			return
	//		}
	//		data, err := hit.Source.MarshalJSON()
	//		if err != nil {
	//			log.Error("Error: ", err)
	//			return
	//		}
	//		item := models.NetLog{}
	//		err = json.Unmarshal(data, &item)
	//		if err != nil {
	//			log.Error("Json Unmarshal error: ", err)
	//			return
	//		}
	//		ac := models.Action{}
	//
	//		ac.Id = item.RecordId
	//		ac.User = item.User
	//		tm, _ := time.Parse("20060102150405", item.RecordTime)
	//		ac.Date = tm.Format("2006-01-02 15:04:05")
	//		ac.TargetIp = item.DstIp
	//		ac.Behavior = item.App
	//		ac.BehaviorType = item.Serv
	//		ac.Url = item.Url
	//		ac.Mac = item.Mac
	//		ac.Port = item.ServPort
	//
	//		list = append(list, ac)
	//		docs++
	//	}
	//}
	//list = list[req.Skip*req.Limit:]

	response.List = list
	response.Total = total

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    response,
		"message": "获取数据成功!",
	})
}
