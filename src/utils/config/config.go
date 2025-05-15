package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//@Desc	动态加载YAML配置文件
func InitConfig() error {
	// config/config.yaml
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}
	
	// 解析配置
	for key, value := range(ConfigMap) {
		if err := viper.UnmarshalKey(key, value); err != nil {
			return fmt.Errorf("解析%s配置失败: %v", key, err)
		}
	}
	
	// 监视配置文件改动
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("检测到配置文件变更:", e.Name)
		for key, value := range(ConfigMap) {
			_ = viper.UnmarshalKey(key, value)
		}
	})

	return nil
}
