package types

import "io"

type UploadedExample struct {
	Index       int
	S3Key       string
	S3VersionID *string
	MimeType    string
}

type ExpectedDocumentTemplateInput struct {
	Title           string
	Description     string
	ExampleFile     io.Reader
	ExampleFileName string
	ExampleMimeType string
	ExampleFileSize int64
}
