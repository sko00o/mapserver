package config

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type Config struct {
	//HTTP
	HTTPServerWin   string `ini:"http_server_win"`
	HTTPServerLinux string `ini:"http_server_linux"`
	//LOG
	LogDirWin   string `ini:"log_dir_win"`
	LogDirLinux string `ini:"log_dir_linux"`
	LogPrefix   string `ini:"log_prefix"`
	//MAP
	CenterLat float64 `ini:"center_latitude"`
	CenterLng float64 `ini:"center_longitude"`
	ZoomLevel int     `ini:"zoom_level"`
}

//Read Server's Config Value from "path"
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
