package pkg

import (
	"email2db/tests"
	"io/ioutil"
	"net/mail"
	"os"
	"testing"
)

func TestParseUUEncode(t *testing.T) {
	raw, err := os.Open(tests.GetLocalPath("../tests/test.eml"))
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	defer raw.Close()

	email, err := mail.ReadMessage(raw)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	bin, err := ioutil.ReadAll(email.Body)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	text := string(bin)

	//t.Log(text)

	uu, err := ParseUUEncode(text)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(uu)
	t.Log("PASS")
}