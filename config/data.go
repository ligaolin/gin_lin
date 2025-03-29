package config

type Config struct {
	Domain             string  `json:"domain" toml:"domain" yaml:"domain"`
	Static             string  `json:"static" toml:"static" yaml:"static"`
	MaxMultipartMemory int64   `json:"max_multipart_memory" toml:"max_multipart_memory" yaml:"max_multipart_memory"`
	Mysql              Mysql   `json:"mysql" toml:"mysql" yaml:"mysql"`
	Captcha            Captcha `json:"captcha" toml:"captcha" yaml:"captcha"`
	Cache              Cache   `json:"cache" toml:"cache" yaml:"cche"`
	Jwt                Jwt     `json:"jwt" toml:"jwt" yaml:"jwt"`
	Pay                Pay     `json:"pay" toml:"pay" yaml:"pay"`
}

type Mysql struct {
	User      string `json:"user" toml:"user" yaml:"user"`
	Password  string `json:"password" toml:"password" yaml:"password"`
	Host      string `json:"host" toml:"host" yaml:"host"`
	Port      int    `json:"port" toml:"port" yaml:"port"`
	DBName    string `json:"db_name" toml:"db_name" yaml:"db_name"`
	Charset   string `json:"charset" toml:"charset" yaml:"charset"`
	ParseTime string `json:"parseTime" toml:"parseTime" yaml:"parseTime"`
	Loc       string `json:"loc" toml:"loc" yaml:"loc"`
}

type Captcha struct {
	Expir      int64 `json:"expir" toml:"expir" yaml:"expir"` // 过期时间
	Width      int   `json:"width" toml:"width" yaml:"width"`
	Height     int   `json:"height" toml:"height" yaml:"height"`
	Length     int   `json:"length" toml:"length" yaml:"length"`
	NoiseCount int   `json:"noiseCount" toml:"noiseCount" yaml:"noiseCount"` // 噪点数量
}
type Cache struct {
	Use   string     `json:"use" toml:"use" yaml:"use"`
	Redis CacheRedis `json:"redis" toml:"redis" yaml:"redis"`
	File  CacheFile  `json:"file" toml:"file" yaml:"file"`
}
type CacheRedis struct {
	Addr     string `json:"addr" toml:"addr" yaml:"addr"`
	Port     int    `json:"port" toml:"port" yaml:"port"`
	Password string `json:"password" toml:"password" yaml:"password"`
}
type CacheFile struct {
	Path string `json:"path" toml:"path" yaml:"path"`
}

type Jwt struct {
	Expir int64 `json:"expir" toml:"expir" yaml:"expir"` // jwt登录过期时间，分钟，1440一天
}

type Pay struct {
	Wechat PayWechat `json:"wechat" toml:"wechat" yaml:"wechat"`
	Ali    PayAli    `json:"ali" toml:"ali" yaml:"ali"`
}
type PayWechat struct {
	AppID      string `json:"app_id" toml:"app_id" yaml:"app_id"`
	MchID      string `json:"mch_id" toml:"mch_id" yaml:"mch_id"`
	SerialNo   string `json:"serial_no" toml:"serial_no" yaml:"serial_no"`
	ApiV3Key   string `json:"api_v3_key" toml:"api_v3_key" yaml:"api_v3_key"`
	PrivateKey string `json:"private_key" toml:"private_key" yaml:"private_key"` // 文件路径
	NotifyUrl  string `json:"notify_url" toml:"notify_url" yaml:"notify_url"`
}

type PayAli struct {
	AppID               string `json:"app_id" toml:"app_id" yaml:"app_id"`
	PrivateKey          string `json:"private_key" toml:"private_key" yaml:"private_key"`                                  // 文件路径
	AppPublicCert       string `json:"app_public_cert" toml:"app_public_cert" yaml:"app_public_cert"`                      // 文件路径
	AlipayRootCert      string `json:"alipay_root_cert" toml:"alipay_root_cert" yaml:"alipay_root_cert"`                   // 文件路径
	AlipayCertPublicKey string `json:"alipay_cert_public_key" toml:"alipay_cert_public_key" yaml:"alipay_cert_public_key"` // 文件路径
}

func DefaultConfig() Config {
	return Config{
		Domain:             "",
		Static:             "data/static",
		MaxMultipartMemory: 100,
		Mysql: Mysql{
			User:      "root",
			Password:  "root",
			Host:      "127.0.0.1",
			Port:      3306,
			DBName:    "",
			Charset:   "utf8mb4",
			ParseTime: "True",
			Loc:       "Local",
		},
		Captcha: Captcha{
			Expir:      5,
			Width:      100,
			Height:     40,
			Length:     4,
			NoiseCount: 1,
		},
		Cache: Cache{
			Use: "file",
			Redis: CacheRedis{
				Addr:     "127.0.0.1",
				Port:     6379,
				Password: "",
			},
			File: CacheFile{
				Path: "data/cache",
			},
		},
		Jwt: Jwt{
			Expir: 10080,
		},
		Pay: Pay{
			Wechat: PayWechat{
				AppID:      "",
				MchID:      "",
				SerialNo:   "",
				ApiV3Key:   "",
				PrivateKey: "",
				NotifyUrl:  "",
			},
			Ali: PayAli{
				AppID:               "",
				PrivateKey:          "",
				AppPublicCert:       "",
				AlipayRootCert:      "",
				AlipayCertPublicKey: "",
			},
		},
	}
}
