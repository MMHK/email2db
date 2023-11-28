package pkg

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

type MateData map[string]string

func (this *MateData) Scan(value interface{}) error {
	bin, ok := value.([]byte)
	if !ok {
		this = &MateData{}
		return nil
	}
	return json.Unmarshal(bin, this)
}

type APIEmailListItem struct {
	ID        int      `json:"id"`
	Subject   string   `json:"subject" gorm:"column:subject"`
	From      string   `json:"from" gorm:"column:from"`
	To        ToList   `json:"to" gorm:"column:to;type:json"`
	Meta      MateData `gorm:"column:meta;type:json" json:"-"`
	Date      string   `json:"date" gorm:"-"`
	CreatedAt string   `json:"created_at" gorm:"column:created_at"`
}

func (this *APIEmailListItem) FillDate() {
	headers, ok := this.Meta["headers"]
	layout := "2006-01-02T15:04:05-07:00"
	if ok {
		message, err := mail.ReadMessage(strings.NewReader(fmt.Sprintf("%s\n", headers)))
		if err != nil {
			Log.Error(err)
		} else {
			dateStr := message.Header.Get("Date")
			date, err := time.Parse(time.RFC1123Z, dateStr)
			if err != nil {
				Log.Error(err)
			} else {
				this.Date = date.Format(layout)
				return
			}
		}
	}
	this.Date = time.Now().Format(layout)
}

type APIEmailDetail struct {
	APIEmailListItem
	HTML        string            `json:"html"`
	Attachments []AttachmentModel `json:"attachments" gorm:"foreignKey:EmailID;references:ID"`
}

func (this *APIEmailDetail) GetAttachments() *[]AttachmentModel {
	if this.Attachments == nil {
		this.Attachments = make([]AttachmentModel, 0)
	}

	return &this.Attachments
}

func (this *APIEmailDetail) FillHTML() {
	html, ok := this.Meta["html"]
	if ok {
		this.HTML = html
	} else {
		text, ok := this.Meta["text"]
		if ok {
			this.HTML = strings.ReplaceAll(text, "\n", "<br>\n")
		}
	}
}

type APIListWrapper struct {
	Items      []APIEmailListItem `json:"items" gorm:"-"`
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
