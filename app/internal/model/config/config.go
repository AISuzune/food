package config

type Config struct {
	App        *App      `mapstructure:"app"  yaml:"app"`
	DataBase   *Database `mapstructure:"database"  yaml:"database"`
	Logger     *Logger   `mapstructure:"logger" yaml:"logger"`
	Server     *Server   `mapstructure:"server"  yaml:"server"`
	Cors       CORS      `mapstructure:"cors" yaml:"cors"`
	Auth       Auth      `mapstructure:"auth" yaml:"auth"`
	YelpApiKey string    `mapstructure:"yelpApiKey" yaml:"yelpApiKey"`
}
