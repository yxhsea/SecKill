package srv_err

import (
	"errors"
)

const (
	ErrInvalidRequest      = 1001
	ErrNotFoundProductId   = 1002
	ErrUserCheckAuthFailed = 1003
	ErrUserServiceBusy     = 1004
	ErrActiveNotStart      = 1005
	ErrActiveAlreadyEnd    = 1006
	ErrActiveSaleOut       = 1007
	ErrProcessTimeout      = 1008
	ErrClientClosed        = 1009
)

const (
	ErrServiceBusy     = 1001
	ErrSecKillSucc     = 1002
	ErrNotFoundProduct = 1003
	ErrSoldout         = 1004
	ErrRetry           = 1005
	ErrAlreadyBuy      = 1006
)

var errMsg = map[int]string{
	ErrServiceBusy:     "服务器错误",
	ErrSecKillSucc:     "抢购成功",
	ErrNotFoundProduct: "没有该商品",
	ErrSoldout:         "商品售罄",
	ErrRetry:           "请重试",
	ErrAlreadyBuy:      "已经抢购",
}

func GetErrMsg(code int) error {
	return errors.New(errMsg[code])
}
