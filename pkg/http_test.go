package pkg

import (
	"os"
	_ "email2db/tests"
	"testing"
)

func loadHTTPConfig() (*Config) {
	return &Config{
		HttpConfig: HttpConfig{
			Listen: os.Getenv("HTTP_LIST"),
			WebRoot: os.Getenv("WEB_ROOT"),
		},
		Storage: &StorageConfig{
			S3: loadS3Config(),
		},
		DB: &DBConfig{
			MySQL: &MySQLConfig{
				DSN: loadMySQLDSN(),
			},
		},
	}
}

func TestHTTPService_Start(t *testing.T) {
	service := NewHttpService(loadHTTPConfig())
	service.Start()

	t.Log("PASS")
}