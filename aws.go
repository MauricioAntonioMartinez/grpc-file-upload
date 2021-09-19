package upload

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// AwsStorage is an implementation of Storage interface
type AwsStorage struct {
	conf       AwsConfig
	uploadId   *string
	partsCount int
	svc        s3iface.S3API
	parts      []*s3.CompletedPart
}

// AwsConfig everyting needed to upload to AWS
type AwsConfig struct {
	Bucket string        //  Bucket to send the file to.
	Key    string        //  key of the file to send.
	Config []*aws.Config // you might skip this filed if you have env variables set or aws roles are already setted for you.
}

type awsPart struct {
	data     *bytes.Reader
	uploadID *string
	partID   int
}

// NewAwsStorage is a constructor to retreive an AwsStorage struct.
func NewAwsStorage(conf AwsConfig) *AwsStorage {
	return &AwsStorage{
		conf:  conf,
		parts: []*s3.CompletedPart{},
	}
}

func (u *AwsStorage) init() error {
	mySession := session.Must(session.NewSession(u.conf.Config...))
	svc := s3.New(mySession)

	u.svc = svc
	multi, err := svc.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket: aws.String(u.conf.Bucket),
		Key:    aws.String(u.conf.Key),
	})
	if err != nil {
		fmt.Println("Creating multi ", err)
		return err
	}
	u.uploadId = multi.UploadId
	return nil
}

func (u *AwsStorage) newPart(content []byte) part {
	fmt.Println("Content of length ", len(content))
	data := bytes.NewReader(content)
	u.partsCount++
	p := awsPart{
		uploadID: u.uploadId,
		partID:   u.partsCount,
		data:     data,
	}
	return p
}

func (p awsPart) getData() *bytes.Reader {
	return p.data
}

func (p awsPart) getID() int64 {
	return int64(p.partID)
}

func (u *AwsStorage) sendPart(p part) error {
	data := p.getData()

	res, err := u.svc.UploadPart(&s3.UploadPartInput{
		Body:       aws.ReadSeekCloser(data),
		UploadId:   u.uploadId,
		Bucket:     aws.String(u.conf.Bucket),
		Key:        aws.String(u.conf.Key),
		PartNumber: aws.Int64(p.getID()),
	})

	if err != nil {
		fmt.Println("error sending part -> ", err)
		return err
	}
	fmt.Printf("Sent part #%d length of %d \n", p.getID(), data.Len())

	u.parts = append(u.parts, &s3.CompletedPart{
		ETag:       res.ETag,
		PartNumber: aws.Int64(p.getID()),
	})

	return nil
}

func (u *AwsStorage) closeUpload() error {

	sort.Slice(u.parts, func(i, j int) bool {
		return *u.parts[i].PartNumber < *u.parts[j].PartNumber
	})
	fmt.Println("Closing multipart .....")
	_, err := u.svc.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		UploadId: u.uploadId,
		Bucket:   aws.String(u.conf.Bucket),
		Key:      aws.String(u.conf.Key),
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: u.parts,
		},
	})

	if err != nil {
		fmt.Println("closing multi ", err)
	}
	fmt.Println("Multipart closed file uploaded.")
	return err
}
