package main

import (
	"context"
	"email2db/pkg"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	conf_path := flag.String("c", "conf.json", "config json file")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	conf, err := pkg.NewConfigFromLocal(*conf_path)
	if err != nil {
		pkg.Log.Error(err)
		conf = &pkg.Config{}
	}

	conf.MarginWithENV()

	pkg.Log.Debug("show config detail:")
	pkg.Log.Debug(conf.ToJSON())

	if conf.IsEnablePop3() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		ctx, cancel := context.WithCancel(context.Background())

		go pkg.StartCheckerWorker(conf, 30 * time.Minute, ctx)

		go func() {
			<-sigs
			cancel()
		}()
	}

	service := pkg.NewHttpService(conf)
	service.Start()
}
