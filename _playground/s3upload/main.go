package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	Bucket   *string
	Client   *s3.Client
	Uploader *manager.Uploader
}

func getClient(endpoint string) *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		slog.Error("Could not connect to s3 server", "Reason", err.Error())
		panic(err)
	}
	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})
}

func (s *S3) Upload(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		slog.Error("Could not open file", "filename", filename, "Reason", err.Error())
		panic(err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}()

	req := &s3.PutObjectInput{
		Bucket: aws.String("remote-localstack"),
		Key:    aws.String(filename),
		Body:   file,
	}
	resp, err := s.Uploader.Upload(context.TODO(), req)
	if err != nil {
		slog.Error("Could not upload file to s3", "filename", filename, "Reason", err.Error())
	} else {
		slog.Info("File uploaded", "filename", filename, "location", resp.Location)
	}
}

func (s *S3) GetAll() {
	resp, err := s.Client.ListObjectsV2(
		context.TODO(),
		&s3.ListObjectsV2Input{Bucket: aws.String("remote-localstack")},
	)
	if err != nil {
		slog.Error("Could not get all files from s3")
	}
	slog.Info("Listing all files: =======================")
	for _, obj := range resp.Contents {
		slog.Info("File", "filename", aws.ToString(obj.Key), "size", *obj.Size)
	}
}

func main() {
	client := getClient(os.Getenv("AWS_ENDPOINT"))

	s3 := &S3{
		Client:   client,
		Uploader: manager.NewUploader(client),
		Bucket:   aws.String("remote-localstack"),
	}
	s3.Upload("main.go")
	s3.GetAll()
}

// AWS_ENDPOINT=http://10.40.160.124:4566 AWS_REGION=us-east-1 AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test AWS_REGION=us-east-1 go run main.go
