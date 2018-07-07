package setup

import (
	"SecKill/sk_admin/config"
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

	config.SecAdminConfCtx.EtcdConf = &config.EtcdConf{
		EtcdConn:          cli,
		EtcdSecProductKey: productKey,
	}
}
