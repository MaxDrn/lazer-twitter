package options

var Current = ConfigWithDefaults()

func ConfigWithDefaults() *Config {
	return &Config{
		HostName:       "localhost",
		RESTListenPort: "8080",
		DBName:         "postgres",
		DBUser:         "postgres",
		DBPassword:     "password",
		DBHost:         "localhost",
		DBPort:         "5432",
	}
}

type Config struct {
	HostName       string
	RESTListenPort string
	DBName         string
	DBUser         string
	DBPassword     string
	DBHost         string
	DBPort         string
}
