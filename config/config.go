package config

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// Config 配置文件
type Config struct {
	// HTTP
	ServerWin   string `ini:"http_server_win"`
	ServerLinux string `ini:"http_server_linux"`
	// LOG
	LogDirWin   string `ini:"log_dir_win"`
	LogDirLinux string `ini:"log_dir_linux"`
	LogPrefix   string `ini:"log_prefix"`
	// DEBUG
	DebugEnable bool `ini:"debug_enable"`
}

func (c Config) String() string {

	http := fmt.Sprintf("HTTP:[%v]/[%v]", c.ServerWin, c.ServerLinux)

	log := fmt.Sprintf("LOG:[win:%v]/[linux:%v]:[prefix:%v]", c.LogDirWin, c.LogDirLinux, c.LogPrefix)

	return http + ", " + log
}

// ReadConfig 读取配置文件
func ReadConfig(path string) (Config, error) {
	var config Config
	conf, err := ini.Load(path)
	if err != nil {
		log.Println("load config file fail!")
		return config, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&config)
	if err != nil {
		log.Println("mapto config file fail!")
		return config, err
	}
	return config, nil
}
