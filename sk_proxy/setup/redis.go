package setup

import (
	"SecKill/sk_proxy/config"
	"SecKill/sk_proxy/service/srv_redis"
	"github.com/Unknwon/com"
	"github.com/go-redis/redis"
	"log"
	"time"
)

//初始化Redis
func InitRedis(host string, passWord string, db int, proxy2layerQueueNameRedis, layer2proxyQueueNameRedis,
	idBlackListHashRedis, ipBlackListHashRedis, idBlackListQueueRedis, ipBlackListQueueRedis string) {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: passWord,
		DB:       db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Printf("Connect redis failed. Error : %v", err)
	}

	config.SecKillConfCtx.RedisConf = &config.RedisConf{
		RedisConn:            client,
		Proxy2layerQueueName: proxy2layerQueueNameRedis,
		Layer2proxyQueueName: layer2proxyQueueNameRedis,
		IdBlackListHash:      idBlackListHashRedis,
		IpBlackListHash:      ipBlackListHashRedis,
		IdBlackListQueue:     idBlackListQueueRedis,
		IpBlackListQueue:     ipBlackListQueueRedis,
	}
	loadBlackList(client)
	initRedisProcess()
}

//加载黑名单列表
func loadBlackList(conn *redis.Client) {
	config.SecKillConfCtx.IPBlackMap = make(map[string]bool, 10000)
	config.SecKillConfCtx.IDBlackMap = make(map[int]bool, 10000)

	//用户Id
	idList, err := conn.HGetAll(config.SecKillConfCtx.RedisConf.IdBlackListHash).Result()
	if err != nil {
		log.Printf("hget all failed. Error : %v", err)
		return
	}

	for _, v := range idList {
		id, err := com.StrTo(v).Int()
		if err != nil {
			log.Printf("invalid user id [%v]", id)
			continue
		}
		config.SecKillConfCtx.IDBlackMap[id] = true
	}

	//用户Ip
	ipList, err := conn.HGetAll(config.SecKillConfCtx.RedisConf.IpBlackListHash).Result()
	if err != nil {
		log.Printf("hget all failed. Error : %v", err)
		return
	}

	for _, v := range ipList {
		config.SecKillConfCtx.IPBlackMap[v] = true
	}

	go syncIpBlackList(conn)
	go syncIdBlackList(conn)
	return
}

//同步用户ID黑名单
func syncIdBlackList(conn *redis.Client) {
	for {
		idArr, err := conn.BRPop(time.Minute, config.SecKillConfCtx.RedisConf.IdBlackListQueue).Result()
		if err != nil {
			log.Printf("brpop id failed, err : %v", err)
			continue
		}
		id, _ := com.StrTo(idArr[1]).Int()
		config.SecKillConfCtx.RWBlackLock.Lock()
		{
			config.SecKillConfCtx.IDBlackMap[id] = true
		}
		config.SecKillConfCtx.RWBlackLock.Unlock()

	}
}

//同步用户IP黑名单
func syncIpBlackList(conn *redis.Client) {
	var ipList []string
	lastTime := time.Now().Unix()

	for {
		ipArr, err := conn.BRPop(time.Minute, config.SecKillConfCtx.RedisConf.IpBlackListQueue).Result()
		if err != nil {
			log.Printf("brpop ip failed, err : %v", err)
			continue
		}

		ip := ipArr[1]
		curTime := time.Now().Unix()
		ipList = append(ipList, ip)

		if len(ipList) > 100 || curTime-lastTime > 5 {
			config.SecKillConfCtx.RWBlackLock.Lock()
			{
				for _, v := range ipList {
					config.SecKillConfCtx.IPBlackMap[v] = true
				}
			}
			config.SecKillConfCtx.RWBlackLock.Unlock()

			lastTime = curTime
			log.Printf("sync ip list from redis success, ip[%v]", ipList)
		}
	}
}

//初始化redis进程
func initRedisProcess() {
	for i := 0; i < config.SecKillConfCtx.WriteProxy2LayerGoroutineNum; i++ {
		go srv_redis.WriteHandle()
	}

	for i := 0; i < config.SecKillConfCtx.ReadProxy2LayerGoroutineNum; i++ {
		go srv_redis.ReadHandle()
	}
}
