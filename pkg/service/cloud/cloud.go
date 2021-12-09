// Package cloud uses for uploading video to cloud
package cloud

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

const presignTime = time.Minute

// AWSS3 represent aws s3 storage
type AWSS3 struct {
	bucketName string
	accessKey  string
	secretKey  string
	region     string
}

// NewS3Storage initialize *AWS3
func NewS3Storage(bucketName, region, accessKey, secretKey string) *AWSS3 {
	return &AWSS3{
		bucketName: bucketName,
		accessKey:  accessKey,
		secretKey:  secretKey,
		region:     region,
	}
}

// Upload file to aws s3
func (cloud *AWSS3) Upload(ctx context.Context, header *multipart.FileHeader) (string, error) {
	file, err := header.Open()
	if err != nil {
		return "", err
	}

	sess, err := cloud.session()

	if err != nil {
		return "", err
	}

	uploader := s3manager.NewUploader(sess)

	newFileName := fmt.Sprintf("%s_%s", uuid.New().String(), header.Filename)
	_, err = uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(cloud.bucketName),
		Key:    aws.String(newFileName),
	})

	if err != nil {
		return "", err
	}

	return newFileName, nil
}

func (cloud *AWSS3) URL(filename string) (string, error) {
	sess, err := cloud.session()

	if err != nil {
		return "", err
	}

	s3svc := s3.New(sess)
	req, _ := s3svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(cloud.bucketName),
		Key:    aws.String(filename),
	})

	return req.Presign(presignTime)
}

func (cloud *AWSS3) session() (*session.Session, error) {
	return session.NewSession(
		&aws.Config{
			Region: aws.String(cloud.region),
			Credentials: credentials.NewStaticCredentials(
				cloud.accessKey,
				cloud.secretKey,
				""), // a token will be created when the session it's used.
		})
}
