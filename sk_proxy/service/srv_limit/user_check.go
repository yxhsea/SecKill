package srv_limit

import (
	"SecKill/sk_proxy/config"
	"crypto/md5"
	"fmt"
	"log"
)

//用户检查
func UserCheck(req *config.SecRequest) (err error) {
	found := false
	for _, refer := range config.SecKillConfCtx.ReferWhiteList {
		if refer == req.ClientRefence {
			found = true
			break
		}
	}

	if !found {
		err = fmt.Errorf("invalid request")
		log.Printf("user[%d] is reject by refer, req[%v]", req.UserId, req)
		return
	}

	authData := fmt.Sprintf("%d:%s", req.UserId, config.SecKillConfCtx.CookieSecretKey)
	authSign := fmt.Sprintf("%x", md5.Sum([]byte(authData)))

	if authSign != req.UserAuthSign {
		err = fmt.Errorf("invalid user cookie auth")
		return
	}

	return
}
