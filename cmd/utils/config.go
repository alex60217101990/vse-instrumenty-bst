package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/alex60217101990/vse-instrumenty-bst/external/configs"
	"github.com/alex60217101990/vse-instrumenty-bst/external/enums"
	"github.com/alex60217101990/vse-instrumenty-bst/external/helpers"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

var format = flag.String("format", "yaml", "Set configs file format. (json, yaml)")

var green = color.New(color.FgGreen, color.Bold)

func main() {
	var err error

	flag.Usage = helpers.PrintFlags
	flag.Parse()

	var currConfigFormat enums.EncodingType
	currConfigFormat.FromString(helpers.StringPtr(format))

	conf := &configs.Configs{
		Ver:         helpers.String("0.0.1"),
		ServiceName: "bst_tree",
		IsDebug:     true,
		Logger: &configs.Logger{
			LoggerType: configs.Zero,
			LogsPath:   "./tmp/logs",
		},
		Server: &configs.Server{
			Host: "",
			Port: 8077,
		},
		BST: &configs.BST{
			SnapshotPath:   "",
			UseCompression: true,
			Limit:          true,
			Size:           30,
		},
	}

	var bts []byte
	switch currConfigFormat {
	case enums.Yaml:
		bts, err = yaml.Marshal(conf)
	case enums.Json:
		bts, err = json.MarshalIndent(conf, "", "\t")
	default:
		log.Fatal("invalid format")
	}

	if err != nil {
		log.Fatal(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dir = filepath.Join(dir, "deploy", "configs")
	_, _ = green.Println(dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			log.Fatal(err)
		}
	}

	filePath := filepath.Join(dir, fmt.Sprintf("app-configs.%s", currConfigFormat))
	_, _ = green.Println(filePath)
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filePath, bts, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
