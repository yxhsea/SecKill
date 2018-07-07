package setup

import (
	"SecKill/sk_proxy/controller"
	"github.com/gin-gonic/gin"
)

//设置路由
func setupRouter(router *gin.Engine) {
	//秒杀管理
	router.GET("/sec/info", controller.SecInfo)
	router.GET("/sec/list", controller.SecInfoList)
	router.POST("/sec/kill", controller.SecKill)
}
