package pkg

import "github.com/op/go-logging"

var Log = logging.MustGetLogger("email2db")

func init() {
	format := logging.MustStringFormatter(
		`Email2DB %{color} %{shortfunc} %{level:.4s} %{shortfile}
%{id:03x}%{color:reset} %{message}`,
	)
	logging.SetFormatter(format)
}
