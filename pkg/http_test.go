package pkg

import (
	"os"
	_ "email2db/tests"
	"testing"
)

func loadHTTPConfig() (*Config) {
	conf := &Config{
		HttpConfig: HttpConfig{
			Listen: os.Getenv("HTTP_LIST"),
			WebRoot: os.Getenv("WEB_ROOT"),
		},
	}

	conf.MarginWithENV()

	return conf
}

func TestHTTPService_Start(t *testing.T) {
	service := NewHttpService(loadHTTPConfig())
	service.Start()

	t.Log("PASS")
}