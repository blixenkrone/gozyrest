package storage

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	s3Region = "eu-north-1"
	s3Bucket = "image-exif-storage"
)

// File places the file in right name and dir in S3
type File struct {
	Body   []byte
	Name   string
	Format string
	Dir    string
}

// InitFile creates a new file for upload
func InitFile(body []byte, fileName string, format, string, dir string) *File {

	return &File{
		Body:   body,
		Name:   fileName,
		Format: format,
		Dir:    dir,
	}
}

// UploadFile S3 uploader
func (file *File) UploadFile() error {
	newSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(s3Region),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS"), os.Getenv("AWS_SECRET"), ""),
	})
	if err != nil {
		return err
	}
	session := session.Must(newSession, err)
	uploader := s3manager.NewUploader(session)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:                 bytes.NewBuffer(file.Body),
		Bucket:               aws.String(s3Bucket),
		Key:                  aws.String(file.Dir + "/" + string(file.Name)),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return fmt.Errorf("Failed to upload file:  %v", err)
	}
	fmt.Printf("Successfully uploaded file to: %s\n", aws.StringValue(&result.Location))
	return nil
}

// GetAWSFile to fetch file from s3
// func GetAWSFile(fileName string) []byte {
// 	buf := &aws.WriteAtBuffer{}
// 	sess, _ := session.NewSession(&aws.Config{
// 		Region:      aws.String(s3Region),
// 		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS"), os.Getenv("AWS_SECRET"), ""),
// 	})
// 	dl := s3manager.NewDownloader(sess)
// 	_, err := dl.Download(buf, &s3.GetObjectInput{
// 		Bucket: aws.String(s3SecretBucket),
// 		Key:    aws.String(fileName),
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return buf.Bytes()
// }
