package pkg

import (
	"email2db/tests"
	"testing"
)

func loadConfig() (*Config, error) {
	conf, err := NewConfigFromLocal(tests.GetLocalPath("../conf.json"))
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func Test_SaveConfig(t *testing.T) {

	conf := &Config{
		Storage: &StorageConfig{},
		DB: &DBConfig{
			MySQL: &MySQLConfig{},
		},
	}
	configPath := tests.GetLocalPath("../config.json")
	err := conf.Save(configPath)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Log("PASS")
}