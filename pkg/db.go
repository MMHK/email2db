package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/url"
	"os"
	"time"
)

type MySQLConfig struct {
	DSN string `json:"dsn"`
}

type DBConfig struct {
	MySQL *MySQLConfig `json:"mysql"`
}

func LoadMySQLDSNWithENV() string {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbname := os.Getenv("MYSQL_DATABASE")
	username := os.Getenv("MYSQL_USERNAME")
	pwd := os.Getenv("MYSQL_PASSWORD")
	tz := os.Getenv("TZ")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		url.PathEscape(username),
		url.PathEscape(pwd),
		host, port, url.PathEscape(dbname),
		url.QueryEscape(tz))
}

type ToList []string

func (ToList) GormDataType() string {
	return "json"
}

func (this *ToList) Scan(value interface{}) error {
	bin, ok := value.([]byte)
	if !ok {
		this = &ToList{}
		return nil
	}
	return json.Unmarshal(bin, this)
}

func (this ToList) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	jsonStr := "[]"
	jsonString, err := json.Marshal(this)
	if err != nil {
		Log.Error(err)
	} else {
		jsonStr = string(jsonString)
	}

	return clause.Expr{
		SQL:  "?",
		Vars: []interface{}{jsonStr},
	}
}

type MetaObject map[string]string

func (MetaObject) GormDataType() string {
	return "json"
}

func (this MetaObject) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	jsonStr := "[]"
	jsonString, err := json.Marshal(this)
	if err != nil {
		Log.Error(err)
	} else {
		jsonStr = string(jsonString)
	}

	return clause.Expr{
		SQL:  "?",
		Vars: []interface{}{jsonStr},
	}
}

type IMailDetail interface {
	GetAttachments() (*[]AttachmentModel)
}

type MailModel struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Subject   string     `gorm:"column:subject" json:"subject"`
	From      string     `gorm:"column:from" json:"from"`
	ReplyTo   string     `gorm:"column:reply_to" json:"reply_to"`
	To        ToList     `gorm:"column:to" gorm:"type:json" json:"to"`
	Meta      MetaObject `gorm:"column:meta" gorm:"type:json" json:"meta"`
	CreatedAt time.Time  `gorm:"column:created_at" gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" gorm:"type:datetime" json:"updated_at"`
	Attachments []AttachmentModel `gorm:"foreignKey:EmailID;references:ID;"`
}

func (MailModel) TableName() string {
	return "em_email"
}

func (this *MailModel) GetAttachments() (*[]AttachmentModel) {
	if this.Attachments == nil {
		this.Attachments = make([]AttachmentModel, 0)
	}

	return &this.Attachments
}

type AttachmentModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Path      string    `gorm:"column:path" json:"path"`
	ContentID string    `gorm:"column:content-id" json:"content-id"`
	EmailID   uint      `gorm:"column:email_id" json:"email_id"`
	MimeType  string    `gorm:"column:mime_type" json:"mime_type"`
	CreatedAt time.Time `gorm:"column:created_at" gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" gorm:"type:datetime" json:"updated_at"`
}

func (AttachmentModel) TableName() string {
	return "em_attachment"
}

type IDBHelper interface {
	SaveMail(target *MailModel) (ID uint, err error)
	SaveAttachment(target *AttachmentModel) (ID uint, err error)
	GetMailList(list interface{}, search string, pageSize uint, currentPage uint) (*Pagination, error)
	GetMailDetail(target IMailDetail, MailID int) (error)
	GetAttachment(id int) (*AttachmentModel, error)
	Close() (error)
}

type DBHelper struct {
	connection *gorm.DB
}

func GetDBHelper(config *DBConfig) (IDBHelper, error) {
	if len(config.MySQL.DSN) > 0 {
		connection, err := gorm.Open(mysql.New(mysql.Config{
			DSN:                       config.MySQL.DSN,
			SkipInitializeWithVersion: false,
		}), &gorm.Config{
			//Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return nil, err
		}
		return &DBHelper{connection}, nil
	}
	return nil, errors.New("database config not found")
}

func (this *DBHelper) SaveMail(target *MailModel) (uint, error) {
	result := this.connection.Create(target)
	if result.Error != nil {
		return 0, result.Error
	}
	return target.ID, nil
}

func (this *DBHelper) SaveAttachment(target *AttachmentModel) (uint, error) {
	result := this.connection.Create(target)
	if result.Error != nil {
		return 0, result.Error
	}
	return target.ID, nil
}

func (this *DBHelper) GetMailList(list interface{}, search string, pageSize uint, currentPage uint) (*Pagination, error) {
	tx := this.connection.Model(&MailModel{})

	if len(search) > 0 {
		s := fmt.Sprintf(`%%%s%%`, search)
		tx = tx.Where("`subject` LIKE ? OR `from` LIKE ? OR `to` LIKE ?", s, s, s)
	}

	var counter int64 = 0
	err := tx.Count(&counter).Error
	if err != nil {
		Log.Error(err)
		return nil, err
	}
	if counter < 0 {
		counter = 0
	}
	pager := NewPagination(uint(counter), currentPage, pageSize)
	tx = tx.Limit(int(pageSize)).Offset(int((pager.CurrentPage() - 1) * pageSize)).Order("`id` desc")

	err = tx.Find(list).Error
	if err != nil {
		Log.Error(err)
		return nil, err
	}
	return pager, nil
}

func (this *DBHelper) GetMailDetail(target IMailDetail, MailID int) (error) {
	tx := this.connection.Model(&MailModel{}).Where("id = ?", MailID)

	err := tx.First(&target).Error
	if err != nil {
		Log.Error(err)
		return err
	}

	attachments := target.GetAttachments()
	err = this.connection.Model(&AttachmentModel{}).Where("email_id=?", MailID).Find(attachments).Error
	if err != nil {
		Log.Error(err)
	}
	return nil
}

func (this *DBHelper) GetAttachment(id int) (*AttachmentModel, error) {
	tx := this.connection.Model(&AttachmentModel{}).Where("id = ?", id)

	target := AttachmentModel{}
	err := tx.First(&target).Error
	if err != nil {
		Log.Error(err)
		return nil, err
	}
	return &target, nil
}

func (this *DBHelper) Close() error {
	db, err := this.connection.DB()
	if err != nil {
		return err
	}
	return db.Close()
}