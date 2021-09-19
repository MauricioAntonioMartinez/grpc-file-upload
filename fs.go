package upload

import (
	"bytes"
	"os"
)

type FileSystemStorage struct {
	conf FileSystemStorageConfig
	buf  *bytes.Buffer
}

type FileSystemStorageConfig struct {
	Path string
}

func NewFileSystemStorage(conf FileSystemStorageConfig) *FileSystemStorage {
	return &FileSystemStorage{
		conf: conf,
	}
}

func (u *FileSystemStorage) init() error {
	newBuf := bytes.Buffer{}
	u.buf = &newBuf
	return nil
}

func (u *FileSystemStorage) newPart(content []byte) part {
	u.buf.Write(content)
	return nil
}

func (u *FileSystemStorage) sendPart(p part) error {
	return nil
}

func (u *FileSystemStorage) closeUpload() error {
	return os.WriteFile(u.conf.Path, u.buf.Bytes(), 0644)
}
