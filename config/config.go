package config

type RDB struct {
	Dbfilename string
	Dir        string
}

type Config struct {
	RDB RDB
}

func NewConfig(dbfilename, dir string) *Config {
	return &Config{
		RDB: RDB{
			Dbfilename: dbfilename,
			Dir:        dir,
		},
	}
}

func (c *Config) GetConfig() map[string]string {
	return map[string]string{
		"dbfilename": c.RDB.Dbfilename,
		"dir": c.RDB.Dir,
	}
}
