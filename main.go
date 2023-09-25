package main

import (
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"go.uber.org/zap"
	"xm-task/api"
	"xm-task/conf"
	"xm-task/dao"
	"xm-task/helpers/modules"
	"xm-task/log"
	"xm-task/services"
)

func main() {
	var configPath string
	var migrate string
	flag.StringVar(&configPath, "config", "./config.json", "Path to the config file")
	flag.StringVar(&migrate, "migrate", "true", "Should backend make migration?")
	flag.Parse()

	if configPath == "" {
		log.Info("Usage: program_name -config <config_path>")
		flag.PrintDefaults()
		return
	}

	if migrate == "" {
		log.Info("Usage: program_name -migrate <boolean>")
		flag.PrintDefaults()
		return
	}

	bMigrate, err := strconv.ParseBool(migrate)
	if err != nil {
		log.Fatal("cannot parse bool param", zap.Error(err))
	}

	config, err := conf.GetNewConfig(configPath)
	if err != nil {
		log.Fatal("cannot read config from file", zap.Error(err))
	}

	d, err := dao.New(config, bMigrate)
	if err != nil {
		log.Fatal("dao.New", zap.Error(err))
	}

	s, err := services.NewService(config, d)
	if err != nil {
		log.Fatal("services.NewService", zap.Error(err))
	}

	a, err := api.NewAPI(config, s)
	if err != nil {
		log.Fatal("api.NewAPI", zap.Error(err))
	}

	mds := []modules.Module{a}

	modules.Run(mds)

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	<-gracefulStop
	modules.Stop(mds)
}
