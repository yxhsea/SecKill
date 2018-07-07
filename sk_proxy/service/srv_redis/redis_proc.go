package srv_redis

import (
	"SecKill/sk_proxy/config"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

//写数据到Redis
func WriteHandle() {
	for {
		fmt.Println("wirter data to redis.")
		req := <-config.SecKillConfCtx.SecReqChan
		fmt.Println("accessTime : ", req.AccessTime)
		conn := config.SecKillConfCtx.RedisConf.RedisConn

		data, err := json.Marshal(req)
		if err != nil {
			log.Printf("json.Marshal req failed. Error : %v, req : %v", err, req)
			continue
		}

		err = conn.LPush(config.SecKillConfCtx.RedisConf.Proxy2layerQueueName, string(data)).Err()
		if err != nil {
			log.Printf("lpush req failed. Error : %v, req : %v", err, req)
			continue
		}
		log.Printf("lpush req success. req : %v", string(data))
	}
}

//从redis读取数据
func ReadHandle() {
	for {
		conn := config.SecKillConfCtx.RedisConf.RedisConn

		//阻塞弹出
		data, err := conn.BRPop(time.Minute, config.SecKillConfCtx.RedisConf.Layer2proxyQueueName).Result()
		if err != nil {
			log.Printf("brpop layer2proxy failed. Error : %v", err)
			continue
		}

		var result *config.SecResult
		err = json.Unmarshal([]byte(data[1]), &result)
		if err != nil {
			log.Printf("json.Unmarshal failed. Error : %v", err)
			continue
		}

		userKey := fmt.Sprintf("%d_%d", result.UserId, result.ProductId)
		fmt.Println("userKey : ", userKey)
		config.SecKillConfCtx.UserConnMapLock.Lock()
		resultChan, ok := config.SecKillConfCtx.UserConnMap[userKey]
		config.SecKillConfCtx.UserConnMapLock.Unlock()
		if !ok {
			log.Printf("user not found : %v", userKey)
			continue
		}

		resultChan <- result
		log.Printf("request result send to chan succeee, userKey : %v", userKey)
	}
}
