package config

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/go-redis/redis"
	"sync"
)

const (
	ProductStatusNormal       = 0 //商品状态正常
	ProductStatusSaleOut      = 1 //商品售罄
	ProductStatusForceSaleOut = 2 //商品强制售罄
)

var SecKillConfCtx = &SecKillConf{
	SecProductInfoMap: make(map[int]*SecProductInfoConf, 1024),
	UserConnMap:       make(map[string]chan *SecResult, 1024),
	SecReqChan:        make(chan *SecRequest, 1024),
}

//redis配置
type RedisConf struct {
	RedisConn            *redis.Client //链接
	Proxy2layerQueueName string        //队列名称
	Layer2proxyQueueName string        //队列名称
	IdBlackListHash      string        //用户黑名单hash表
	IpBlackListHash      string        //IP黑名单Hash表
	IdBlackListQueue     string        //用户黑名单队列
	IpBlackListQueue     string        //IP黑名单队列
}

//Etcd配置
type EtcdConf struct {
	EtcdConn          *clientv3.Client //链接
	EtcdSecProductKey string           //商品键
}

//访问限制
type AccessLimitConf struct {
	IPSecAccessLimit   int //IP每秒钟访问限制
	UserSecAccessLimit int //用户每秒钟访问限制
	IPMinAccessLimit   int //IP每分钟访问限制
	UserMinAccessLimit int //用户每分钟访问限制
}

type SecKillConf struct {
	RedisConf *RedisConf
	EtcdConf  *EtcdConf

	SecProductInfoMap map[int]*SecProductInfoConf
	RWSecProductLock  sync.RWMutex

	CookieSecretKey string

	ReferWhiteList []string //白名单

	IPBlackMap map[string]bool
	IDBlackMap map[int]bool

	AccessLimitConf AccessLimitConf

	RWBlackLock                  sync.RWMutex
	WriteProxy2LayerGoroutineNum int
	ReadProxy2LayerGoroutineNum  int

	SecReqChan     chan *SecRequest
	SecReqChanSize int

	UserConnMap     map[string]chan *SecResult
	UserConnMapLock sync.Mutex
}

//商品信息配置
type SecProductInfoConf struct {
	ProductId int   `json:"product_id"` //商品ID
	StartTime int64 `json:"start_time"` //开始时间
	EndTime   int64 `json:"end_time"`   //结束时间
	Status    int   `json:"status"`     //状态
	Total     int   `json:"total"`      //商品总数量
	Left      int   `json:"left"`       //商品剩余数量
}

type SecResult struct {
	ProductId int    `json:"product_id"` //商品ID
	UserId    int    `json:"user_id"`    //用户ID
	Token     string `json:"token"`      //Token
	TokenTime int64  `json:"token_time"` //Token生成时间
	Code      int    `json:"code"`       //状态码
}

type SecRequest struct {
	ProductId     int             `json:"product_id"` //商品ID
	Source        string          `json:"source"`
	AuthCode      string          `json:"auth_code"`
	SecTime       string          `json:"sec_time"`
	Nance         string          `json:"nance"`
	UserId        int             `json:"user_id"`
	UserAuthSign  string          `json:"user_auth_sign"` //用户授权签名
	AccessTime    int64           `json:"access_time"`
	ClientAddr    string          `json:"client_addr"`
	ClientRefence string          `json:"client_refence"`
	CloseNotify   <-chan bool     `json:"-"`
	ResultChan    chan *SecResult `json:"-"`
}
