package pkg

import (
	"net/mail"
	"testing"
)

func TestConvertCharset(t *testing.T) {
	raw := "\xa4\xa4\xa4\xe5\xc1c\xc5\xe9"

	t.Log(raw)

	dist, err := ConvertCharset("big5", raw)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(dist)
	t.Log("PASS")
}

func TestParseEmailAddress(t *testing.T) {
	raw := `fwd@uat-zurich.driver.com.hk`

	addrList, err := mail.ParseAddressList(raw)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	for _, contact := range addrList {
		t.Log(contact.Name)
		t.Log(contact.Address)
	}
	t.Log("PASS")
}