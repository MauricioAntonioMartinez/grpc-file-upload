package upload

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AwsStorage struct {
	conf       AwsConfig
	uploadId   *string
	partsCount int
	svc        *s3.S3
	parts      []*s3.CompletedPart
}

type AwsConfig struct {
	Bucket string
	Key    string
}

type awsPart struct {
	data     bytes.Buffer
	uploadID *string
	partID   int
}

func NewAwsStorage(conf AwsConfig) *AwsStorage {
	return &AwsStorage{
		conf:  conf,
		parts: []*s3.CompletedPart{},
	}
}

func (u *AwsStorage) init() error {
	mySession := session.Must(session.NewSession())
	svc := s3.New(mySession, aws.NewConfig().WithRegion("us-east-1"))
	u.svc = svc
	multi, err := svc.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket: aws.String(u.conf.Bucket),
		Key:    aws.String(u.conf.Key),
	})
	if err != nil {
		return err
	}
	u.uploadId = multi.UploadId
	return nil
}

func (u *AwsStorage) newPart(data bytes.Buffer) part {
	u.partsCount++
	return awsPart{
		data:     data,
		uploadID: u.uploadId,
		partID:   u.partsCount,
	}
}

func (p awsPart) getData() bytes.Buffer {
	return p.data
}

func (p awsPart) getID() int64 {
	return int64(p.partID)
}

func (u *AwsStorage) sendPart(p part) error {
	data := p.getData()
	res, err := u.svc.UploadPart(&s3.UploadPartInput{
		Body:       aws.ReadSeekCloser(&data),
		UploadId:   u.uploadId,
		Bucket:     aws.String(u.conf.Bucket),
		Key:        aws.String(u.conf.Key),
		PartNumber: aws.Int64(p.getID()),
	})

	if err != nil {
		return err
	}

	fmt.Println(res.ETag)

	return nil
}

func (u *AwsStorage) closeUpload() error {
	_, err := u.svc.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		UploadId: u.uploadId,
		Bucket:   aws.String(u.conf.Bucket),
		Key:      aws.String(u.conf.Key),
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: u.parts,
		},
	})

	return err
}
