package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
	"strings"
)

type HttpConfig struct {
	Listen  string `json:"listen"`
	WebRoot string `json:"web_root"`
}

type ServiceError struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

type HTTPService struct {
	conf *Config
}

func NewHttpService(conf *Config) *HTTPService {
	return &HTTPService{
		conf: conf,
	}
}

func (this *HTTPService) parseSendGridWebhook(writer http.ResponseWriter, request *http.Request) {
	Log.Debug("new http request incoming")
	Log.Debugf("%+v", request)

	if strings.EqualFold(request.Method, "POST") {
		Log.Debug("start parse send grid webhook")
		sendgrid, err := GetSendGridParser(request)
		if err != nil {
			this.ResponseError(err, writer)
			return
		}
		Log.Debug("end parse send grid webhook")
		go (func() {
			Log.Debug("start save to DB")
			db, err := GetDBHelper(this.conf.DB)
			if err != nil {
				Log.Error(err)
				return
			}
			Log.Debug("start save mail to DB")
			mail_id, err := db.SaveMail(&MailModel{
				Subject: sendgrid.GetSubject(),
				From: sendgrid.GetFrom(),
				ReplyTo: sendgrid.GetFrom(),
				To: sendgrid.GetToList(),
				Meta: sendgrid.GetMate(),
			})
			if err != nil {
				Log.Error(err)
				return
			}
			Log.Debug("start upload attachments to S3")
			storage, err := GetStorage(this.conf.Storage)
			if err != nil {
				Log.Error(err)
				return
			}
			attachments := sendgrid.GetAttachments()
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
				_, err = db.SaveAttachment(&AttachmentModel{
					Name: attachment.FileName,
					Path: remotePath,
					ContentID: attachment.ContentID,
					EmailID: mail_id,
					MimeType: attachment.MimeType,
				})
				if err != nil {
					Log.Error(err)
					continue
				}
			}
			Log.Debug("end upload attachments to S3")
			Log.Debug("end save to DB")
		})()

		this.ResponseJSON(&map[string]string{
			"msg": "OK",
		}, writer)
		return
	}

	this.ResponseError(errors.New("invalid request"), writer)
}

func (s *HTTPService) Start() {
	rHandler := mux.NewRouter()

	rHandler.HandleFunc("/", s.RedirectSwagger)
	rHandler.HandleFunc("/webhook/sendgrid", s.parseSendGridWebhook)
	rHandler.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/",
		http.FileServer(http.Dir(fmt.Sprintf("%s/swagger", s.conf.WebRoot)))))
	rHandler.NotFoundHandler = http.HandlerFunc(s.NotFoundHandle)

	Log.Info("http service starting")
	Log.Infof("Please open http://%s\n", s.conf.Listen)
	err := http.ListenAndServe(s.conf.Listen, rHandler)
	if err != nil {
		Log.Error(err)
	}
}

func (s *HTTPService) ResponseJSON(source interface{}, writer http.ResponseWriter) {
	json_str, err := json.Marshal(source)
	if err != nil {
		s.ResponseError(err, writer)
	}
	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, string(json_str))
}

func (s *HTTPService) NotFoundHandle(writer http.ResponseWriter, request *http.Request) {
	server_error := &ServiceError{Error: "handle not found!", Status: false}
	json_str, _ := json.Marshal(server_error)
	http.Error(writer, string(json_str), 404)
}

func (s *HTTPService) ResponseError(err error, writer http.ResponseWriter) {
	server_error := &ServiceError{Error: err.Error(), Status: false}
	json_str, _ := json.Marshal(server_error)
	http.Error(writer, string(json_str), 200)
}


func (s *HTTPService) RedirectSwagger(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/swagger", 301)
}
