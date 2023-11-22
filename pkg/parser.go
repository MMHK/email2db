package pkg

import (
	"email2db/pkg/parser"
	"io"
	"net/http"
)

type AttachmentRaw struct {
	FileName string
	File io.ReadCloser
	MimeType string
}

type IParser interface {
	Parse(request *http.Request) (IParser, error)
	GetSubject() string
	GetToList() *[]string
	GetMate() *map[string]string
	GetReplyTo() string
	GetFrom() string
	GetAttachments() *[]AttachmentRaw
}

func GetSendGridParser(request *http.Request) (IParser, error) {
	target := &parser.SendGridParser{}

	return target.Parse(request)
}
