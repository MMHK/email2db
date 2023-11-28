package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pedroalbanese/uuencode"
	"io"
	"os"
	"regexp"
	"strings"
)

const UUENCODE_RULE = `begin ([0-9]{3}) ([^ \n\r]+)([\r\n]+)([\s\S\n\r]+)end`

type UUEncode struct {
	Data io.Reader
	Name string
	Mode os.FileMode
}

type UUEncodeResult struct {
	Embeds []*UUEncode
	SplitBody string
}

func ParseUUEncode(text string) (*UUEncodeResult, error) {
	// 定義正則表達式以匹配 uuencode 部分
	re := regexp.MustCompile(UUENCODE_RULE)

	// 匹配 uuencode 部分
	matchGroup := re.FindAllStringSubmatch(text, -1)
	if len(matchGroup) < 0 {
		err := errors.New("no uuencode section found")
		Log.Error(err)
		return nil, err
	}

	list := make([]*UUEncode, 0)
	for _, match := range matchGroup {
		if len(match) < 2 {
			err := errors.New("no uuencode section found")
			Log.Error(err)
			continue
		}

		reader := strings.NewReader(fmt.Sprintf("%s\n", match[0]))
		buf := new(bytes.Buffer)
		//Log.Info(buf.String())

		decoder := uuencode.NewReader(reader, nil)
		_, err := io.Copy(buf, decoder)
		if err != nil {
			Log.Error(err)
			continue
		}

		filename, ok := decoder.File()
		if !ok {
			err := errors.New("filename not found")
			Log.Error(err)
			continue
		}
		mode, ok := decoder.Mode()
		if !ok {
			err := errors.New("mode not found")
			Log.Error(err)
			continue
		}

		list = append(list, &UUEncode{
			Data: buf,
			Name: strings.TrimSpace(filename),
			Mode: mode,
		})
	}

	if len(list) <= 0 {
		return nil, errors.New("no uuencode section found")
	}

	// 返回解析結果
	return &UUEncodeResult{
		Embeds: list,
		SplitBody: re.ReplaceAllString(text, ""),
	}, nil
}