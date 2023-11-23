package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

type S3Storage struct {
	Conf    *S3Config
	session *session.Session
}

func LoadS3ConfigWithEnv() *S3Config {
	remotePathPrefix := "email2db"
	prefix := os.Getenv("S3_PREFIX")
	if len(prefix) > 0 {
		remotePathPrefix = prefix
	}

	return &S3Config{
		AccessKey: os.Getenv("S3_KEY"),
		SecretKey: os.Getenv("S3_SECRET"),
		Bucket: os.Getenv("S3_BUCKET"),
		Region: os.Getenv("S3_REGION"),
		PrefixPath: remotePathPrefix,
	}
}

func NewS3Storage(conf *S3Config) (IStorage, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(conf.Region),
		Credentials: credentials.NewStaticCredentials(conf.AccessKey, conf.SecretKey, ""),
	})
	if err != nil {
		Log.Error(err)
		return nil, err
	}

	return &S3Storage{
		Conf:    conf,
		session: sess,
	}, nil
}

func (this *S3Storage) Upload(localPath string, Key string, opt *UploadOptions) (path string, url string, err error) {
	file, err := os.Open(localPath)
	if err != nil {
		return "", "", err
	}

	defer file.Close()

	uploader := s3manager.NewUploader(this.session)
	path = filepath.ToSlash(filepath.Join(this.Conf.PrefixPath, Key))

	var publicflag *string
	if opt.PublicRead {
		publicflag = aws.String("public-read")
	}

	info, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(this.Conf.Bucket),
		Key:         aws.String(path),
		Body:        file,
		ACL:         publicflag,
		ContentType: aws.String(mime.TypeByExtension(localPath)),
	})

	return path, info.Location, err
}

func (this *S3Storage) PutContent(content string, Key string, opt *UploadOptions) (path string, url string, err error) {
	return this.PutStream(strings.NewReader(content), Key, opt)
}

func (this *S3Storage) PutStream(reader io.Reader, Key string, opt *UploadOptions) (path string, url string, err error) {
	uploader := s3manager.NewUploader(this.session)

	contentType := "application/octet-stream"
	if len(opt.ContentType) > 0 {
		contentType = opt.ContentType
	}

	path = filepath.ToSlash(filepath.Join(this.Conf.PrefixPath, Key))

	var publicflag *string
	if opt.PublicRead {
		publicflag = aws.String("public-read")
	}

	info, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(this.Conf.Bucket),
		Key:         aws.String(path),
		Body:        reader,
		ACL:         publicflag,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		Log.Error(err)
		return path, "", err
	}

	return path, info.Location, err
}
