package configs

import (
	"encoding/json"

	"github.com/pkg/errors"
)

var Conf *Configs

type Configs struct {
	Ver         *string `yaml:"ver"`
	ServiceName string  `yaml:"service-name" json:"service_name"`
	IsDebug     bool    `yaml:"-" json:"-"`
	Logger      *Logger `yaml:"logger"`
	Server      *Server `yaml:"http-server" json:"http_server"`
	BST         *BST    `yaml:"bst"`
}

type Server struct {
	Host string `yaml:"server-host" json:"server_host"`
	Port uint16 `yaml:"server-port" json:"server_port"`
}

type BST struct {
	SnapshotPath   string `yaml:"snapshot-path" json:"snapshot_path"`
	UseCompression bool   `yaml:"use-compression" json:"use_compression"`
	Limit          bool   `yaml:"use-limit" json:"use_limit"`
	Size           uint64 `yaml:"size" json:"size"`
}

type Logger struct {
	LoggerType LoggerType `yaml:"logger-type" json:"logger_type"`
	LogsPath   string     `yaml:"logs-path" json:"logs_path"`
}

func (l *Logger) MarshalJSON() ([]byte, error) {
	type alias struct {
		LoggerType string `json:"logger_type"`
		LogsPath   string `json:"logs_path"`
	}
	if l == nil {
		l = &Logger{}
	}
	return json.Marshal(alias{
		LoggerType: l.LoggerType.String(),
		LogsPath:   l.LogsPath,
	})
}

func (l *Logger) UnmarshalJSON(data []byte) (err error) {
	type alias struct {
		LoggerType string `json:"logger_type"`
		LogsPath   string `json:"logs_path"`
	}
	var tmp alias
	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if l == nil {
		l = &Logger{}
	}

	err = l.LoggerType.Set(tmp.LoggerType)
	if err != nil {
		return errors.WithMessagef(err, "failed to parse '%s'", tmp.LoggerType)
	}

	l.LogsPath = tmp.LogsPath

	return nil
}

func (l *Logger) MarshalYAML() (interface{}, error) {
	type alias struct {
		LoggerType string `yaml:"logger-type"`
		LogsPath   string `yaml:"logs-path"`
	}
	if l == nil {
		l = &Logger{}
	}
	return alias{
		LoggerType: l.LoggerType.String(),
		LogsPath:   l.LogsPath,
	}, nil
}

func (l *Logger) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias struct {
		LoggerType string `yaml:"logger-type"`
		LogsPath   string `yaml:"logs-path"`
	}
	var tmp alias
	if err := unmarshal(&tmp); err != nil {
		return err
	}

	if l == nil {
		l = &Logger{}
	}

	err := l.LoggerType.Set(tmp.LoggerType)
	if err != nil {
		return errors.WithMessagef(err, "failed to parse '%s'", tmp.LoggerType)
	}

	l.LogsPath = tmp.LogsPath

	return nil
}
