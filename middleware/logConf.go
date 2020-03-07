package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

//func loggerToFile(logPath,logFile string) *lfshook.LfsHook {
//	// 日志文件
//	fileName := path.Join(logPath, logFile)
//	// 写入文件
//	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
//	if err != nil {
//		fmt.Print("写入文件失败: ", err)
//	}
//	// 实例化
//	logger := logrus.New()
//	// 设置输出
//	logger.Out = src
//	// 设置日志级别
//	logger.SetLevel(logrus.DebugLevel)
//	// 设置 rotate logs
//	logWrite, err := rotatelogs.New(
//		// 分割后的文件名
//		fileName+".%Y%m%d.log",
//		// 生成软链，只想最新的日志文件
//		rotatelogs.WithLinkName(fileName),
//		// 设置最大保存时间
//		rotatelogs.WithMaxAge(7*24*time.Hour),
//		// 设置日志切割时间间隔
//		rotatelogs.WithRotationTime(24*time.Hour),
//	)
//	writeMap := lfshook.WriterMap{
//		logrus.InfoLevel:  logWrite,
//		logrus.FatalLevel: logWrite,
//		logrus.DebugLevel: logWrite,
//		logrus.WarnLevel:  logWrite,
//		logrus.ErrorLevel: logWrite,
//		logrus.PanicLevel: logWrite,
//	}
//	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
//		TimestampFormat:  "2006-01-02 15:04:05",
//		DisableTimestamp: false,
//		DataKey:          "",
//		FieldMap:         nil,
//		CallerPrettyfier: nil,
//		PrettyPrint:      false,
//	})
//	return lfHook
//}

func RequstLog() gin.HandlerFunc {
	//// 日志文件
	//fileName := path.Join(db.LC.LogPath, db.LC.LogFile+"_req")
	//// 写入文件
	//src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	//if err != nil {
	//	fmt.Println("写入文件失败: ", err)
	//
	//}
	//// 实例化
	//logger := logrus.New()
	//// 设置输出
	//logrus.SetOutput(src)
	//// 设置日志级别
	//logrus.SetLevel(logrus.DebugLevel)
	//// 打印该条日志的函数名
	//logrus.SetReportCaller(true)
	//// 设置 rotate logs
	//logWrite, err := rotatelogs.New(
	//	// 分割后的文件名
	//	fileName+".%Y%m%d.log",
	//	// 生成软链，只想最新的日志文件
	//	rotatelogs.WithLinkName(fileName),
	//	// 设置最大保存时间
	//	rotatelogs.WithMaxAge(7*24*time.Hour),
	//	// 设置日志切割时间间隔
	//	rotatelogs.WithRotationTime(24*time.Hour),
	//)
	//writeMap := lfshook.WriterMap{
	//	logrus.InfoLevel:  logWrite,
	//	logrus.FatalLevel: logWrite,
	//	logrus.DebugLevel: logWrite,
	//	logrus.WarnLevel:  logWrite,
	//	logrus.ErrorLevel: logWrite,
	//	logrus.PanicLevel: logWrite,
	//}
	//lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
	//	TimestampFormat:  "2006-01-02 15:04:05",
	//	DisableTimestamp: false,
	//	DataKey:          "",
	//	FieldMap:         nil,
	//	CallerPrettyfier: nil,
	//	PrettyPrint:      false,
	//})
	//logger.AddHook(lfHook)

	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		// 执行时间
		latency := end.Sub(start)
		// 请求方式
		method := c.Request.Method
		// 请求路由
		uri := c.Request.RequestURI
		// 状态码
		code := c.Writer.Status()
		// 请求IP
		clientIp := c.ClientIP()
		// 日志格式
		log.WithFields(log.Fields{
			"code":    code,
			"latency": latency,
			"IP":      clientIp,
			"method":  method,
			"uri":     uri,
		}).Info()
	}
}
