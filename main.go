package email2db

import (
	"email2db/pkg"
	"flag"
	"runtime"
)

func main() {
	conf_path := flag.String("c", "conf.json", "config json file")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	conf, err := pkg.NewConfigFromLocal(*conf_path)
	if err != nil {
		pkg.Log.Error(err)

		conf := &pkg.Config{}
		conf.MarginWithENV()
	}

	service := pkg.NewHttpService(conf)
	service.Start()
}
