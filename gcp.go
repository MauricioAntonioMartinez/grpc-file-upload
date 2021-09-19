package upload

import (
	"bytes"
	"os"
)

type GoogleCloudStorage struct {
	conf GoogleCloudStorageConfig
	buf  *bytes.Buffer
}

type GoogleCloudStorageConfig struct {
	Path string
}

func NewGoogleCloudStorage(conf GoogleCloudStorageConfig) *GoogleCloudStorage {
	return &GoogleCloudStorage{
		conf: conf,
	}
}

func (u *GoogleCloudStorage) init() error {
	newBuf := bytes.Buffer{}
	u.buf = &newBuf
	return nil
}

func (u *GoogleCloudStorage) newPart(content []byte) part {
	u.buf.Write(content)
	return nil
}

func (u *GoogleCloudStorage) sendPart(p part) error {
	return nil
}

func (u *GoogleCloudStorage) closeUpload() error {
	return os.WriteFile(u.conf.Path, u.buf.Bytes(), 0644)
}
