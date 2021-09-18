package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Uploads a file to a specific bucket in S3 with the file name
// as the object's key. After it's uploaded, a message is sent
// to a queue.
func main() {
	mySession := session.Must(session.NewSession())

	svc := s3.New(mySession, aws.NewConfig().WithRegion("us-east-1"))

	multi, err := svc.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket: aws.String("obok-test"),
		Key:    aws.String("my-file"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fi, err := os.ReadFile("../client/gopher.png")
	if err != nil {
		log.Fatal(err)
	}

	// st := strings.NewReader(string(fi))

	st := bytes.NewReader(fi)

	res, err := svc.UploadPart(&s3.UploadPartInput{
		Body:       aws.ReadSeekCloser(st),
		UploadId:   multi.UploadId,
		Bucket:     aws.String("obok-test"),
		Key:        aws.String("my-file"),
		PartNumber: aws.Int64(1),
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.ETag)

	res2, err := svc.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		UploadId: multi.UploadId,
		Bucket:   aws.String("obok-test"),
		Key:      aws.String("my-file"),
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: []*s3.CompletedPart{
				{ETag: res.ETag, PartNumber: aws.Int64(1)},
			},
		},
	})

	fmt.Println(res2.ETag)

	if err != nil {
		log.Fatal(err)
	}

}
