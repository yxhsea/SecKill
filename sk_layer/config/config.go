package config

import (
	"SecKill/sk_layer/service/srv_limit"
	"SecKill/sk_layer/service/srv_product"
	"SecKill/sk_layer/service/srv_user"
	"github.com/coreos/etcd/clientv3"
	"github.com/go-redis/redis"
	"sync"
)

var AppConfig = &SecLayerConf{}
var SecLayerCtx = &SecLayerContext{}

//秒杀商品信息配置
type SecProductInfoConf struct {
	ProductId         int                 `json:"product_id"`           //商品ID
	StartTime         int64               `json:"start_time"`           //秒杀开始时间
	EndTime           int64               `json:"end_time"`             //秒杀结束时间
	Status            int                 `json:"status"`               //状态
	Total             int                 `json:"total"`                //商品总数
	Left              int                 `json:"left"`                 //商品剩余数量
	OnePersonBuyLimit int                 `json:"one_person_buy_limit"` //单个用户购买数量限制
	BuyRate           float64             `json:"buy_rate"`             //购买频率限制
	SoldMaxLimit      int                 `json:"sold_max_limit"`
	SecLimit          *srv_limit.SecLimit `json:"sec_limit"` //限速控制
}

//Redis配置
type RedisConf struct {
	RedisConn            *redis.Client //链接
	Proxy2layerQueueName string        //队列名称
	Layer2proxyQueueName string        //队列名称
}

//Etcd配置
type EtcdConf struct {
	EtcdConn          *clientv3.Client //链接
	EtcdSecProductKey string           //商品键
}

//秒杀逻辑层配置
type SecLayerConf struct {
	WriteGoroutineNum int //写操作goroutine数量控制
	ReadGoroutineNum  int //读操作goroutine数量控制

	HandleUserGoroutineNum int //处理用户goroutine数量控制

	Read2HandleChanSize  int //
	Handle2WriteChanSize int //

	MaxRequestWaitTimeout int //最大请求等待时间

	SendToWriteChanTimeout  int //
	SendToHandleChanTimeout int //

	SecProductInfoMap map[int]*SecProductInfoConf
	TokenPassWd       string //Token
}

type SecLayerContext struct {
	RedisConf *RedisConf
	EtcdConf  *EtcdConf

	RWSecProductLock sync.RWMutex

	SecLayerConf *SecLayerConf

	WaitGroup sync.WaitGroup

	Read2HandleChan  chan *SecRequest
	Handle2WriteChan chan *SecResponse

	HistoryMap     map[int]*srv_user.UserBuyHistory
	HistoryMapLock sync.Mutex

	ProductCountMgr *srv_product.ProductCountMgr //商品计数
}

//请求
type SecRequest struct {
	ProductId     int    `json:"product_id"` //商品ID
	Source        string `json:"source"`
	AuthCode      string `json:"auth_code"`
	SecTime       string `json:"sec_time"`
	Nance         string `json:"nance"`
	UserId        int    `json:"user_id"`
	UserAuthSign  string `json:"user_auth_sign"`
	AccessTime    int64  `json:"access_time"`
	ClientAddr    string `json:"client_addr"`
	ClientRefence string `json:"client_refence"`
}

//响应
type SecResponse struct {
	ProductId int    `json:"product_id"` //商品ID
	UserId    int    `json:"user_id"`    //用户ID
	Token     string `json:"token"`      //Token
	TokenTime int64  `json:"token_time"` //Token生成时间
	Code      int    `json:"code"`       //状态码
}
