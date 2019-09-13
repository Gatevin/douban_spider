package config

import (
    "github.com/spf13/viper"
)

var conf *viper.Viper = nil

func LoadConfig() *viper.Viper {
    if conf == nil {
        v := viper.New()
        v.AddConfigPath("./conf/")
        v.SetConfigType("yaml")
        v.SetConfigName("config")
        if err := v.ReadInConfig(); err != nil {
            return nil
        }
        conf = v
    }
    return conf
}
