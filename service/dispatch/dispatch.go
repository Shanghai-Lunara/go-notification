package dispatch

import (
	"flag"
	"fmt"
	"github.com/nevercase/go-notification/config"
	"github.com/nevercase/go-notification/service/dispatch/common"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

var (
	s *common.Service
)

func Init() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := config.Init(); err != nil {
		log.Fatal(fmt.Sprintf("conf.Init err:(%v)", err))
	}
	logPath := fmt.Sprintf("%s/../../log/%s", config.GetConfigPath(), config.GetDispatchLogFile())
	path, _ := filepath.Abs(logPath)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic("err: ", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Print("close log file err:", err)
		}
	}()
	log.SetOutput(file)
	s = common.New(config.GetConfig())
	signalHandler()
}

func signalHandler() {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-ch
		log.Printf("get a signal %s, stop the go-notification dispatch serivce \n", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			s.Close()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
