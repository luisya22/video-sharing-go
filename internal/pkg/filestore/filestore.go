package filestore

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"luismatosgarcia.dev/video-sharing-go/internal/tests"
	"mime/multipart"
)

type Config struct {
	AwsAccessKeyId string
	AwsSecretKey   string
	AwsBucketName  string
	AwsRegion      string
	AwsEndpoint    string
}

type FileStore interface {
	Set(id int64, file *io.Reader, fileHeader *multipart.FileHeader) (string, error)
	Get()
}

type S3Bucket struct {
	region         string
	bucketName     string
	endpoint       string
	awsAccessKeyId string
	awsSecretKey   string
}

func (s S3Bucket) Set(id int64, file *io.Reader, fileHeader *multipart.FileHeader) (string, error) {
	//Config: Region, Credentials, Config.EndpointResolverWithOptions

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           s.endpoint,
			SigningRegion: s.region,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s.awsAccessKeyId, s.awsSecretKey, "")),
	)
	if err != nil {
		return "", err
	}

	s3client := s3.NewFromConfig(cfg)
	_, err = s.UploadFile(s3client, file, fileHeader)
	if err != nil {
		return "", err
	}

	return "videos/" + fileHeader.Filename, nil
}

func (s S3Bucket) UploadFile(client *s3.Client, file *io.Reader, fileHeader *multipart.FileHeader) (*s3.PutObjectOutput, error) {
	uploadOutput, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        &s.bucketName,
		Key:           aws.String("videos/" + fileHeader.Filename),
		Body:          *file,
		ContentLength: fileHeader.Size,
	})
	if err != nil {
		return nil, err
	}

	return uploadOutput, nil
}

func (s S3Bucket) Get() {
	//TODO implement me
	panic("implement me")
}

func NewFileStore(fileStoreType string, cfg Config) (FileStore, error) {
	//TODO Change to registry
	switch fileStoreType {
	case "s3":
		return S3Bucket{
			region:         cfg.AwsRegion,
			bucketName:     cfg.AwsBucketName,
			awsAccessKeyId: cfg.AwsAccessKeyId,
			awsSecretKey:   cfg.AwsSecretKey,
			endpoint:       cfg.AwsEndpoint,
		}, nil
	default:
		return S3Bucket{
			region:         cfg.AwsRegion,
			bucketName:     cfg.AwsBucketName,
			awsAccessKeyId: cfg.AwsAccessKeyId,
			awsSecretKey:   cfg.AwsSecretKey,
			endpoint:       cfg.AwsEndpoint,
		}, nil
	}
}

// Mocks

type Mock struct {
	FnCalls map[string]int
	Str     string
	Err     error
}

func (f Mock) Set(id int64, file *io.Reader, fileHeader *multipart.FileHeader) (string, error) {
	tests.Called(f.FnCalls, "Set")
	return f.Str, f.Err
}

func (f Mock) Get() {
	tests.Called(f.FnCalls, "Get")
}

func (f Mock) GetFnCalls(fnName string) int {
	value, exists := f.FnCalls[fnName]

	if !exists {
		return 0
	}

	return value
}
