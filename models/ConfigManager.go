package models

import (
	"fmt"

	"github.com/spf13/viper"
)

type ConfigManager struct{}

type ConfigManagerStruct struct {
	DbConf  ConfigDbStruct  `mapstructure:"db"`
	AppConf ConfigAppStruct `mapstructure:"app"`
}

type ConfigDbStruct struct {
	Address  string
	Port     int
	Name     string
	User     string
	Password string
}

type ConfigAppStruct struct {
	Salt           string
	CookieLifetime int
	Port           int
	FullName       string
	ShortName      string
	Titile         string
}

var configVariables ConfigManagerStruct
var v *viper.Viper

func (p *ConfigManager) CheckConfig() bool {
	v = viper.New()
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(".")
	v.AddConfigPath("./")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		return false
	}
	err := v.Unmarshal(&configVariables)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
		return false
	}
	return true
}

func (p *ConfigManager) GetProps() *ConfigManagerStruct {
	return &configVariables
}
