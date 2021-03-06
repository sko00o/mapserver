//go:generate go-bindata -prefix ./public/ -o=assets/public/public_gen.go -pkg=public public/...
//go:generate go-bindata -prefix ./templates/ -o=assets/templates/templates_gen.go -pkg=templates templates/...

package main

import (
	"fmt"
	"mapserver/applog"
	"mapserver/config"
	"mapserver/httpserver"
	"runtime"

	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

// package flags
var (
	TimeStamp string
	Version   string
	Revision  string
)

func run(c *cli.Context) error {
	conf, err := config.ReadConfig(c.String("conf"))
	if err != nil {
		log.Error("read from conf fail!", c.String("conf"))
		return err
	}
	log.Infoln("conf = ", conf)
	log.Infoln("runtime.GOOS = ", runtime.GOOS)

	var logger *applog.AutoDailyLoger
	if runtime.GOOS == "windows" {
		logger = applog.NewAutoDailyLoger(conf.LogDirWin, conf.LogPrefix)
	} else {
		logger = applog.NewAutoDailyLoger(conf.LogDirLinux, conf.LogPrefix)
	}

	logger.Start()
	defer logger.Stop()

	//start http server
	go func() {
		if runtime.GOOS == "windows" {
			httpserver.StartServer(conf.ServerWin, conf.DebugEnable)
		} else {
			httpserver.StartServer(conf.ServerLinux, conf.DebugEnable)
		}
	}()

	//quit when receive end signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.Infof("signal received signal %v", <-sigChan)
	log.Warn("shutting down server")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "mapserver"
	app.Usage = "local map server"
	app.Copyright = "lowan@lowan-cn.com"
	app.Version = fmt.Sprintf("Version @ %s; Revision @ %s; TimeStamp @ %s", Version, Revision, TimeStamp)
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "conf,c",
			Usage:  "Set conf path here",
			Value:  "appserver_dev.conf",
			EnvVar: "APP_CONF",
		},
	}
	app.Run(os.Args)
}
