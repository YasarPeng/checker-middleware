package verify

type Result struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type RDBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"db"`
	Driver   string `json:"driver"` // 如 "mysql"、"dm"、"goldendb"
}

type RDBResult struct {
	Connect map[string]string `json:"connect"`
	Write   map[string]string `json:"write"`
	Delete  map[string]string `json:"delete"`
}

type RDBConnection struct {
	Driver   string `json:"driver"`
	DSN      string `json:"dsn"`
	TimeoutS int    `json:"timeout"`
}

type CacheConfig struct {
	Host      string   `json:"host"`
	Port      int      `json:"port"`
	Password  string   `json:"password"`
	DB        int      `json:"db"`
	Mode      string   // "redis", "credis", "sentinel"
	Sentinels []string // sentinel 地址列表
	Master    string   // sentinel 主名
	Timeout   int      `json:"timeout"`
}

type CacheResult struct {
	Connect map[string]string `json:"connect"`
	Write   map[string]string `json:"write"`
	Delete  map[string]string `json:"delete"`
}

type StorageConfig struct {
	Endpoint     string
	AccessKey    string
	SecretKey    string
	Bucket       string
	Region       string // s3/oss用
	Provider     string // minio, oss, s3
	Secure       bool   // minio用
	UsePathStyle bool   // s3用
	Timeout      int    // 秒
}

type StorageResult struct {
	Connect map[string]string `json:"connect"`
	Write   map[string]string `json:"write"`
	Delete  map[string]string `json:"delete"`
}

type MQConfig struct {
	Provider string // "kafka" or "rabbitmq"
	// Kafka
	Brokers []string
	Topic   string
	// RabbitMQ
	Host     string
	Port     int
	User     string
	Password string
	Vhost    string
}

type MQResult struct {
	Connect map[string]string `json:"connect"`
	Write   map[string]string `json:"write"`
	Delete  map[string]string `json:"delete"`
}
