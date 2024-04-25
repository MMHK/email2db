package pkg

import (
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
		list, err := reader.PullMailList(conn, 3)
		if err != nil {
			return err
		}

		for _, mail := range list {
			t.Logf("%+v", mail.GetMate())
			t.Logf("%+v", mail.GetAttachments())
		}

		return nil
	})

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	t.Log("PASS")
}
