package setup

import (
	"github.com/gin-gonic/gin"
	"log"
)

//初始化Http服务
func InitServer(host string) {
	router := gin.Default()
	setupRouter(router)
	err := router.Run(host)
	if err != nil {
		log.Printf("Init http server. Error : %v", err)
	}
}
