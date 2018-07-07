package controller

import (
	"SecKill/sk_proxy/service/srv_sec"
	"fmt"
	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"time"
)

//获取商品列表
func SecInfoList(ctx *gin.Context) {
	list, _, err := srv_sec.SecInfoList()
	if err != nil {
		log.Printf("SecInfoList, err : %v", err)
		ctx.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "failed",
			"data": list,
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": list,
	})
	return
}

//获取商品详情
func SecInfo(ctx *gin.Context) {
	//商品ID
	productId, _ := com.StrTo(ctx.Query("product_id")).Int()
	fmt.Println("productId : ", productId)
	data := srv_sec.SecInfo(productId)

	ctx.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
	return
}

func SecKill(ctx *gin.Context) {
	productId := ctx.PostForm("product_id") //商品Id
	userId := ctx.PostForm("user_id")       //用户id
	source := ctx.PostForm("src")           //来源
	authCode := ctx.PostForm("auth_code")   //授权码
	secTime := ctx.PostForm("time")         //时间戳
	nance := ctx.PostForm("nance")          //随机字符串

	secRequest := srv_sec.NewSecRequest()
	secRequest.AuthCode = authCode
	secRequest.Nance = nance
	secRequest.ProductId, _ = com.StrTo(productId).Int()
	secRequest.SecTime = secTime
	secRequest.Source = source
	secRequest.UserAuthSign = ctx.Request.Header.Get("AuthSign")
	secRequest.UserId, _ = com.StrTo(userId).Int()
	secRequest.AccessTime = time.Now().Unix()
	if len(ctx.Request.RemoteAddr) > 0 {
		secRequest.ClientAddr = strings.Split(ctx.Request.RemoteAddr, ":")[0]
	}
	secRequest.ClientRefence = ctx.Request.Referer()
	secRequest.CloseNotify = ctx.Writer.CloseNotify()

	data, code, err := srv_sec.SecKill(secRequest)
	if err != nil {
		ctx.JSON(400, map[string]interface{}{
			"msg":  err.Error(),
			"code": code,
			"data": "",
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
	return
}
