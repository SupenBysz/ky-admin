package config

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `yaml:"host" json:"host"`
	Port            int    `yaml:"port" json:"port"`
	User            string `yaml:"user" json:"user"`
	Password        string `yaml:"password" json:"password"`
	Name            string `yaml:"name" json:"name"`
	Driver          string `yaml:"driver" json:"driver"`
	Database        string `yaml:"database" json:"database"`
	ShowSQL         bool   `yaml:"show_sql" json:"show_sql"`
	LogLevel        string `yaml:"log_level" json:"log_level"`
	MaxOpenConns    int    `yaml:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns" json:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime" json:"conn_max_lifetime"`
}
