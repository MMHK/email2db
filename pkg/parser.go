package pkg

import (
	"io"
	"net/http"
)

type AttachmentRaw struct {
	FileName  string `json:"filename"`
	File      io.ReadCloser
	MimeType  string `json:"mime_type"`
	ContentID string `json:"content_id"`
}

type IParser interface {
	Parse(request *http.Request) (IParser, error)
	GetSubject() string
	GetToList() []string
	GetMate() map[string]string
	GetFrom() string
	GetAttachments() []*AttachmentRaw
}

func GetSendGridParser(request *http.Request) (IParser, error) {
	target := &SendGridParser{}

	return target.Parse(request)
}
