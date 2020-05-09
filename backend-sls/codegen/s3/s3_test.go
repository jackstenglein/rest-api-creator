package s3

import (
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

// --------------- Download Tests ----------------------

type downloadFunc func(io.WriterAt, *s3.GetObjectInput, ...func(*s3manager.Downloader)) (int64, error)

func (f downloadFunc) Download(w io.WriterAt, i *s3.GetObjectInput, o ...func(*s3manager.Downloader)) (int64, error) {
	return f(w, i, o...)
}

func downloadMock(mockInput *s3.GetObjectInput, mockErr error) downloadFunc {
	return func(_ io.WriterAt, input *s3.GetObjectInput, _ ...func(*s3manager.Downloader)) (int64, error) {
		if !reflect.DeepEqual(input, mockInput) {
			return 0, errors.NewServer("Incorrect mock input")
		}
		return 0, mockErr
	}
}

var downloadTests = []struct {
	name string

	// Input
	localPath string
	s3Key     string

	// Mock data
	mockInput *s3.GetObjectInput
	mockErr   error

	// Expected output
	wantErr error
}{
	{
		name:    "InvalidPath",
		wantErr: errors.Wrap(errors.NewServer("open : no such file or directory"), "Failed to create file ``"),
	},
	{
		name:      "DownloadError",
		localPath: "tempTestFile",
		s3Key:     "testKey",
		mockInput: &s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("BUCKET_NAME")),
			Key:    aws.String("testKey"),
		},
		mockErr: errors.NewServer("Download failure"),
		wantErr: errors.Wrap(errors.NewServer("Download failure"), "Failed to download file from S3"),
	},
	{
		name:      "SuccessfulInvocation",
		localPath: "tempTestFile",
		s3Key:     "testKey",
		mockInput: &s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("BUCKET_NAME")),
			Key:    aws.String("testKey"),
		},
	},
}

func TestDownload(t *testing.T) {
	for _, test := range downloadTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			downloadSvc = downloadMock(test.mockInput, test.mockErr)
			defer func() {
				downloadSvc = defaultDownloader
			}()

			// Execute
			err := Download(test.localPath, test.s3Key)

			// Verify
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got err %v; want %v", err, test.wantErr)
			}
		})
	}
}

// --------------- Upload Tests ----------------------

type uploadFunc func(*s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)

func (f uploadFunc) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	return f(input, options...)
}

func uploadMock(mockInput *s3manager.UploadInput, mockErr error) uploadFunc {
	return func(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
		if !reflect.DeepEqual(input.Bucket, mockInput.Bucket) || !reflect.DeepEqual(input.Key, mockInput.Key) {
			return nil, errors.NewServer("Incorrect mock input")
		}
		return nil, mockErr
	}
}

var uploadTests = []struct {
	name string

	// Input
	localPath string
	s3Key     string

	// Mock data
	mockInput *s3manager.UploadInput
	mockErr   error

	// Expected output
	wantErr error
}{
	{
		name:    "InvalidLocalPath",
		wantErr: errors.Wrap(errors.NewServer("open : no such file or directory"), "Failed to open file ``"),
	},
	{
		name:      "UploadError",
		localPath: "/tmp",
		s3Key:     "testKey",
		mockInput: &s3manager.UploadInput{
			Bucket: aws.String(os.Getenv("BUCKET_NAME")),
			Key:    aws.String("testKey"),
		},
		mockErr: errors.NewServer("Failed to upload"),
		wantErr: errors.Wrap(errors.NewServer("Failed to upload"), "Failed to upload file to S3"),
	},
	{
		name:      "SuccessfulInvocation",
		localPath: "/tmp",
		s3Key:     "testKey",
		mockInput: &s3manager.UploadInput{
			Bucket: aws.String(os.Getenv("BUCKET_NAME")),
			Key:    aws.String("testKey"),
		},
	},
}

func TestUpload(t *testing.T) {
	for _, test := range uploadTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			uploadSvc = uploadMock(test.mockInput, test.mockErr)
			defer func() {
				uploadSvc = defaultUploader
			}()

			// Execute
			err := Upload(test.localPath, test.s3Key)

			// Verify
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got err %v; want %v", err, test.wantErr)
			}
		})
	}
}
