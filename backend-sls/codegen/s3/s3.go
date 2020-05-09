package s3

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type downloader interface {
	Download(io.WriterAt, *s3.GetObjectInput, ...func(*s3manager.Downloader)) (int64, error)
}

type uploader interface {
	Upload(*s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

// The session the S3 Downloader will use
var sess = session.New()

var defaultDownloader = s3manager.NewDownloader(sess)
var defaultUploader = s3manager.NewUploader(sess)
var downloadSvc downloader = defaultDownloader
var uploadSvc uploader = defaultUploader

// Download retrieves the file with the given key from AWS S3 and saves it in the local filesystem
// at the specified local path.
func Download(localPath string, s3Key string) error {

	// Create the file to write the S3 Object contents to.
	file, err := os.Create(localPath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create file `%s`", localPath))
	}
	defer file.Close()

	// Write the contents of S3 Object to the file
	_, err = downloadSvc.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(s3Key),
	})
	return errors.Wrap(err, "Failed to download file from S3")
}

// Presign returns a presigned URL for the specified AWS S3 key. The key expires in 10 minutes.
func Presign(s3Key string) (string, error) {
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(s3Key),
	})
	url, err := req.Presign(10 * time.Minute)
	return url, errors.Wrap(err, "Failed to presign S3 request")
}

// Upload uploads a local file to AWS S3. The file is stored with the given key in the bucket
// specified by the environment variable BUCKET_NAME.
func Upload(localPath string, s3Key string) error {

	file, err := os.Open(localPath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to open file `%s`", localPath))
	}
	defer file.Close()

	_, err = uploadSvc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(s3Key),
		Body:   file,
	})
	return errors.Wrap(err, "Failed to upload file to S3")
}
