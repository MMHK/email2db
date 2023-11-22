package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type MySQLConfig struct {
	DSN string `json:"dsn"`
}

type DBConfig struct {
	MySQL *MySQLConfig `json:"mysql"`
}

type ToList []string

func (ToList) GormDataType() string {
	return "json"
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

type MailModel struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	Subject   string `gorm:"column:subject" json:"subject"`
	From      string `gorm:"column:from" json:"from"`
	ReplyTo   string `gorm:"column:reply_to" json:"reply_to"`
	To        ToList `gorm:"column:to" gorm:"type:json" json:"to"`
	Meta      MetaObject `gorm:"column:meta" gorm:"type:json" json:"meta"`
	CreatedAt time.Time `gorm:"column:created_at" gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" gorm:"type:datetime" json:"updated_at"`
}

func (MailModel) TableName() string {
	return "em_email"
}

type AttachmentModel struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	Name   string `gorm:"column:name" json:"name"`
	Path      string `gorm:"column:path" json:"path"`
	ContentID   string `gorm:"column:content-id" json:"content-id"`
	MimeType   string `gorm:"column:mime_type" json:"mime_type"`
	CreatedAt time.Time `gorm:"column:created_at" gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" gorm:"type:datetime" json:"updated_at"`
}

func (AttachmentModel) TableName() string {
	return "em_attachment"
}

type IDBHelper interface {
	SaveMail(target *MailModel) (ID uint, err error)
	SaveAttachment(target *AttachmentModel) (ID uint, err error)
}

type DBHelper struct {
	connection *gorm.DB
}

func GetDBHelper(config *DBConfig) (IDBHelper, error) {
	if len(config.MySQL.DSN) > 0 {
		connection, err := gorm.Open(mysql.New(mysql.Config{
			DSN: config.MySQL.DSN,
			SkipInitializeWithVersion: false,
		}), &gorm.Config{})
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