package main

import (
	"github.com/spf13/viper"
	"gladiators/src/ProxyController"
	"log"
	"os"
	"path/filepath"
)

func main() {
	connectionType := viper.GetString("proxy.connectionType")
	proxyHost := viper.GetString("proxy.host")
	proxyPort := viper.GetString("proxy.port")
	localHost := viper.GetString("local.host")
	localPort := viper.GetString("local.port")

	pm := ProxyController.NewProxyManager(connectionType, proxyHost, proxyPort, localHost, localPort)

	if err := pm.Listen(); err != nil {
		log.Print(err)
		panic(err)
	}

}

func init() {

	if err := initConfig(); err != nil {
		log.Print(err)
		panic(err)
	}

}

func initConfig() error {
	curDir, err := os.Getwd()
	if err != nil {
		return err
	}
	parentPath := filepath.Dir(curDir)
	configPath := filepath.Join(parentPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil

}
