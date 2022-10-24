package filestore

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"luismatosgarcia.dev/video-sharing-go/internal/server/http"
	"mime/multipart"
)

type Filestore interface {
	Set(file *multipart.File, fileHeader *multipart.FileHeader) (string, error)
	Get()
}

type S3Bucket struct {
	region     string
	bucketName string
	endpoint   string
	creds      *credentials.Credentials
}

func (s S3Bucket) Set(file *multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region:      aws.String(s.region),
			Credentials: s.creds,
			Endpoint:    &s.endpoint,
		},
	})
	if err != nil {
		return "", err
	}

	uploader := s3manager.NewUploader(sess)

	uploadOutput, err := s.UploadFile(uploader, file, fileHeader)

	return uploadOutput.Location, err
}

func (s S3Bucket) UploadFile(uploader *s3manager.Uploader, file *multipart.File, fileHeader *multipart.FileHeader) (*s3manager.UploadOutput, error) {
	uploadOutput, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileHeader.Filename),
		Body:   *file,
	})
	if err != nil {
		return uploadOutput, err
	}

	return uploadOutput, err
}

func (s S3Bucket) Get() {
	//TODO implement me
	panic("implement me")
}

func NewFileStore(fileStoreType string, cfg *http.Config) (Filestore, error) {
	//TODO Change to registry
	switch fileStoreType {
	case "s3":
		return S3Bucket{
			region:     cfg.FileStore.AwsRegion,
			bucketName: cfg.FileStore.AwsBucketName,
			creds:      credentials.NewStaticCredentials(cfg.FileStore.AwsAccessKeyId, cfg.FileStore.AwsSecretKey, ""),
			endpoint:   cfg.FileStore.AwsEndpoint,
		}, nil
	default:
		return nil, nil
	}
}
