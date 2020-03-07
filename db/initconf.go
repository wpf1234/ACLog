package db

import (
	"GAIOpsAcLog/models"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"sync"
	"time"
)

var (
	defaultRoot = "app"
	//key="mysql"
	m  sync.RWMutex
	KC models.KafkaConf
	EC models.ESConf
	LC models.LogConf
	HC models.HBaseConf
	MC models.MysqlConf
	//NetDb *gorm.DB
	//IntranetDb *gorm.DB
)

func LoggerToFile(logPath, logFile string) {
	// 日志文件
	fileName := path.Join(logPath, logFile)
	// 写入文件
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("写入文件失败: ", err)
		return
	}
	// 实例化
	//logger := log.New()
	// 设置输出
	log.SetOutput(src)
	// 设置日志级别
	log.SetLevel(log.DebugLevel)
	// 设置显示打印该条日志的函数名
	log.SetReportCaller(true)
	// 设置 rotate logs
	logWrite, err := rotatelogs.New(
		// 分割后的文件名
		fileName+".%Y%m%d.log",
		// 生成软链，只想最新的日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割时间间隔
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err!=nil{
		fmt.Println("Error: ",err)
		return
	}
	writeMap := lfshook.WriterMap{
		log.InfoLevel:  logWrite,
		log.FatalLevel: logWrite,
		log.DebugLevel: logWrite,
		log.WarnLevel:  logWrite,
		log.ErrorLevel: logWrite,
		log.PanicLevel: logWrite,
	}
	lfHook := lfshook.NewHook(writeMap, &log.JSONFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableTimestamp: false,
		DataKey:          "",
		FieldMap:         nil,
		CallerPrettyfier: nil,
		PrettyPrint:      false,
	})
	log.AddHook(lfHook)
}

func Init() {
	m.Lock()
	defer m.Unlock()

	err := config.Load(file.NewSource(
		file.WithPath("./conf/application.yml"),
	))

	if err != nil {
		log.Error("加载配置文件出错: ", err)
		return
	}

	if err := config.Get(defaultRoot, "log").Scan(&LC); err != nil {
		log.Error("获取Log配置文件失败: ", err)
		return
	}
	log.Info("读取Log配置成功!")

	//logger := log.New()
	LoggerToFile(LC.LogPath, LC.LogFile)
	//logger.AddHook(hook)

	//if err := config.Get(defaultRoot, "kafka").Scan(&KC); err != nil {
	//	logs.Error("获取kafka配置文件失败: ", err)
	//	return
	//}
	//logrus.Info("读取kafka配置成功!")

	if err := config.Get(defaultRoot, "elastic").Scan(&EC); err != nil {
		log.Error("获取ES配置文件失败: ", err)
		return
	}
	log.Info("读取ES配置成功!")

	if err := config.Get(defaultRoot, "hbase").Scan(&HC); err != nil {
		log.Error("获取HBase配置文件失败: ", err)
		return
	}
	log.Info("读取HBase配置成功!")

	if err := config.Get(defaultRoot, "mysql").Scan(&MC); err != nil {
		log.Error("获取配置失败: ", err)
		return
	}
	log.Info("读取数据库配置成功!")

	// 二套网
	//NetDb,err=gorm.Open("mysql",MC.NetDsn)
	//if err!=nil{
	//	logrus.Error("Open net log database failed: ",err)
	//	return
	//}
	//
	//NetDb.DB().SetMaxIdleConns(10)
	//NetDb.DB().SetMaxOpenConns(20)
	//NetDb.DB().SetConnMaxLifetime(20*time.Second)
	//
	//if err:=NetDb.DB().Ping();err!=nil{
	//	logrus.Error("数据库连接失败: ",err)
	//	return
	//}
	//logrus.Info("连接two_net成功!")
	//
	//// 内网
	//IntranetDb,err=gorm.Open("mysql",MC.IntranetDsn)
	//if err!=nil{
	//	logrus.Error("Open intranet log database failed: ",err)
	//	return
	//}
	//
	//IntranetDb.DB().SetMaxIdleConns(10)
	//IntranetDb.DB().SetMaxOpenConns(10)
	//IntranetDb.DB().SetConnMaxLifetime(20*time.Second)
	//
	//if err:=IntranetDb.DB().Ping();err!=nil{
	//	logrus.Error("数据库连接失败: ",err)
	//	return
	//}
	//logrus.Info("连接intranet成功!")
}
