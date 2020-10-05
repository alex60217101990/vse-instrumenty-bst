package logger

import (
	"log"
	"os"
	"sync"

	"github.com/alex60217101990/vse-instrumenty-bst/external/configs"

	"github.com/fatih/color"
)

var (
	PackageOnceLoad sync.Once

	CmdServer = color.New(color.FgHiGreen, color.Bold)
	CmdError  = color.New(color.FgHiRed, color.Bold)
	CmdInfo   = color.New(color.FgHiBlue, color.Bold)

	AppLogger     Logger
	DefaultLogger *log.Logger = log.New(os.Stderr, color.RedString("[SYSTEM] "), log.LstdFlags|log.Lshortfile|log.LUTC)
)

func InitLoggerSettings() {
	PackageOnceLoad.Do(func() {
		CmdServer.Println("Run once - 'logger' package loading...")

		switch configs.Conf.Logger.LoggerType {
		case configs.Zero:
			AppLogger = NewZeroLogger()
		case configs.Zap:
			AppLogger = NewZapLogger()
		default:
			log.Fatalf("invalid logger type: %+v\n", configs.Conf.Logger.LoggerType)
		}
	})
}

func CloseLoggers() {
	AppLogger.Close()

	CmdServer.Println("'logger' package stoped...")
}
