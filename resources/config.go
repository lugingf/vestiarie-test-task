package resources

type Config struct {
	Server  *ServerConfig
	Storage *DataBaseConfig
}

type DataBaseConfig struct {
	Driver   string
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type ServerConfig struct {
	Host string
	Port string
}

func NewConfig() (*Config, error) {
	c := Config{
		Server:  serverConfig(),
		Storage: dbConfig(),
	}

	return &c, nil
}

func serverConfig() *ServerConfig {
	return &ServerConfig{
		Host: "localhost",
		Port: "8080",
	}
}

func dbConfig() *DataBaseConfig {
	return &DataBaseConfig{
		Driver:   "mysql",
		User:     "gotest",
		Password: "gotest",
		Host:     "database",
		Port:     "3306",
		Name:     "local_gotest",
	}
}
