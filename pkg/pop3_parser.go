package pkg

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/DusanKasan/parsemail"
	"github.com/knadh/go-pop3"
	"io"
	"mime"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type IPop3Config interface {
	GetHost() string
	GetPort() int
	GetEmail() string
	GetPwd() string
	GetTLS() bool
}

type ZohoPopConfig struct {
	Email     string `json:"email"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	AppSecret string `json:"app_secret"`
	TLS       bool   `json:"tls"`
}

func LoadPop3ConfigWithEnv() IPop3Config {
	intVal, err := strconv.Atoi(os.Getenv("ZOHO_POP3_PORT"))
	if err != nil {
		intVal = 995
	}
	boolVal, err := strconv.ParseBool(os.Getenv("ZOHO_POP3_TLS"))
	if err != nil {
		boolVal = false
	}

	return &ZohoPopConfig{
		Host:      os.Getenv("ZOHO_POP3_HOST"),
		Email:     os.Getenv("ZOHO_EMAIL"),
		Port:      intVal,
		TLS:       boolVal,
		AppSecret: os.Getenv("ZOHO_APP_SECRET"),
	}
}

func (this *ZohoPopConfig) GetHost() string {
	return this.Host
}

func (this *ZohoPopConfig) GetPort() int {
	return this.Port
}

func (this *ZohoPopConfig) GetEmail() string {
	return this.Email
}

func (this *ZohoPopConfig) GetPwd() string {
	return this.AppSecret
}

func (this *ZohoPopConfig) GetTLS() bool {
	return this.TLS
}

type Pop3Parser struct {
	rawBody     map[string]string
	contacts    *ContactList
	attachments []*AttachmentRaw
	Date        time.Time
}

func (s *Pop3Parser) Parse(request *http.Request) (IParser, error) {
	return s, nil
}

func NewPop3ParserFromRaw(reader io.Reader) (*Pop3Parser, error) {
	s := &Pop3Parser{
		rawBody:     map[string]string{},
		contacts:    new(ContactList),
		attachments: []*AttachmentRaw{},
	}

	eml, err := parsemail.Parse(reader)
	if err != nil {
		return nil, err
	}
	s.rawBody["subject"] = eml.Subject
	fromList := ""
	for i, addr := range eml.From {
		if i == 0 {
			fromList = addr.Address
			continue
		}
		fromList = fmt.Sprintf("%s,%s", fromList, addr.Address)
	}
	s.contacts.From = eml.From[0]
	s.rawBody["from"] = fromList

	contactList := []*mail.Address{}
	toList := ""
	for i, addr := range eml.To {
		if i == 0 {
			toList = addr.Address
			continue
		}
		toList = fmt.Sprintf("%s,%s", toList, addr.Address)
		contactList = append(contactList, addr)
	}
	s.contacts.To = eml.To
	s.rawBody["to"] = toList

	ccList := ""
	for i, addr := range eml.Cc {
		if i == 0 {
			ccList = addr.Address
			continue
		}
		ccList = fmt.Sprintf("%s,%s", ccList, addr.Address)
	}
	s.contacts.CC = eml.Cc
	s.rawBody["cc"] = ccList

	replyList := ""
	for i, addr := range eml.ReplyTo {
		if i == 0 {
			replyList = addr.Address
			continue
		}
		replyList = fmt.Sprintf("%s,%s", replyList, addr.Address)
	}
	s.rawBody["ReplyTo"] = replyList

	s.rawBody["text"] = eml.TextBody
	s.rawBody["html"] = eml.HTMLBody

	s.Date = eml.Date

	r := bufio.NewReader(strings.NewReader(eml.TextBody))
	line, _ := r.ReadString('\n')
	// 檢查是否 base64 編碼的body
	if len(line) == 78 {
		rawBody, err := base64.StdEncoding.DecodeString(eml.TextBody)
		if err == nil {
			s.rawBody["text"] = string(rawBody)
		}
	}

	r = bufio.NewReader(strings.NewReader(eml.HTMLBody))
	line, _ = r.ReadString('\n')
	// 檢查是否 base64 編碼的body
	if len(line) == 78 {
		rawBody, err := base64.StdEncoding.DecodeString(eml.HTMLBody)
		if err == nil {
			s.rawBody["html"] = string(rawBody)
		}
	}

	s.rawBody["MessageID"] = eml.MessageID

	fileList := []*AttachmentRaw{}
	for _, file := range eml.Attachments {
		fileList = append(fileList, &AttachmentRaw{
			FileName:  file.Filename,
			File:      io.NopCloser(file.Data),
			MimeType:  file.ContentType,
			ContentID: "",
		})
	}
	for _, file := range eml.EmbeddedFiles {
		fileList = append(fileList, &AttachmentRaw{
			FileName:  file.CID,
			File:      io.NopCloser(file.Data),
			MimeType:  file.ContentType,
			ContentID: file.CID,
		})
	}

	result, err := ParseUUEncode(eml.TextBody)
	if err == nil && result != nil && len(result.Embeds) > 0 {
		for _, embed := range result.Embeds {

			mimeType := "application/octet-stream"
			mimeType = mime.TypeByExtension(filepath.Ext(embed.Name))

			fileList = append(fileList, &AttachmentRaw{
				FileName:  embed.Name,
				MimeType:  mimeType,
				File:      io.NopCloser(embed.Data),
				ContentID: "",
			})
		}

		s.rawBody["text"] = result.SplitBody
	}
	if err != nil {
		Log.Error(err)
	}

	s.attachments = fileList

	return s, nil
}

func (s *Pop3Parser) GetSubject() string {
	return s.rawBody["subject"]
}

func (s *Pop3Parser) GetToList() []string {
	out := []string{}

	for _, addr := range s.contacts.To {
		out = append(out, addr.Address)
	}

	return out
}

func (s *Pop3Parser) GetMate() map[string]string {
	return s.rawBody
}

func (s *Pop3Parser) GetFrom() string {
	return s.rawBody["from"]
}

func (s *Pop3Parser) GetDate() time.Time {
	return s.Date
}

func (s *Pop3Parser) GetAttachments() []*AttachmentRaw {
	return s.attachments
}

type Pop3Reader struct {
	config IPop3Config
	client *pop3.Client
}

func (this *Pop3Reader) StartConnection(callback func(conn *pop3.Conn) error) error {
	connection, err := this.client.NewConn()
	if err != nil {
		Log.Error(err)
		return err
	}
	defer connection.Quit()

	// Authenticate.
	if err := connection.Auth(this.config.GetEmail(), this.config.GetPwd()); err != nil {
		Log.Error(err)
		return err
	}

	return callback(connection)
}

func (this *Pop3Reader) GetCounter(conn *pop3.Conn) (int, int, error) {
	return conn.Stat()
}

func (this *Pop3Reader) EachMail(conn *pop3.Conn, limit int, callback func(parser *Pop3Parser)) error {
	maiList, err := conn.List(0)
	if err != nil {
		return err
	}

	total := len(maiList)
	end := total - limit
	if end < 0 {
		end = 0
	}
	for i := total - 1; i >= end; i-- {
		msgID := maiList[i]
		buf, err := conn.RetrRaw(msgID.ID)

		parser, err := NewPop3ParserFromRaw(buf)
		if err != nil {
			Log.Error(err)
			continue
		}
		callback(parser)
	}

	return nil
}

func (this *Pop3Reader) PullMailList(conn *pop3.Conn, limit int) ([]*Pop3Parser, error) {
	parserList := make([]*Pop3Parser, 0)
	maiList, err := conn.List(0)
	if err != nil {
		return parserList, err
	}

	total := len(maiList)
	end := total - limit
	if end < 0 {
		end = 0
	}

	for i := total - 1; i >= end; i-- {
		msgID := maiList[i]
		buf, err := conn.RetrRaw(msgID.ID)
		if err != nil {
			Log.Error(err)
			continue
		}

		//os.WriteFile(tests.GetLocalPath(fmt.Sprintf("../tests/%s.eml", MakeUUID())), buf.Bytes(), os.ModePerm)

		parser, err := NewPop3ParserFromRaw(buf)
		if err != nil {
			Log.Error(err)
			continue
		}
		parserList = append(parserList, parser)
	}

	return parserList, nil
}

func NewPop3Reader(config IPop3Config) *Pop3Reader {
	client := pop3.New(pop3.Opt{
		Host:       config.GetHost(),
		Port:       config.GetPort(),
		TLSEnabled: config.GetTLS(),
	})

	return &Pop3Reader{
		config: config,
		client: client,
	}
}

func RunPop3Checker(limit int, pop3Config IPop3Config, dbConfig *DBConfig, storageConf *StorageConfig) {
	reader := NewPop3Reader(pop3Config)

	helper, err := GetDBHelper(dbConfig)
	if err != nil {
		Log.Error(err)
		return
	}
	defer helper.Close()

	Log.Debug("connect to pop3 server")
	reader.StartConnection(func(conn *pop3.Conn) error {
		defer Log.Debug("disconnect to pop3 server")

		return reader.EachMail(conn, limit, func(parser *Pop3Parser) {
			Log.Debug("pulling email from pop3 server")

			mail_id := uint(0)

			meta := parser.GetMate()
			ReplyTo, ok := meta["ReplyTo"]
			if !ok || ReplyTo == "" {
				ReplyTo = parser.GetFrom()
			}

			mail := &MailModel{
				Subject: parser.GetSubject(),
				From:    parser.GetFrom(),
				To:      parser.GetToList(),
				Meta:    meta,
				ReplyTo: ReplyTo,
				CreatedAt: parser.GetDate(),
			}

			if !helper.Exist(mail) {
				Log.Debug("start save mail to DB")
				mail_id, err = helper.SaveMail(mail)
				if err != nil {
					Log.Error(err)
					return
				}
				Log.Debugf("Mail Save to DB, ID:%d", mail_id)
			} else {
				Log.Debug("mail already exist")
			}

			attachments := parser.GetAttachments()
			if len(attachments) > 0 && mail_id > 0 {
				Log.Debug("start upload attachments to S3")
				storage, err := GetStorage(storageConf)
				if err != nil {
					Log.Error(err)
					return
				}

				for _, attachment := range attachments {
					distKey := fmt.Sprintf("%s%s", MakeUUID(), filepath.Ext(attachment.FileName))
					Log.Debugf("uploading attachment[%s] to S3", attachment.FileName)
					remotePath, _, err := storage.PutStream(attachment.File, distKey, &UploadOptions{
						ContentType: attachment.MimeType,
					})
					defer attachment.File.Close()
					if err != nil {
						Log.Error(err)
						continue
					}
					Log.Debugf("start save attachment[%s] to DB", attachment.FileName)
					_, err = helper.SaveAttachment(&AttachmentModel{
						Name:      attachment.FileName,
						Path:      remotePath,
						ContentID: attachment.ContentID,
						EmailID:   mail_id,
						MimeType:  attachment.MimeType,
					})
					if err != nil {
						Log.Error(err)
						continue
					}
				}
			} else {
				Log.Debug("no attachments or mail already exist")
			}
		})
	})

	helper.Close()
}


func StartCheckerWorker(conf *Config, sleepTime time.Duration, ctx context.Context) {
	Log.Info("Start Checker Worker")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker exiting.")
			return
		default:
			RunPop3Checker(15, conf.Pop3Config, conf.DB, conf.Storage)
			time.Sleep(sleepTime)
		}
	}
}
