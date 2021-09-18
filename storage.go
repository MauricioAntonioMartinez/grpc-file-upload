package upload

import "bytes"

type SourceStorage interface {
	init() error
	sendPart(part) error
	newPart(data []byte) part
	closeUpload() error
}

type part interface {
	getData() *bytes.Reader
	getID() int64
}

func initStorage(storage SourceStorage) error {
	return storage.init()
}
