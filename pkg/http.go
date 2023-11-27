package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type HttpConfig struct {
	Listen  string `json:"listen"`
	WebRoot string `json:"web_root"`
}

type APIStandardError struct {
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

type ParamsInPath struct {
	params map[string]string
}

func (this *ParamsInPath) GetInt(key string, val int) int {
	value, ok := this.params[key]
	if !ok {
		return val
	}
	out, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return val
	}
	return int(out)
}

func GetParamsInPath(request *http.Request) (*ParamsInPath) {
	vars := mux.Vars(request)
	return &ParamsInPath{
		params: vars,
	}
}

type ParamsInQuery ParamsInPath

func (this *ParamsInQuery) GetUint(key string, val uint) uint {
	value, ok := this.params[key]
	if !ok {
		return val
	}
	out, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return val
	}
	return uint(out)
}

func (this *ParamsInQuery) GetString(key string, val string) string {
	value, ok := this.params[key]
	if !ok {
		return val
	}

	return value
}

func GetParamsInQuery(request *http.Request) (*ParamsInQuery) {
	vars := request.URL.Query()
	p := make(map[string]string, 0)

	for k, ls := range vars {
		for _, v := range ls {
			_, ok := p[k]
			if ok {
				continue
			}
			p[k] = v
		}
	}

	return &ParamsInQuery{
		params: p,
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

func (this *HTTPService) getEmailList(writer http.ResponseWriter, request *http.Request) {
	params := GetParamsInQuery(request)
	pageSize := params.GetUint("pageSize", 10)
	page := params.GetUint("page", 10)
	search := params.GetString("s", "")

	db, err := GetDBHelper(this.conf.DB)
	if err != nil {
		Log.Error(err)
		this.ResponseError(err, writer)
		return
	}
	list := []APIEmailListItem{}
	pager, err := db.GetMailList(&list, search, pageSize, page)
	if err != nil {
		Log.Error(err)
		this.ResponseError(err, writer)
		return
	}

	this.ResponseJSON(&APIListResponse{
		APIStandardError: APIStandardError{
			Status: true,
		},
		Data: APIListWrapper{
			Items: list,
			Pagination: *pager,
		},
	}, writer)
}

func (this *HTTPService) getEmailDetail(writer http.ResponseWriter, request *http.Request) {
	params := GetParamsInPath(request)
	mail_id := params.GetInt("id", 0)
	if mail_id <= 0 {
		this.ResponseError(errors.New("Mail ID must be greater than zero"), writer)
		return
	}

	db, err := GetDBHelper(this.conf.DB)
	if err != nil {
		Log.Error(err)
		this.ResponseError(err, writer)
		return
	}

	result := APIEmailDetail{}
	err = db.GetMailDetail(&result, mail_id)
	if err != nil {
		Log.Error(err)
		this.ResponseError(err, writer)
		return
	}

	this.ResponseJSON(&APIEmailDetailResponse{
		APIStandardError: APIStandardError{
			Status: true,
		},
		Data: result,
	}, writer)
}

func (this *HTTPService) downloadAttachment(writer http.ResponseWriter, request *http.Request)  {
	params := GetParamsInPath(request)
	attachment_id := params.GetInt("id", 0)
	if attachment_id <= 0 {
		this.ResponseError(errors.New("Attachment ID must be greater than zero"), writer)
		return
	}

	db, err := GetDBHelper(this.conf.DB)
	if err != nil {
		Log.Error(err)
		this.ResponseError(err, writer)
		return
	}

	attachment, err := db.GetAttachment(attachment_id)
	if err != nil {
		Log.Error(err)
		this.ResponseError(err, writer)
		return
	}
	storage, err := GetStorage(this.conf.Storage)
	if err != nil {
		Log.Error(err)
		this.ResponseError(err, writer)
		return
	}

	url, err := storage.GetDownloadLink(attachment.Path)
	if err != nil {
		Log.Error(err)
		this.ResponseError(err, writer)
		return
	}

	http.Redirect(writer, request, url, 302)
}

func (s *HTTPService) Start() {
	rHandler := mux.NewRouter()

	rHandler.HandleFunc("/", s.RedirectSwagger)
	rHandler.HandleFunc("/api/mail", s.getEmailList)
	rHandler.HandleFunc("/api/mail/{id}", s.getEmailDetail)
	rHandler.HandleFunc("/api/attachment/{id}", s.downloadAttachment)
	rHandler.HandleFunc("/webhook/sendgrid", s.parseSendGridWebhook)
	rHandler.PathPrefix("/").Handler(http.StripPrefix("/",
		http.FileServer(http.Dir(fmt.Sprintf("%s", s.conf.WebRoot)))))
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
	server_error := &APIStandardError{Error: "handle not found!", Status: false}
	json_str, _ := json.Marshal(server_error)
	http.Error(writer, string(json_str), 404)
}

func (s *HTTPService) ResponseError(err error, writer http.ResponseWriter) {
	server_error := &APIStandardError{Error: err.Error(), Status: false}
	json_str, _ := json.Marshal(server_error)
	http.Error(writer, string(json_str), 200)
}


func (s *HTTPService) RedirectSwagger(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/swagger/", 301)
}
