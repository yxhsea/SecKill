package setup

import (
	"SecKill/sk_layer/config"
	"SecKill/sk_layer/logic"
	"github.com/coreos/etcd/clientv3"
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

	config.SecLayerCtx.EtcdConf = &config.EtcdConf{
		EtcdConn:          cli,
		EtcdSecProductKey: productKey,
	}

	logic.LoadProductFromEtcd()
}
