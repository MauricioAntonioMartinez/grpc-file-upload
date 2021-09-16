package upload

import "bytes"

type SourceStorage interface {
	init() error
	sendPart(part) error
	newPart(data bytes.Buffer) part
}

type part interface {
	getData() bytes.Buffer
	getID() int64
}

func initStorage(storage SourceStorage) error {
	return storage.init()
}
