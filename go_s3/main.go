package main

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession())

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	//f, err := os.Open(filename)
	//if err != nil {
	//	return fmt.Errorf("failed to open file %q, %v", filename, err)
	//}

	var myBucket = "tsuboya1"
	var objectKey = "bbb/ccc.yaml"

	var buff bytes.Buffer
	buff.WriteString("This is a test.\nThe second line.")

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(myBucket),
		Key:    aws.String(objectKey),
		Body:   &buff,
	})
	if err != nil {
		fmt.Printf("Failed to upload file, %v", err)
		return
	}
	fmt.Printf("File uploaded to, %s\n", result.Location)
}
