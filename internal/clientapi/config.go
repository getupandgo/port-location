package clientapi

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type GRPCServer struct {
	Host string `yaml:"host" validate:"required"`
	Port string `yaml:"port" validate:"required"`
}

type HTTPServer struct {
	Host string `yaml:"host" validate:"required"`
	Port string `yaml:"port" validate:"required"`
}

type Config struct {
	HTTPServer   `yaml:"http_server" validate:"required"`
	GRPCServer   `yaml:"grpc_server" validate:"required"`
	PortFilePath string `yaml:"port_file_path" validate:"required"`
}

func (c *Config) Read(path string) error {
	if path == "" {
		return errors.New("no path to config file provided")
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	if err = yaml.Unmarshal(b, c); err != nil {
		return errors.Wrap(err, "failed to parse config file")
	}

	return nil
}
