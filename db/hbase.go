package db

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"io"
	"time"
)

var HClient gohbase.Client

func InitHBase() {
	option := gohbase.EffectiveUser(HC.User)
	HClient = gohbase.NewClient(HC.Host, option)
	log.Info("HBase: " + HC.Host + " connected!")
}

func PutsData(table, rowKey string, values map[string]map[string][]byte) error {
	putRequest, err := hrpc.NewPutStr(context.Background(), table, rowKey, values)
	if err != nil {
		log.Error("Hrpc new put error: ", err)
		return err
	}
	_, err = HClient.Put(putRequest)
	if err != nil {
		log.Error("HBase put error: ", err)
		return err
	}
	return nil
}

func GetData(table, index string) {
	for {
		tm := time.Now()
		endRow := tm.Format("20060102150405") + "_"
		m, _ := time.ParseDuration("-21m")
		startRow := tm.Add(m).Format("20060102150405") + "_"
		//log.Info("Start Row: ", startRow, "\nEnd Row: ", endRow)

		getRequest, err := hrpc.NewScanRangeStr(context.Background(), table, startRow, endRow)
		if err != nil {
			log.Println("创建数据请求体失败: ", err)
			return
		}
		scan := HClient.Scan(getRequest)
		var res []*hrpc.Result
		for {
			rsp, err := scan.Next()
			if err == io.EOF || rsp == nil {
				log.Println(err, rsp)
				break
			}
			if err != nil {
				log.Println("Scan next error: ", err)
				return
			}
			res = append(res, rsp)
		}
		myMp := make(map[string]map[string]string, 0)
		for _, cells := range res {

			mp := make(map[string]string)
			for _, cell := range cells.Cells {

				rk := string(cell.Row)
				//cf:=string(cell.Family)
				//tag:=string(cell.Tags)
				qu := string(cell.Qualifier)
				val := string(cell.Value)
				//	fmt.Println("value: ",val)
				mp[qu] = val
				myMp[rk] = mp
			}
		}
		fmt.Println("Data length: ", len(myMp))
		CopyDataToEs(index, myMp)

		ticker := time.NewTicker(20 * time.Minute)
		<-ticker.C
	}
}
