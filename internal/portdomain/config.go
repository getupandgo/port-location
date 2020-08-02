package portdomain

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
