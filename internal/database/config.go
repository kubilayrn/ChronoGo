package database

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (c *Config) GetDSN() string {
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}
	if c.Port == "" {
		c.Port = "5432"
	}
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/" + c.DBName + "?sslmode=" + c.SSLMode
}
