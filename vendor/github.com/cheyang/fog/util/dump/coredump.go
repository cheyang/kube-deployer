package dump

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
)

func StackTrace(all bool) string {
	buf := make([]byte, 10240)

	for {
		size := runtime.Stack(buf, all)

		if size == len(buf) {
			buf = make([]byte, len(buf)<<1)
			continue
		}
		break

	}

	return string(buf)
}

func InstallCoreDumpGenerator() {

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGQUIT)

	go func() {
		for {
			sig := <-signals

			switch sig {
			case syscall.SIGQUIT:
				t := time.Now()
				timestamp := fmt.Sprint(t.Format("20060102150405"))
				fmt.Println("User told me to generate core dump")
				coredump("/tmp/go_" + timestamp + ".txt")
			// case syscall.SIGTERM:
			// 	fmt.Println("User told me to exit")
			// 	os.Exit(0)
			default:
				continue
			}
		}

	}()
}

func coredump(fileName string) {
	logrus.Infoln("Dump stacktrace to ", fileName)
	ioutil.WriteFile(fileName, []byte(StackTrace(true)), 0644)
}
