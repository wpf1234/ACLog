package db

import (
	"GAIOpsAcLog/models"
	"GAIOpsAcLog/utils"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

func SelectData(db *gorm.DB){
	for{
		//var acLog []models.AcLog
		var maxId int
		var recordDate string
		date:=time.Now().Format("20060102")
		d:=db.Raw("select date_log,record_id from log_id")
		d.Row().Scan(&recordDate,&maxId)
		if recordDate != date{
			maxId=0
			goto LOOP
		}else{
			goto LOOP
		}
	LOOP:
		table:=date+MC.Table
		t:=time.Now()
		end:=t.Format("15:04:05")
		m,_:=time.ParseDuration("-10m")
		// 十分钟前
		start:=t.Add(m).Format("15:04:05")
		fmt.Println("Start: ",start,"\nEnd: ",end)
		sql:=fmt.Sprintf("select record_id,user,host_ip,dst_ip,serv,app," +
			"serv_port,record_time,result from %s where record_id > %d " +
			"and record_time between '%s' and '%s'",
			table,maxId,start,end)

		raw:=db.Raw(sql)
		rows,err:=raw.Rows()
		if err!=nil{
			logs.Error("获取AC数据失败: ",err)
			return
		}

		for rows.Next(){
			var acLog models.AcLog
			var tm string
			rows.Scan(&maxId,&acLog.User,&acLog.HostIp,&acLog.DstIp,
				&acLog.Server,&acLog.App,&acLog.ServPort,
				&tm,&acLog.Result)
			tm=time.Now().Format("2006-01-02")+" "+tm
			rt,_:=time.Parse("2006-01-02 15:04:05",tm)
			acLog.RecordTime=rt.Format("20060102150405")
			//inNets=append(inNets,inner)
			intIP:=strconv.FormatInt(utils.InetAtoN(acLog.HostIp),10)
			intIP=utils.ReverseString(intIP)
			rowKey:=acLog.RecordTime+"_"+intIP
			value:=utils.Handle(acLog,HC.Cf)
			fmt.Println("RowKey: ",rowKey,"\nValue: ",value)
			//maxId=acLog.RecordId
		}
		fmt.Println("最后一条的记录ID: ",maxId)
		d=db.Exec("update log_id set record_id=?,date_log=?",maxId,date)
		fmt.Println("Rows affected: ",d.RowsAffected)
		//db2:=NetDb.Raw(sql)
		//rows2,err:=db2.Rows()
		//if err!=nil{
		//	logs.Error("获取二套网数据失败: ",err)
		//	return
		//}
		//for rows2.Next(){
		//	var net models.AcLog
		//	var tm string
		//	rows1.Scan(&net.RecordId,&net.User,&net.HostIp,&net.DstIp,
		//		&net.Server,&net.App,&net.ServPort,
		//		&tm,&net.Result)
		//	net.RecordTime=time.Now().Format("2006-01-02")+" "+tm
		//	nets=append(nets,net)
		//}

		//inner:=utils.Handle(inNets)
		//two:=utils.Handle(nets)

		// 写入 HBase

		//SendToKafka(KC.InnerTopic,KC.InnerKey,inner)
		//SendToKafka(KC.NetTopic,KC.NetKey,two)

		ticker:=time.NewTicker(10*time.Minute)
		<-ticker.C
	}
}
