package helpers

import (
	"fmt"
	"os"

	"github.com/alex60217101990/vse-instrumenty-bst/external/configs"
	"github.com/alex60217101990/vse-instrumenty-bst/external/logger"
	"github.com/integralist/go-findroot/find"
)

func InitConfigs(path string) {
	if len(path) == 0 {
		root, err := find.Repo()
		if err != nil {
			logger.DefaultLogger.Fatal(err)
		}

		path = fmt.Sprintf("%s%sdeploy%sconfigs%sapp-configs.yaml", root.Path, string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator))

		logger.CmdInfo.Printf("root project directory: %+v; configs dir: %s\n", root, path)
	}

	err := configs.ReadConfigFile(path)
	if err != nil {
		logger.DefaultLogger.Fatal(err)
	}
}
