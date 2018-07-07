package setup

import (
	"SecKill/sk_layer/config"
	"github.com/go-redis/redis"
	"log"
)

//初始化redis
func InitRedis(host string, passWord string, db int, proxy2layerQueueName, layer2proxyQueueName string) {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: passWord,
		DB:       db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Printf("Connect redis failed. Error : %v", err)
	}

	config.SecLayerCtx.RedisConf = &config.RedisConf{
		RedisConn:            client,
		Proxy2layerQueueName: proxy2layerQueueName,
		Layer2proxyQueueName: layer2proxyQueueName,
	}
}
