package setup

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"testing"
	"time"
)

func TestInitEtcd(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.199.159:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Printf("Connect etcd failed. Error : %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := cli.Get(ctx, "sec_kill_product")
	if err != nil {
		log.Printf("Get falied. Error : %v", err)
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}
