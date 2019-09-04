package options


var Current = ConfigWithDefaults()

func ConfigWithDefaults() *Config {
	return &Config{
		RESTListenPort:     "8080",
		DBName: 			"postgres",
		DBUser: 			"postgres",
		DBPassword: 		"password",
		DBHost:				"localhost",
		DBPort: 			"5432",
	}
}

type Config struct {
	RESTListenPort string
	DBName         string
	DBUser         string
	DBPassword     string
	DBHost         string
	DBPort         string
}