package database

type MysqlConfig struct {
	User      string `json:"user" toml:"user" yaml:"user"`
	Password  string `json:"password" toml:"password" yaml:"password"`
	Host      string `json:"host" toml:"host" yaml:"host"`
	Port      int    `json:"port" toml:"port" yaml:"port"`
	DBName    string `json:"db_name" toml:"db_name" yaml:"db_name"`
	Charset   string `json:"charset" toml:"charset" yaml:"charset"`
	ParseTime string `json:"parse_time" toml:"parse_time" yaml:"parse_time"`
	Loc       string `json:"loc" toml:"loc" yaml:"loc"`
}
