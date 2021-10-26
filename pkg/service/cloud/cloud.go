// Package cloud uses for uploading video to cloud
package cloud

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"math/rand"
	"mime/multipart"
	"time"
)

const (
	randLetterLength = 8
	minLetterIndex   = 97  // 'a'
	maxLetterIndex   = 122 // 'z'
)

// AWSS3 represent aws s3 storage
type AWSS3 struct {
	bucketName string
	accessKey  string
	secretKey  string
	region     string
}

func init() {
	rand.Seed(time.Now().UnixNano())
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

	sess, err := cloud.createSession()
	if err != nil {
		return "", err
	}

	uploader := s3manager.NewUploader(sess)

	newFileName := cloud.generateName(header.Filename)
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

// generateName for file in aws
func (cloud *AWSS3) generateName(fileName string) string {
	var letters []byte
	buf := bytes.NewBuffer(letters)
	for i := 0; i < randLetterLength; i++ {
		buf.WriteByte(byte(minLetterIndex + (rand.Intn(maxLetterIndex - minLetterIndex))))
	}
	newName := fmt.Sprintf("original_%s_%s", buf.String(), fileName)
	return newName
}

// createSession for aws
func (cloud *AWSS3) createSession() (*session.Session, error) {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(cloud.region),
			Credentials: credentials.NewStaticCredentials(
				cloud.accessKey,
				cloud.secretKey,
				""), // a token will be created when the session it's used.
		})
	if err != nil {
		return nil, err
	}

	return sess, nil
}
