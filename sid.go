package upload

import (
	"context"
	"fmt"
	"io"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Stream interface {
	grpc.ServerStream
}

type Uploader struct {
	conf    UploaderConfig
	msg     proto.Message
	parts   chan part
	wait    sync.WaitGroup
	close   chan interface{}
	Storage SourceStorage
}

type UploaderConfig struct {
	MessageType   proto.Message
	MessageNumber int
	MaxSize       int
	ChuckSize     int
	Verbose       bool
}

func NewUploader(conf UploaderConfig, storage SourceStorage) *Uploader {
	return &Uploader{
		conf:    conf,
		msg:     conf.MessageType,
		parts:   make(chan part),
		wait:    sync.WaitGroup{},
		close:   make(chan interface{}),
		Storage: storage,
	}
}

type UploaderResponse struct{}

func (u *Uploader) Upload(stream Stream) (*UploaderResponse, error) {
	partSize, fileSize, partData := 0, 0, []byte{}

	if err := initStorage(u.Storage); err != nil {
		return nil, err
	}

	go u.partsSender()

	for {
		err := contextErr(stream.Context())
		if err != nil {
			return nil, err
		}

		err = stream.RecvMsg(u.msg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		data := u.getData()
		partData = append(partData, data...)
		size := len(data)
		fileSize += size
		partSize += size

		if fileSize > u.conf.MaxSize {
			return nil, ErrorImageSizeTooLarge
		}

		if partSize >= u.conf.ChuckSize {
			u.wait.Add(1)
			u.parts <- u.Storage.newPart(partData)
			partData = []byte{}
			partSize = 0
		}
	}

	u.close <- true
	u.wait.Wait()
	if err := u.Storage.closeUpload(); err != nil {
		return nil, err
	}

	return &UploaderResponse{}, nil
}

func (u *Uploader) getData() []byte {
	ref := u.msg.ProtoReflect()
	field := ref.Descriptor().Fields().Get(u.conf.MessageNumber)
	data := ref.Get(field).Bytes()
	return data
}

func (u *Uploader) partsSender() {
lis:
	for {
		select {
		case part := <-u.parts:
			go u.sendPart(part)
		case <-u.close:
			fmt.Println("Done reading bytes, clean up ....")
			break lis
		}
	}
}

func (u *Uploader) sendPart(p part) {
	defer u.wait.Done()
	u.Storage.sendPart(p)
}

func contextErr(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return ErrorRequestCanceled
	case context.DeadlineExceeded:
		return ErrorDeadlineExceeded
	default:
		return nil
	}
}
