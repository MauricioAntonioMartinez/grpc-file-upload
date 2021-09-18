package upload

import "errors"

var (
	ErrorDeadlineExceeded  = errors.New("dead line exceeded")
	ErrorRequestCanceled   = errors.New("request canceled")
	ErrorImageSizeTooLarge = errors.New("image size too large")
)
