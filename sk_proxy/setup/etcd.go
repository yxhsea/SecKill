package setup

import (
	"SecKill/sk_proxy/config"
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"log"
	"time"
)

//初始化Etcd
func InitEtcd(host, productKey string) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{host},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Printf("Connect etcd failed. Error : %v", err)
	}

	config.SecKillConfCtx.EtcdConf = &config.EtcdConf{
		EtcdConn:          cli,
		EtcdSecProductKey: productKey,
	}

	loadSecConf(cli)
	go waterSecProductKey(cli, config.SecKillConfCtx.EtcdConf.EtcdSecProductKey)
}

//加载秒杀商品信息
func loadSecConf(cli *clientv3.Client) {
	rsp, err := cli.Get(context.Background(), config.SecKillConfCtx.EtcdConf.EtcdSecProductKey)
	if err != nil {
		log.Printf("get product info failed, err : %v", err)
		return
	}

	var secProductInfo []*config.SecProductInfoConf
	for _, v := range rsp.Kvs {
		err := json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			log.Printf("unmarshal json failed, err : %v", err)
			return
		}
	}

	updateSecProductInfo(secProductInfo)
}

//监听秒杀商品配置
func waterSecProductKey(cli *clientv3.Client, key string) {
	for {
		rch := cli.Watch(context.Background(), key)
		var secProductInfo []*config.SecProductInfoConf
		var getConfSucc = true

		for wrsp := range rch {
			for _, ev := range wrsp.Events {
				//删除事件
				if ev.Type == mvccpb.DELETE {
					continue
				}

				//更新事件
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &secProductInfo)
					if err != nil {
						getConfSucc = false
						continue
					}
				}
			}

			if getConfSucc {
				updateSecProductInfo(secProductInfo)
			}
		}
	}
}

//更新秒杀商品信息
func updateSecProductInfo(secProductInfo []*config.SecProductInfoConf) {
	tmp := make(map[int]*config.SecProductInfoConf, 1024)
	for _, v := range secProductInfo {
		tmp[v.ProductId] = v
	}

	config.SecKillConfCtx.RWSecProductLock.Lock()
	config.SecKillConfCtx.SecProductInfoMap = tmp
	config.SecKillConfCtx.RWSecProductLock.Unlock()
}
