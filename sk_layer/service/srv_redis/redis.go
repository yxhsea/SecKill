package srv_redis

import (
	"SecKill/sk_layer/config"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func RunProcess() {
	for i := 0; i < config.SecLayerCtx.SecLayerConf.ReadGoroutineNum; i++ {
		config.SecLayerCtx.WaitGroup.Add(1)
		go HandleReader()
	}

	for i := 0; i < config.SecLayerCtx.SecLayerConf.WriteGoroutineNum; i++ {
		config.SecLayerCtx.WaitGroup.Add(1)
		go HandleWrite()
	}

	for i := 0; i < config.SecLayerCtx.SecLayerConf.HandleUserGoroutineNum; i++ {
		config.SecLayerCtx.WaitGroup.Add(1)
		go HandleUser()
	}

	log.Printf("all process goroutine started")
	config.SecLayerCtx.WaitGroup.Wait()
	log.Printf("wait all goroutine exited")
	return
}

func HandleReader() {
	log.Println("read goroutine running")
	for {
		conn := config.SecLayerCtx.RedisConf.RedisConn
		for {
			//从Redis队列中取出数据
			data, err := conn.BRPop(time.Minute, config.SecLayerCtx.RedisConf.Proxy2layerQueueName).Result()
			if err != nil {
				log.Printf("blpop from data failed, err : %v", err)
				continue
			}
			log.Printf("brpop from proxy to layer queue, data : %s\n", data)

			//转换数据结构
			var req config.SecRequest
			err = json.Unmarshal([]byte(data[1]), &req)
			if err != nil {
				log.Printf("unmarshal to secrequest failed, err : %v", err)
				continue
			}

			//判断是否超时
			nowTime := time.Now().Unix()
			fmt.Println(nowTime, " ", req.AccessTime, " ", int64(config.SecLayerCtx.SecLayerConf.MaxRequestWaitTimeout))
			if nowTime-req.AccessTime >= int64(config.SecLayerCtx.SecLayerConf.MaxRequestWaitTimeout) {
				log.Printf("req[%v] is expire", req)
				continue
			}

			//设置超时时间
			timer := time.NewTicker(time.Millisecond * time.Duration(config.SecLayerCtx.SecLayerConf.SendToHandleChanTimeout))
			select {
			case config.SecLayerCtx.Read2HandleChan <- &req:
			case <-timer.C:
				log.Printf("send to handle chan timeout, req : %v", req)
				break
			}
		}
	}
}

func HandleWrite() {
	log.Println("handle write running")

	for res := range config.SecLayerCtx.Handle2WriteChan {
		fmt.Println("===", res)
		err := sendToRedis(res)
		if err != nil {
			log.Printf("send to redis, err : %v, res : %v", err, res)
			continue
		}
	}
}

//将数据推入到Redis队列
func sendToRedis(res *config.SecResponse) (err error) {
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("marshal failed, err : %v", err)
		return
	}

	fmt.Println("推入队列前~~")
	conn := config.SecLayerCtx.RedisConf.RedisConn
	err = conn.LPush(config.SecLayerCtx.RedisConf.Layer2proxyQueueName, string(data)).Err()
	fmt.Println("推入队列后~~")
	if err != nil {
		log.Printf("rpush layer to proxy redis queue failed, err : %v", err)
		return
	}
	log.Printf("lpush layer to proxy success. data[%v]", string(data))

	return
}
