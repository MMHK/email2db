package pkg

import (
	"email2db/tests"
	"testing"
)

func loadConfig() (*Config, error) {
	conf, err := NewConfigFromLocal(tests.GetLocalPath("../config.json"))
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func TestNewConfigFromLocal(t *testing.T) {
	conf, err := NewConfigFromLocal(tests.GetLocalPath("../config.json"))
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Logf("%+v", conf)
	t.Log("PASS")
}

func Test_SaveConfig(t *testing.T) {

	conf := &Config{
		Storage: &StorageConfig{
			S3: &S3Config{},
		},
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