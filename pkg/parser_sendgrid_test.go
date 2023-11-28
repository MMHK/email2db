package pkg

import (
	"email2db/tests"
	"encoding/json"
	"io/ioutil"
	"net/mail"
	"os"
	"path/filepath"
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

func TestWriteFile(t *testing.T) {
	filename := filepath.ToSlash(tests.GetLocalPath("../tests/20231127190028.196167_ZPP0165992ZC_2023_000_DC_062.pdf"))
	err := ioutil.WriteFile(filename, []byte(""), 0777)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}

func TestSupportUUEncode(t *testing.T) {
	raw, err := os.Open(tests.GetLocalPath("../tests/request.json"))
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	defer raw.Close()

	rawBody := make(map[string]string, 0)
	jsondecoder := json.NewDecoder(raw)
	err = jsondecoder.Decode(&rawBody)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	//t.Log(tests.ToJSON(rawBody))

	sendgrid := &SendGridParser{
		rawBody: rawBody,
	}

	sendgrid.parseMeta()
	sendgrid.parseEmbed()

	t.Log(sendgrid.attachments)
	t.Log(tests.ToJSON(sendgrid.rawBody))

	//attachments := sendgrid.GetAttachments()
	//for _, attachment := range attachments {
	//	t.Log(attachment.MimeType)
	//	filePath := tests.GetLocalPath(fmt.Sprintf(`../tests/%s`, attachment.FileName))
	//	filePath = filepath.ToSlash(filePath)
	//	buf := new(bytes.Buffer)
	//	_, err = buf.ReadFrom(attachment.File)
	//	if err != nil {
	//		t.Error(err)
	//		continue
	//	}
	//	err := ioutil.WriteFile(filePath, buf.Bytes(), 0777)
	//	if err != nil {
	//		t.Error(err)
	//		continue
	//	}
	//
	//	defer attachment.File.Close()
	//}

	t.Log("PASS")

}