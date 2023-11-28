package pkg

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"mime"
	"net/http"
	"net/mail"
	"path/filepath"
	"strings"
)

const CHARSET_UTF8 = "UTF-8"

type SendGridParser struct {
	subjectDistChatSet  string
	htmlDistChatSet     string
	textDistChatSet     string
	toDistChatSet       string
	fromDistChatSet     string
	ccDistChatSet       string
	filenameDistChatSet string
	rawBody             map[string]string
	contacts            *ContactList
	attachments         []*AttachmentRaw
}

type ContactList struct {
	To   []*mail.Address
	From *mail.Address
	CC   []*mail.Address
}

func ConvertCharset(encoderName string, src string) (string, error) {
	if encoderName == CHARSET_UTF8 {
		return src, nil
	}
	reader := bytes.NewBufferString(src)
	out, err := charset.NewReaderLabel(encoderName, reader)
	if err != nil {
		return src, err
	}
	bin, err := ioutil.ReadAll(out)
	if err != nil {
		return src, err
	}
	return string(bin), nil
}

func (s *SendGridParser) Parse(request *http.Request) (IParser, error) {
	// 解析多部分表單，10 << 20 指定最大 10 MB
	err := request.ParseMultipartForm(10 << 20)
	if err != nil {
		Log.Error(err)
		return nil, err
	}
	//parse form data
	s.rawBody = map[string]string{}
	for key, list := range request.MultipartForm.Value {
		for _, v := range list {
			s.rawBody[key] = v
		}
	}
	s.parseMeta()
	s.parseContacts()
	s.parseAttachments(request)
	s.parseEmbed()

	return s, nil
}

func (s *SendGridParser) parseMeta() {
	// fill charset settings
	s.subjectDistChatSet = CHARSET_UTF8
	s.htmlDistChatSet = CHARSET_UTF8
	s.textDistChatSet = CHARSET_UTF8
	s.toDistChatSet = CHARSET_UTF8
	s.fromDistChatSet = CHARSET_UTF8
	s.ccDistChatSet = CHARSET_UTF8
	s.filenameDistChatSet = CHARSET_UTF8

	charsetsConfig := map[string]string{}
	charsetsJSONStr, ok := s.rawBody["charsets"]
	if ok {
		err := json.Unmarshal([]byte(charsetsJSONStr), &charsetsConfig)
		if err != nil {
			Log.Error(err)
		}
	}

	for key, value := range charsetsConfig {
		switch key {
		case "to":
			s.toDistChatSet = value
			break
		case "html":
			s.htmlDistChatSet = value
			break
		case "cc":
			s.ccDistChatSet = value
			break
		case "subject":
			s.subjectDistChatSet = value
			break
		case "from":
			s.fromDistChatSet = value
			break
		case "text":
			s.textDistChatSet = value
			break
		case "filename":
			s.filenameDistChatSet = value
			break
		}
	}

	//parse html
	html, ok := s.rawBody["html"]
	if ok {
		html, err := ConvertCharset(s.htmlDistChatSet, html)
		if err != nil {
			Log.Error(err)
		}
		s.rawBody["html"] = html
	}
	//parse text
	text, ok := s.rawBody["text"]
	if ok && s.textDistChatSet != CHARSET_UTF8 {
		text, err := ConvertCharset(s.textDistChatSet, text)
		if err != nil {
			Log.Error(err)
		}
		s.rawBody["text"] = text
	}
	//parse to
	to, ok := s.rawBody["to"]
	if ok && s.toDistChatSet != CHARSET_UTF8 {
		to, err := ConvertCharset(s.toDistChatSet, to)
		if err != nil {
			Log.Error(err)
		}
		s.rawBody["to"] = to
	}
	//parse from
	from, ok := s.rawBody["from"]
	if ok && s.fromDistChatSet != CHARSET_UTF8 {
		from, err := ConvertCharset(s.fromDistChatSet, from)
		if err != nil {
			Log.Error(err)
		}
		s.rawBody["from"] = from
	}
	//parse subject
	subject, ok := s.rawBody["subject"]
	if ok && s.subjectDistChatSet != CHARSET_UTF8 {
		subject, err := ConvertCharset(s.subjectDistChatSet, subject)
		if err != nil {
			Log.Error(err)
		}
		s.rawBody["subject"] = subject
	}
	//parse cc
	cc, ok := s.rawBody["cc"]
	if ok && s.ccDistChatSet != CHARSET_UTF8 {
		cc, err := ConvertCharset(s.ccDistChatSet, cc)
		if err != nil {
			Log.Error(err)
		}
		s.rawBody["cc"] = cc
	}
}

