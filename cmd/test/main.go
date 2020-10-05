package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alex60217101990/vse-instrumenty-bst/external/configs"
	"github.com/alex60217101990/vse-instrumenty-bst/external/helpers"
	fast_http "github.com/alex60217101990/vse-instrumenty-bst/external/http-server/fast-http"
	"github.com/alex60217101990/vse-instrumenty-bst/external/logger"
)

var (
	debug      = flag.Bool("d", true, "Print debug logs ?")
	confFile   = flag.String("c", "", "Config file path")
	loggerType configs.LoggerType
)

func main() {
	flag.Var(&loggerType, "lt", "Type of logger usage (zap | zero)")
	flag.Usage = helpers.PrintFlags
	flag.Parse()

	helpers.InitConfigs(helpers.StringPtr(confFile))
	configs.Conf.IsDebug = *debug
	if loggerType != configs.Default {
		configs.Conf.Logger.LoggerType = loggerType
	}

	logger.InitLoggerSettings()

	server := fast_http.NewFastHttpServer()
	server.Init()
	go server.Run()

	logger.CmdServer.Printf("🚀 %s service started...\n", strings.ToUpper(configs.Conf.ServiceName))

	var Stop = make(chan os.Signal, 1)
	signal.Notify(Stop,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGABRT,
	)
	for range Stop {
		server.Close()
		logger.CmdServer.Printf("🚫 %s service stoped...\n", strings.ToUpper(configs.Conf.ServiceName))
		logger.CloseLoggers()
		return
	}
}
