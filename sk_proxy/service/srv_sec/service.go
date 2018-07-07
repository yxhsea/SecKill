package srv_sec

import (
	"SecKill/sk_proxy/config"
	"SecKill/sk_proxy/service/srv_err"
	"SecKill/sk_proxy/service/srv_limit"
	"fmt"
	"log"
	"time"
)

func SecInfo(productId int) (date map[string]interface{}) {
	config.SecKillConfCtx.RWSecProductLock.RLock()
	defer config.SecKillConfCtx.RWSecProductLock.RUnlock()

	v, ok := config.SecKillConfCtx.SecProductInfoMap[productId]
	if !ok {
		return nil
	}

	data := make(map[string]interface{})
	data["product_id"] = productId
	data["start_time"] = v.StartTime
	data["end_time"] = v.EndTime
	data["status"] = v.Status

	return data
}

func SecKill(req *config.SecRequest) (map[string]interface{}, int, error) {
	//对Map加锁处理
	config.SecKillConfCtx.RWSecProductLock.RLock()
	defer config.SecKillConfCtx.RWSecProductLock.RUnlock()

	var code int
	err := srv_limit.UserCheck(req)
	if err != nil {
		code = srv_err.ErrUserCheckAuthFailed
		log.Printf("userId[%d] invalid, check failed, req[%v]", req.UserId, req)
		return nil, code, err
	}

	err = srv_limit.AntiSpam(req)
	if err != nil {
		code = srv_err.ErrUserServiceBusy
		log.Printf("userId[%d] invalid, check failed, req[%v]", req.UserId, req)
		return nil, code, err
	}

	data, code, err := SecInfoById(req.ProductId)
	if err != nil {
		log.Printf("userId[%d] secInfoById Id failed, req[%v]", req.UserId, req)
		return nil, code, err
	}

	userKey := fmt.Sprintf("%d_%d", req.UserId, req.ProductId)
	fmt.Println("userKey : ", userKey)
	config.SecKillConfCtx.UserConnMap[userKey] = req.ResultChan
	//将请求送入通道并推入到redis队列当中
	config.SecKillConfCtx.SecReqChan <- req

	ticker := time.NewTicker(time.Second * 10)

	defer func() {
		ticker.Stop()
		config.SecKillConfCtx.UserConnMapLock.Lock()
		delete(config.SecKillConfCtx.UserConnMap, userKey)
		config.SecKillConfCtx.UserConnMapLock.Unlock()
	}()

	select {
	case <-ticker.C:
		code = srv_err.ErrProcessTimeout
		err = fmt.Errorf("request timeout")
		return nil, code, err
	case <-req.CloseNotify:
		code = srv_err.ErrClientClosed
		err = fmt.Errorf("client already closed")
		return nil, code, err
	case result := <-req.ResultChan:
		code = result.Code
		if code != 1002 {
			return data, code, srv_err.GetErrMsg(code)
		}
		data["product_id"] = result.ProductId
		data["token"] = result.Token
		data["user_id"] = result.UserId
		return data, code, nil
	}
}

func NewSecRequest() *config.SecRequest {
	secRequest := &config.SecRequest{
		ResultChan: make(chan *config.SecResult, 1),
	}
	return secRequest
}

func SecInfoList() ([]map[string]interface{}, int, error) {
	config.SecKillConfCtx.RWSecProductLock.RLock()
	defer config.SecKillConfCtx.RWSecProductLock.RUnlock()

	var data []map[string]interface{}
	for _, v := range config.SecKillConfCtx.SecProductInfoMap {
		item, _, err := SecInfoById(v.ProductId)
		if err != nil {
			log.Printf("get sec info, err : %v", err)
			continue
		}
		data = append(data, item)
	}
	return data, 0, nil
}

func SecInfoById(productId int) (map[string]interface{}, int, error) {
	//对Map加锁处理
	config.SecKillConfCtx.RWSecProductLock.RLock()
	defer config.SecKillConfCtx.RWSecProductLock.RUnlock()

	var code int
	v, ok := config.SecKillConfCtx.SecProductInfoMap[productId]
	if !ok {
		return nil, srv_err.ErrNotFoundProductId, fmt.Errorf("not found product_id:%d", productId)
	}

	start := false      //秒杀活动是否开始
	end := false        //秒杀活动是否结束
	status := "success" //状态

	nowTime := time.Now().Unix()

	//秒杀活动没有开始
	if nowTime-v.StartTime < 0 {
		start = false
		end = false
		status = "second kill not start"
		code = srv_err.ErrActiveNotStart
	}

	//秒杀活动已经开始
	if nowTime-v.StartTime > 0 {
		start = true
	}

	//秒杀活动已经结束
	if nowTime-v.EndTime > 0 {
		start = false
		end = true
		status = "second kill is already end"
		code = srv_err.ErrActiveAlreadyEnd
	}

	//商品已经被停止或售磬
	if v.Status == config.ProductStatusForceSaleOut || v.Status == config.ProductStatusSaleOut {
		start = false
		end = false
		status = "product is sale out"
		code = srv_err.ErrActiveSaleOut
	}

	//组装数据
	data := map[string]interface{}{
		"product_id": productId,
		"start":      start,
		"end":        end,
		"status":     status,
	}
	return data, code, nil
}
