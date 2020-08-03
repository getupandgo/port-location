package portdomain

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type GRPCServer struct {
	Host string `yaml:"host" validate:"required"`
	Port string `yaml:"port" validate:"required"`
}

type DB struct {
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	Name          string `yaml:"name"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	MigrationsDir string `yaml:"migrations_dir"`
}

type Config struct {
	GRPCServer `yaml:"grpc_server" validate:"required"`
	DB         `yaml:"db" validate:"required"`
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
