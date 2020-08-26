package server

import (
	"errors"

	_ "github.com/gfes980615/Diana/transport/http/controller"

	"github.com/gfes980615/Diana/glob/common/log"

	"github.com/gfes980615/Diana/db"
	"github.com/gfes980615/Diana/glob/config"
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/injection/controller"

	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
)

// Run ...
func Run() error {
	stop := make(chan error)

	if config.Config == nil {
		return errors.New("server run fail with nil config")
	} else {
		configInfo := pp.Sprintln(config.Config)
		log.Debug(configInfo)
	}

	// Error handling setup
	log.Init(config.Config.LogConfig.Env,
		config.Config.LogConfig.Level,
		config.Config.LogConfig.HistoryPath,
		config.Config.LogConfig.Duration,
		"",
		"",
		false,
		config.Config.LogConfig.FullColor,
		config.Config.LogConfig.FullTimestamp)

	if err := db.InitMysql(config.Config.DatabaseConfig.Mysql); err != nil {
		return err
	}

	if err := injection.InitInject(); err != nil {
		return err
	}

	if router, err := controller.InitController(); err != nil {
		return err
	} else {
		go run(router, stop)
	}

	log.Info("Server start.")
	err := <-stop
	log.Stop()

	return err
}

func run(router *gin.Engine, stop chan error) {
	if err := router.Run(config.Config.GINConfig.Address); err != nil {
		stop <- errors.New(" Doesn't has valid port. ")
	}
}