func (s *SendGridParser) parseContacts() {
	s.contacts = &ContactList{
		To:   []*mail.Address{},
		CC:   []*mail.Address{},
		From: nil,
	}
	fromStr, ok := s.rawBody["from"]
	if ok {
		from, err := mail.ParseAddress(fromStr)
		if err != nil {
			Log.Error(err)
		} else {
			s.contacts.From = from
		}
	}

	ccStr, ok := s.rawBody["cc"]
	if ok {
		cc, err := mail.ParseAddressList(ccStr)
		if err != nil {
			Log.Error(err)
		} else {
			s.contacts.CC = cc
		}
	}

	toStr, ok := s.rawBody["to"]
	if ok {
		to, err := mail.ParseAddressList(toStr)
		if err != nil {
			Log.Error(err)
		} else {
			s.contacts.To = to
		}
	}
}

func (s *SendGridParser) parseAttachments(request *http.Request) {
	s.attachments = []*AttachmentRaw{}

	contentIDMappings := map[string]string{}
	contentIdsStr, ok := s.rawBody["content-ids"]
	if ok {
		err := json.Unmarshal([]byte(contentIdsStr), &contentIDMappings)
		if err != nil {
			Log.Error(err)
		}
	}

	for fieldName, fileList := range request.MultipartForm.File {
		for _, header := range fileList {
			filename, err := ConvertCharset(s.filenameDistChatSet, header.Filename)
			if err != nil {
				Log.Error(err)
				continue
			}
			reader, err := header.Open()
			if err != nil {
				Log.Error(err)
				continue
			}
			mimeType := header.Header.Get("Content-Type")
			contentID := ""
			for k, v := range contentIDMappings {
				if strings.EqualFold(v, fieldName) {
					contentID = k
				}
			}
			s.attachments = append(s.attachments, &AttachmentRaw{
				FileName:  filename,
				MimeType:  mimeType,
				File:      reader,
				ContentID: contentID,
			})
		}
	}
}

func (s *SendGridParser) parseEmbed() {
	text, ok:= s.rawBody["text"]
	if ok && len(text) > 0 {
		result, err := ParseUUEncode(text)
		if err == nil && result != nil && len(result.Embeds) > 0 {
			for _, embed := range result.Embeds {

				mimeType := "application/octet-stream"
				mimeType = mime.TypeByExtension(filepath.Ext(embed.Name))

				s.attachments = append(s.attachments, &AttachmentRaw{
					FileName:  embed.Name,
					MimeType:  mimeType,
					File:      ioutil.NopCloser(embed.Data),
					ContentID: "",
				})
			}

			s.rawBody["text"] = result.SplitBody
		}
		if err != nil {
			Log.Error(err)
		}
	}
}

func (s *SendGridParser) GetSubject() string {
	subject, ok := s.rawBody["subject"]
	if !ok {
		return ""
	}

	return subject
}

func (s *SendGridParser) GetToList() []string {
	out := []string{}

	for _, to := range s.contacts.To {
		out = append(out, to.Address)
	}

	return out
}

func (s *SendGridParser) GetMate() map[string]string {
	return s.rawBody
}

func (s *SendGridParser) GetFrom() string {
	if s.contacts.From != nil {
		return s.contacts.From.Address
	}

	return ""
}

func (s *SendGridParser) GetAttachments() []*AttachmentRaw {
	return s.attachments
}
