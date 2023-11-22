package parser

import (
	"bytes"
	"email2db/pkg"
	"golang.org/x/net/html/charset"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const CHARSET_UTF8 = "UTF-8"

type SendGridParser struct {
	subjectDistChatSet string
	htmlDistChatSet string
	textDistChatSet string
	toDistChatSet string
	fromDistChatSet string
	ccDistChatSet string
	RawBody *map[string]string
}

func (s SendGridParser) Parse(request *http.Request) (pkg.IParser, error) {
	// 解析多部分表單，10 << 20 指定最大 10 MB
	err := request.ParseMultipartForm(10 << 20)
	if err != nil {
		pkg.Log.Error(err)
		return nil, err
	}
	//parse form data
	rawBody := map[string]string{}
	for key, list := range request.MultipartForm.Value {
		for _, v := range list {
			rawBody[key] = v
		}
	}
	s.RawBody = &rawBody

	// fill charset settings
	s.subjectDistChatSet = CHARSET_UTF8
	s.htmlDistChatSet = CHARSET_UTF8
	s.textDistChatSet = CHARSET_UTF8
	s.toDistChatSet = CHARSET_UTF8
	s.fromDistChatSet = CHARSET_UTF8
	s.ccDistChatSet = CHARSET_UTF8

	charsetsConfig := map[string]string{}
	charsetsJSONStr, ok := rawBody["charsets"]
	if ok {
		err = json.Unmarshal([]byte(charsetsJSONStr), &charsetsConfig)
		if err != nil {
			pkg.Log.Error(err)
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
		}
	}

	//parse html
	html, ok := rawBody["html"]
	if ok {
		html, err = ConvertCharset(s.htmlDistChatSet, html)
		if err != nil {
			pkg.Log.Error(err)
		}
		rawBody["html"] = html
	}
	//parse text
	text, ok := rawBody["text"]
	if ok && s.textDistChatSet != CHARSET_UTF8 {
		text, err = ConvertCharset(s.textDistChatSet, text)
		if err != nil {
			pkg.Log.Error(err)
		}
		rawBody["text"] = text
	}
	//parse to
	to, ok := rawBody["to"]
	if ok && s.toDistChatSet != CHARSET_UTF8 {
		to, err = ConvertCharset(s.toDistChatSet, to)
		if err != nil {
			pkg.Log.Error(err)
		}
		rawBody["to"] = to
	}
	//parse from
	from, ok := rawBody["from"]
	if ok && s.fromDistChatSet != CHARSET_UTF8 {
		from, err = ConvertCharset(s.fromDistChatSet, from)
		if err != nil {
			pkg.Log.Error(err)
		}
		rawBody["from"] = from
	}
	//parse subject
	subject, ok := rawBody["subject"]
	if ok && s.subjectDistChatSet != CHARSET_UTF8 {
		subject, err = ConvertCharset(s.subjectDistChatSet, subject)
		if err != nil {
			pkg.Log.Error(err)
		}
		rawBody["subject"] = subject
	}
	//parse cc
	cc, ok := rawBody["cc"]
	if ok && s.ccDistChatSet != CHARSET_UTF8 {
		cc, err = ConvertCharset(s.ccDistChatSet, cc)
		if err != nil {
			pkg.Log.Error(err)
		}
		rawBody["cc"] = cc
	}

	panic("implement me")
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
	bytes, err := ioutil.ReadAll(out)
	if err != nil {
		return src, err
	}
	return string(bytes), nil
}

func (s SendGridParser) GetSubject() string {
	panic("implement me")
}

func (s SendGridParser) GetToList() *[]string {
	panic("implement me")
}

func (s SendGridParser) GetMate() *map[string]string {
	panic("implement me")
}

func (s SendGridParser) GetReplyTo() string {
	panic("implement me")
}

func (s SendGridParser) GetFrom() string {
	panic("implement me")
}

func (s SendGridParser) GetAttachments() *[]pkg.AttachmentRaw {
	panic("implement me")
}


