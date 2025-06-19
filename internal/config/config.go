package config

type Config struct {
	DNS          string   `mapstructure:"DNS" json:"DNS" yaml:"DNS"`
	Database     database `mapstructure:"DATABASE" json:"DATABASE" yaml:"DATABASE"`
	Port         string   `mapstructure:"PORT" yaml:"PORT" json:"PORT" default:"3000"`
	PostsPerPage int      `mapstructure:"PORT_PER_PAGE" json:"PORT_PER_PAGE" yaml:"PORT_PER_PAGE" default:"12"`
}

type database struct {
	URL      string `mapstructure:"URL" yaml:"URL"`
	Name     string `mapstructure:"NAME" yaml:"NAME"`
	User     string `mapstructure:"USER" yaml:"USER"`
	Password string `mapstructure:"PASSWORD" yaml:"PASSWORD"`
}
