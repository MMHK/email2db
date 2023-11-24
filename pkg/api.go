package pkg

import (
	"net/mail"
	"strings"
	"time"
)

type APIEmailListItem struct {
	ID        int               `json:"id"`
	Subject   string            `json:"subject" gorm:"column:subject"`
	From      string            `json:"from" gorm:"column:from"`
	To        ToList            `json:"to" gorm:"column:to;type:json"`
	meta      map[string]string `gorm:"column:meta"`
	CreatedAt string            `json:"created_at" gorm:"column:created_at"`
}

func (this *APIEmailListItem) Date() time.Time {
	headers, ok := this.meta["headers"]
	if ok {
		message, err := mail.ReadMessage(strings.NewReader(headers))
		if err != nil {
			Log.Error(err)
		} else {
			dateStr := message.Header.Get("Date")
			date, err := time.Parse(time.RFC1123Z, dateStr)
			if err != nil {
				Log.Error(err)
			} else {
				return date
			}
		}
	}
	return time.Now()
}

type APIEmailDetail struct {
	APIEmailListItem
	Attachments []AttachmentModel `json:"attachments" gorm:"foreignKey:EmailID;references:ID"`
}

func (this *APIEmailDetail) GetAttachments() *[]AttachmentModel {
	if this.Attachments == nil {
		this.Attachments = make([]AttachmentModel, 0)
	}

	return &this.Attachments
}

func (this *APIEmailDetail) HTML() string {
	html, ok := this.meta["html"]
	if ok {
		return html
	}
	return ""
}

type APIListWrapper struct {
	Items      []APIEmailListItem `json:"items"`
	Pagination Pagination         `json:"pagination"`
}

type APIListResponse struct {
	APIStandardError
	Data APIListWrapper `json:"data"`
}

type APIEmailDetailResponse struct {
	APIStandardError
	Data APIEmailDetail `json:"data"`
}
