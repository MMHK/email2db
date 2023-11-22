package storage

import (
	"email2db/pkg"
	"email2db/tests"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"testing"
)

func loadS3Config() *pkg.S3Config {
	return &pkg.S3Config{
		AccessKey: os.Getenv("S3_KEY"),
		SecretKey: os.Getenv("S3_SECRET"),
		Bucket: os.Getenv("S3_BUCKET"),
		Region: os.Getenv("S3_REGION"),
		PrefixPath: "email2db",
	}
}

func getStorage(t *testing.T) pkg.IStorage {
	s3, err := NewS3Storage(loadS3Config())

	if err != nil {
		t.Error(err)
		return nil
	}

	return s3
}

func TestPutStream(t *testing.T) {
	disk := getStorage(t)

	filename := tests.GetLocalPath("../tests/sample.jpeg")

	file, err := os.Open(filename)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	defer file.Close()

	distPath := fmt.Sprintf("%s%s", pkg.MakeUUID(), filepath.Ext(filename))

	path, url, err := disk.PutStream(file, distPath, &pkg.UploadOptions{
		ContentType: mime.TypeByExtension(filename),
	})

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(path)
	t.Log(url)
}
