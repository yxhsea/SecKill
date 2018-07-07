package activity

import (
	"SecKill/sk_admin/model"
	"SecKill/sk_admin/service"
	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
	"log"
)

func CreateActivity(ctx *gin.Context) {
	activity := &model.Activity{}

	//活动名称
	activity.ActivityName = ctx.PostForm("activity_name")
	//商品Id
	activity.ProductId, _ = com.StrTo(ctx.PostForm("product_id")).Int()
	//活动开始时间
	activity.StartTime, _ = com.StrTo(ctx.PostForm("start_time")).Int64()
	//活动结束时间
	activity.EndTime, _ = com.StrTo(ctx.PostForm("end_time")).Int64()
	//商品数量
	activity.Total, _ = com.StrTo(ctx.PostForm("total")).Int()
	//商品速度
	activity.Speed, _ = com.StrTo(ctx.PostForm("speed")).Int()
	//购买限制
	activity.BuyLimit, _ = com.StrTo(ctx.PostForm("buy_limit")).Int()
	activity.BuyRate, _ = com.StrTo(ctx.PostForm("buy_rate")).Float64()

	activityServer := service.NewActivityService()
	if err := activityServer.CreateActivity(activity); err != nil {
		log.Printf("ActivityServer.CreateActivity, Error : %v", err)
		ctx.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "failed",
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "success",
	})
	return
}

func GetActivityList(ctx *gin.Context) {
	ActivityService := service.NewActivityService()
	activityList, err := ActivityService.GetActivityList()
	if err != nil {
		log.Printf("ActivityService.GetActivityList, err : %v", err)
		ctx.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "failed",
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": activityList,
	})
	return
}
