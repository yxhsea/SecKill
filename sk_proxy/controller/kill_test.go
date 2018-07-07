package controller

import (
	"crypto/md5"
	"fmt"
	"testing"
	"time"
)

func TestKill(t *testing.T) {
	data := map[string]interface{}{
		"product_id": 1,
		"user_id":    1,
		"src":        "192.168.199.1",
		"auth_code":  "userauthcode",
		"time":       time.Now().Unix(),
		"nance":      "dsdsdjkdjskdjksdjhuieurierei",
	}
	authData := fmt.Sprintf("%d:%s", data["user_id"], "WK5wJOiuYaXRUlPsxo3LZEbpCNSyvm8T")
	authSign := fmt.Sprintf("%x", md5.Sum([]byte(authData)))

	fmt.Println("product_id:", data["product_id"])
	fmt.Println("user_id:", data["user_id"])
	fmt.Println("src:", data["src"])
	fmt.Println("auth_code:", data["auth_code"])
	fmt.Println("time:", data["time"])
	fmt.Println("nance:", data["nance"])
	fmt.Println("==============================")
	fmt.Println("AuthSign:", authSign)
}
