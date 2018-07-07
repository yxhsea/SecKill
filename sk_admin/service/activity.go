package service

import (
	"SecKill/sk_admin/config"
	"SecKill/sk_admin/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Unknwon/com"
	"log"
	"time"
)

type ActivityService struct {
}

func NewActivityService() *ActivityService {
	return &ActivityService{}
}

func (p *ActivityService) GetActivityList() ([]map[string]interface{}, error) {
	activityEntity := model.NewActivityModel()
	activityList, err := activityEntity.GetActivityList()
	if err != nil {
		log.Printf("ActivityEntity.GetActivityList, err : %v", err)
		return nil, err
	}

	for _, v := range activityList {
		startTime, _ := com.StrTo(fmt.Sprint(v["start_time"])).Int64()
		v["start_time_str"] = time.Unix(startTime, 0).Format("2006-01-02 15:04:05")

		endTime, _ := com.StrTo(fmt.Sprint(v["end_time"])).Int64()
		v["end_time_str"] = time.Unix(endTime, 0).Format("2006-01-02 15:04:05")

		nowTime := time.Now().Unix()
		if nowTime > endTime {
			v["status_str"] = "已结束"
			continue
		}

		status, _ := com.StrTo(fmt.Sprint(v["status"])).Int()
		if status == model.ActivityStatusNormal {
			v["status_str"] = "正常"
		} else if status == model.ActivityStatusDisable {
			v["status_str"] = "已禁用"
		}
	}

	log.Printf("get activity success, activity list is [%v]", activityList)
	return activityList, nil
}

func (p *ActivityService) CreateActivity(activity *model.Activity) error {
	//写入到数据库
	activityEntity := model.NewActivityModel()
	err := activityEntity.CreateActivity(activity)
	if err != nil {
		log.Printf("ActivityModel.CreateActivity, err : %v", err)
		return err
	}

	//写入到Etcd
	err = p.syncToEtcd(activity)
	if err != nil {
		log.Printf("activity product info sync to etcd failed, err : %v", err)
		return err
	}
	return nil
}

//将商品活动数据同步到Etcd
func (p *ActivityService) syncToEtcd(activity *model.Activity) error {
	etcdKey := config.SecAdminConfCtx.EtcdConf.EtcdSecProductKey
	secProductInfoList, err := p.loadProductFromEtcd(etcdKey)
	if err != nil {
		return err
	}

	var secProductInfo = &model.SecProductInfoConf{}
	secProductInfo.EndTime = activity.EndTime
	secProductInfo.OnePersonBuyLimit = activity.BuyLimit
	secProductInfo.ProductId = activity.ProductId
	secProductInfo.SoldMaxLimit = activity.Speed
	secProductInfo.StartTime = activity.StartTime
	secProductInfo.Status = activity.Status
	secProductInfo.Total = activity.Total
	secProductInfo.BuyRate = activity.BuyRate
	secProductInfoList = append(secProductInfoList, secProductInfo)

	data, err := json.Marshal(secProductInfoList)
	if err != nil {
		log.Printf("json marshal failed, err : %v", err)
		return err
	}

	conn := config.SecAdminConfCtx.EtcdConf.EtcdConn
	_, err = conn.Put(context.Background(), etcdKey, string(data))
	if err != nil {
		log.Printf("put to etcd failed, err : %v, data = [%v]", err, string(data))
		return err
	}

	log.Printf("put to etcd success, data = [%v]", string(data))
	return nil
}

//从Ectd中取出原来的商品数据
func (p *ActivityService) loadProductFromEtcd(key string) ([]*model.SecProductInfoConf, error) {
	log.Println("start get from etcd success")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	rsp, err := config.SecAdminConfCtx.EtcdConf.EtcdConn.Get(ctx, key)
	if err != nil {
		log.Printf("get [%s] from etcd failed, err : %v", key, err)
		return nil, err
	}
	log.Printf("get from etcd success, rsp : %v", rsp)

	var secProductInfo []*model.SecProductInfoConf
	for k, v := range rsp.Kvs {
		log.Printf("key = [%v], value = [%v]", k, v)
		err := json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			log.Printf("Unmsharl second product info failed, err : %v", err)
			return nil, err
		}
		log.Printf("second info conf is [%v]", secProductInfo)
	}

	return secProductInfo, nil
}
