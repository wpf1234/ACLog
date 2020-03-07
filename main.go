package main

import (
	"GAIOpsAcLog/db"
	"GAIOpsAcLog/handler"
	"GAIOpsAcLog/middleware"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/web"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token，验证token

		c.Next()
	}
}

func main() {
	service := web.NewService(
		web.Name("aiops.micro.api.v2.userBehavior"),
		web.Version("latest"),
		//web.Address(":12345"),
	)

	service.Init()
	db.Init()
	db.InitES()
	db.InitHBase()
	//go db.SelectData(db.NetDb)
	go db.GetData(db.HC.InnerTable, db.EC.InnerTp)
	go db.GetData(db.HC.NetTable, db.EC.NetTp)

	//go db.Consumer(db.KC.InnerTopic, db.EC.InnerTp)
	//go db.Consumer(db.KC.NetTopic, db.EC.NetTp)

	g := new(handler.Gin)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(Cors(), middleware.RequstLog())

	router.POST("/v2/userBehavior", g.GetData)

	//注册 handler
	service.Handle("/", router)

	//运行 api
	err := service.Run()
	if err != nil {
		log.Error("服务启动错误!", err)
		return
	}
}
