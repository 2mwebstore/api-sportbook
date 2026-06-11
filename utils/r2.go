package utils

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var r2Client *s3.Client

// InitR2 initialises the Cloudflare R2 S3-compatible client.
// Call once from main() after loading env vars.
func InitR2() {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKey := os.Getenv("R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")

	r2Endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
		config.WithRegion("auto"),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to init R2 config: %v", err))
	}

	r2Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(r2Endpoint)
		o.UsePathStyle = true
	})
}

// UploadFile uploads a multipart file header to R2 and returns the public URL.
func UploadFile(file *multipart.FileHeader, folder string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("cannot open uploaded file: %w", err)
	}
	defer src.Close()

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(src); err != nil {
		return "", fmt.Errorf("cannot read uploaded file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	key := fmt.Sprintf("%s/%d%s", folder, time.Now().UnixNano(), ext)

	bucket := os.Getenv("R2_BUCKET_NAME")
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = r2Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("R2 upload failed: %w", err)
	}

	publicBase := strings.TrimRight(os.Getenv("R2_PUBLIC_URL"), "/")
	return fmt.Sprintf("%s/%s", publicBase, key), nil
}

// DeleteFile removes an object from R2 by its full public URL.
func DeleteFile(publicURL string) error {
	publicBase := strings.TrimRight(os.Getenv("R2_PUBLIC_URL"), "/")
	key := strings.TrimPrefix(publicURL, publicBase+"/")
	bucket := os.Getenv("R2_BUCKET_NAME")

	_, err := r2Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}

// UploadMultipleFiles uploads a slice of file headers and returns their public URLs.
func UploadMultipleFiles(files []*multipart.FileHeader, folder string) ([]string, error) {
	urls := make([]string, 0, len(files))
	for _, f := range files {
		url, err := UploadFile(f, folder)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}
