package config

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/gohouse/gorose"
)

var SecAdminConfCtx = &SecAdminConf{}

type SecAdminConf struct {
	DbConf   *DbConf
	EtcdConf *EtcdConf
}

//数据库配置
type DbConf struct {
	DbConn gorose.Connection //链接
}

//Etcd配置
type EtcdConf struct {
	EtcdConn          *clientv3.Client //链接
	EtcdSecProductKey string           //商品键
}
