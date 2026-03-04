package types

type UploadedExample struct {
	Index       int
	S3Key       string
	S3VersionID *string
	MimeType    string
}
