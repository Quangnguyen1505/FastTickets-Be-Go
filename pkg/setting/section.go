package setting

type Config struct {
	Postgresql   PostgresqlSettings `mapstructure:"postgresql"`
	Security     SecuritySettings   `mapstructure:"security"`
	Logger       LoggerSettings     `mapstructure:"logger"`
	Redis        RedisSettings      `mapstructure:"redis"`
	Server       ServerSettings     `mapstructure:"server"`
	JWT          JWTSettings        `mapstructure:"jwt"`
	Oauth2Google Oauth2Google       `mapstructure:"oauth2Google"`
	RabbitMQ     RabbitMQSettings   `mapstructure:"rabbitmq"`
	Grpc         GrpcSettings       `mapstructure:"grpc"`
}

type ServerSettings struct {
	Port   int    `mapstructure:"port"`
	Mode   string `mapstructure:"mode"`
	Domain string `mapstructure:"domainBe"`
}

type RedisSettings struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"pool_size"`
}

type PostgresqlSettings struct {
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Dbname   string `mapstructure:"dbname"`
}

type SecuritySettings struct {
	Jwt struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`
}

type LoggerSettings struct {
	File_name   string `mapstructure:"file_name"`
	Max_size    int    `mapstructure:"max_size"`
	Max_backups int    `mapstructure:"max_backups"`
	Max_age     int    `mapstructure:"max_age"`
	Compress    bool   `mapstructure:"compress"`
	Loglevel    string `mapstructure:"loglevel"`
}

type JWTSettings struct {
	TOKEN_HOUR_LIFESPAN uint   `mapstructure:"TOKEN_HOUR_LIFESPAN"`
	JWT_EXPIRATION      string `mapstructure:"JWT_EXPIRATION"`
	API_SERCERT_KEY     string `mapstructure:"API_SERCERT_KEY"`
}

type Oauth2Google struct {
	CLIENT_ID     string `mapstructure:"clientId"`
	CLIENT_SECRET string `mapstructure:"clientSecret"`
}

type RabbitMQSettings struct {
	Username string `mapstructure:"rabbitmq_username"`
	Password string `mapstructure:"rabbitmq_password"`
	Host     string `mapstructure:"rabbitmq_host"`
	Port     int    `mapstructure:"rabbitmq_port"`
}

type GrpcSettings struct {
	Client struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"client"`
}
