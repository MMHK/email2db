package pkg

import (
	"email2db/tests"
	"github.com/knadh/go-pop3"
	"os"
	"strconv"
	"testing"
)

func GetTestPop3Config() IPop3Config {
	intVal, err := strconv.Atoi(os.Getenv("ZOHO_POP3_PORT"))
	if err != nil {
		intVal = 995
	}
	boolVal, err := strconv.ParseBool(os.Getenv("ZOHO_POP3_TLS"))
	if err != nil {
		boolVal = false
	}

	return &ZohoPopConfig{
		Host: os.Getenv("ZOHO_POP3_HOST"),
		Email: os.Getenv("ZOHO_EMAIL"),
		Port: intVal,
		TLS: boolVal,
		AppSecret: os.Getenv("ZOHO_APP_SECRET"),
	}
}

func TestPop3Reader_GetCounter(t *testing.T) {
	reader := NewPop3Reader(GetTestPop3Config())

	err := reader.StartConnection(func(conn *pop3.Conn) error {
		count, size, err := reader.GetCounter(conn)
		if err != nil {
			return err
		}

		t.Log("total messages=", count, "size=", size)

		return nil
	})

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	t.Log("PASS")
}

func TestPop3Reader_PullMailList(t *testing.T) {
	reader := NewPop3Reader(GetTestPop3Config())

	err := reader.StartConnection(func(conn *pop3.Conn) error {
		list, err := reader.PullMailList(conn, 5)
		if err != nil {
			return err
		}

		for _, mail := range list {
			//meta := mail.GetMate()
			//MessageID, ok := meta["MessageID"];
			//if !ok {
			//	MessageID = MakeUUID()
			//}
			//html, ok := meta["html"]
			//if ok {
			//	t.Log(html)
			//}

			t.Log("-------Subject-------");
			t.Logf("%+v", mail.GetSubject())
			t.Log("-------Subject-------");
			t.Log("-------To-------");
			t.Logf("%+v", mail.GetToList())
			t.Log("-------To-------");
			t.Log("-------Attachments-------");
			for _, file := range mail.GetAttachments() {
				t.Logf("%s, type: %s", file.FileName, file.MimeType)
			}
			t.Log("-------Attachments-------");
		}

		return nil
	})

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	t.Log("PASS")
}

func TestPop3Parser_EachMail(t *testing.T) {
	reader := NewPop3Reader(GetTestPop3Config())

	err := reader.StartConnection(func(conn *pop3.Conn) error {
		return reader.EachMail(conn, 40, func(parser *Pop3Parser) {
			t.Log("-------Subject-------");
			t.Logf("%+v", parser.GetSubject())
			t.Log("-------Subject-------");
			t.Log("-------To-------");
			t.Logf("%+v", parser.GetToList())
			t.Log("-------To-------");
			t.Log("-------Attachments-------");
			for _, file := range parser.GetAttachments() {
				t.Logf("%s, type: %s", file.FileName, file.MimeType)
			}
			t.Log("-------Attachments-------");
		})
	})

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	t.Log("PASS")
}

func TestPop3Parser_SaveToDB(t *testing.T) {
	dsn := loadMySQLDSN()

	helper, err := GetDBHelper(&DBConfig{
		&MySQLConfig{
			DSN: dsn,
		},
	})
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	reader := NewPop3Reader(GetTestPop3Config())

	err = reader.StartConnection(func(conn *pop3.Conn) error {
		return reader.EachMail(conn, 1, func(parser *Pop3Parser) {
			meta := parser.GetMate()
			ReplyTo, ok := meta["ReplyTo"]
			if !ok {
				ReplyTo = parser.GetFrom()
			}

			mail := &MailModel{
				Subject: parser.GetSubject(),
				From: parser.GetFrom(),
				To: parser.GetToList(),
				Meta: meta,
				ReplyTo: ReplyTo,
			}

			id, err := helper.SaveMail(mail)
			if err != nil {
				t.Error(err)
				return
			}
			t.Log(id)
		})
	})

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	t.Log("PASS")
}

func TestPop3Parser_DBExist(t *testing.T) {
	dsn := loadMySQLDSN()

	t.Log(dsn)

	helper, err := GetDBHelper(&DBConfig{
		&MySQLConfig{
			DSN: dsn,
		},
	})
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	reader := NewPop3Reader(GetTestPop3Config())

	err = reader.StartConnection(func(conn *pop3.Conn) error {
		return reader.EachMail(conn, 1, func(parser *Pop3Parser) {
			meta := parser.GetMate()
			ReplyTo, ok := meta["ReplyTo"]
			if !ok {
				ReplyTo = parser.GetFrom()
			}

			mail := &MailModel{
				Subject: parser.GetSubject(),
				From: parser.GetFrom(),
				To: parser.GetToList(),
				Meta: meta,
				ReplyTo: ReplyTo,
			}

			exist := helper.Exist(mail)
			t.Log(exist)
		})
	})

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	t.Log("PASS")
}

func TestRunPop3Checker(t *testing.T) {
	conf, err := NewConfigFromLocal(tests.GetLocalPath("../config.json"))
	if err != nil {
		t.Error(err)

		conf = &Config{}
	}

	conf.MarginWithENV()

	RunPop3Checker(15, conf.Pop3Config, conf.DB, conf.Storage)

	t.Log("PASS")
}

