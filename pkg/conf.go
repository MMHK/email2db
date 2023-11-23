package pkg

import (
	"encoding/json"
	"os"
)

type Config struct {
	HttpConfig
	Storage *StorageConfig `json:"storages"`
	DB      *DBConfig      `json:"db"`
	TempDir string         `json:"tmp_path"`
}

func NewConfigFromLocal(filename string) (*Config, error) {
	conf := &Config{}
	err := conf.load(filename)
	if err == nil && len(conf.TempDir) <= 0 {
		conf.TempDir = os.TempDir()
	}
	return conf, err
}

func (this *Config) MarginWithENV()  {
	if this.Storage == nil {
		this.Storage = &StorageConfig{
			S3: LoadS3ConfigWithEnv(),
		}
	}
	if this.DB == nil {
		this.DB = &DBConfig{
			MySQL: &MySQLConfig{
				DSN: LoadMySQLDSNWithENV(),
			},
		}
	}
	if len(this.WebRoot) <= 0 {
		this.WebRoot = os.Getenv("WEB_ROOT")
	}
	if len(this.Listen) <= 0 {
		this.Listen = os.Getenv("HTTP_LIST")
	}
}

func (c *Config) load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		Log.Error(err)
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		Log.Error(err)
	}
	return err
}

func (c *Config) Save(saveAs string) error {
	file, err := os.Create(saveAs)
	if err != nil {
		Log.Error(err)
		return err
	}
	defer file.Close()
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		Log.Error(err)
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		Log.Error(err)
	}
	return err
}