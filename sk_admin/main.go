package main

import (
	"SecKill/sk_admin/setup"
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
		Use:   "SKAdmin Server",
		Short: "SKAdmin Server",
		Long:  "SKAdmin Server",
		Run: func(cmd *cobra.Command, args []string) {
			mysqlMap := viper.GetStringMap("mysql")
			hostMysql, _ := mysqlMap["host"].(string)
			portMysql, _ := mysqlMap["port"].(string)
			userMysql, _ := mysqlMap["user"].(string)
			pwdMysql, _ := mysqlMap["pass_wd"].(string)
			dbMysql, _ := mysqlMap["db"].(string)
			setup.InitMysql(hostMysql, portMysql, userMysql, pwdMysql, dbMysql)

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
		viper.AddConfigPath("./sk_admin/")
		viper.AddConfigPath(dir)
		viper.AutomaticEnv()
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}
}
