package clientapi

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
