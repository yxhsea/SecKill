package main

import (
	"SecKill/sk_proxy/setup"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var cfgFile string
var Verbose bool

func main() {

	var RootCmd = &cobra.Command{
		Use:   "SKProxy Server",
		Short: "SKProxy Server",
		Long:  "SKProxy Server",
		Run: func(cmd *cobra.Command, args []string) {
			serviceMap := viper.GetStringMap("service")
			ipSecAccessLimit, _ := serviceMap["ip_sec_access_limit"].(int64)
			ipMinAccessLimit, _ := serviceMap["ip_min_access_limit"].(int64)
			userSecAccessLimit, _ := serviceMap["user_sec_access_limit"].(int64)
			userMinAccessLimit, _ := serviceMap["user_min_access_limit"].(int64)
			writeProxy2layerGoroutineNum, _ := serviceMap["write_proxy2layer_goroutine_num"].(int64)
			readProxy2layerGoroutineNum, _ := serviceMap["read_proxy2layer_goroutine_num"].(int64)
			cookieSecretKey, _ := serviceMap["cookie_secretkey"].(string)
			referWhitelist, _ := serviceMap["refer_whitelist"].(string)
			setup.InitServiceConfig(ipSecAccessLimit, ipMinAccessLimit, userSecAccessLimit, userMinAccessLimit,
				writeProxy2layerGoroutineNum, readProxy2layerGoroutineNum, cookieSecretKey, referWhitelist)

			redisMap := viper.GetStringMap("redis")
			hostRedis, _ := redisMap["host"].(string)
			pwdRedis, _ := redisMap["password"].(string)
			dbRedis, _ := redisMap["db"].(int)
			proxy2layerQueueNameRedis, _ := redisMap["proxy2layer_queue_name"].(string)
			layer2proxyQueueNameRedis, _ := redisMap["layer2proxy_queue_name"].(string)
			idBlackListHashRedis, _ := redisMap["id_black_list_hash"].(string)
			ipBlackListHashRedis, _ := redisMap["ip_black_list_hash"].(string)
			idBlackListQueueRedis, _ := redisMap["id_black_list_queue"].(string)
			ipBlackListQueueRedis, _ := redisMap["ip_black_list_queue"].(string)
			setup.InitRedis(hostRedis, pwdRedis, dbRedis, proxy2layerQueueNameRedis, layer2proxyQueueNameRedis,
				idBlackListHashRedis, ipBlackListHashRedis, idBlackListQueueRedis, ipBlackListQueueRedis)

			etcdMap := viper.GetStringMap("etcd")
			hostEtcd, _ := etcdMap["host"].(string)
			productKey, _ := etcdMap["product_key"].(string)
			setup.InitEtcd(hostEtcd, productKey)

			httpMap := viper.GetStringMap("http")
			hostHttp, _ := httpMap["host"].(string)
			setup.InitServer(hostHttp)
		},
	}

	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}

func initConfig() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("conf")
		viper.AddConfigPath("./sk_proxy/")
		viper.AddConfigPath(dir)
		viper.AutomaticEnv()
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}
}
