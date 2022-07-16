package orascopeLogger

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

type ClsLogger struct {
	fs              *os.File
	seq             uint32
	fileName        string
	fullLogDir      string
	fullLogFilePath string
	console         bool
	debug           bool
	Version         string
	log             *logrus.Entry
}

func BornClsLogger() (this *ClsLogger) {
	this = new(ClsLogger)

	//	this.Version = "20190420"
	//	this.Version = "20220619"
	this.Version = "20220716" // remove file fullpath
	this.fs = nil
	this.seq = 0 // default

	this.fileName = ""
	this.fullLogDir = ""
	this.fullLogFilePath = ""
	this.console = true // 디폴트콘솔임 StartFile 전까지!!
	this.debug = false  // 디폴트 디버깅모드임.

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	this.log = logrus.WithFields(logrus.Fields{"name": "console"})
	return
}

/*
// for SetOutput
// 헤더 장착!!
func (this *ClsLogger) Write(bytes []byte) (int, error) {
	this.seq = this.seq + 1
	return fmt.Fprintf(this.fs, "[%d,%s,%s", this.seq, time.Now().Format("15:04:05"), string(bytes))
}

func (this *ClsLogger) StartConsole() {
	this.console = true
}
*/

func (this *ClsLogger) TracingMode(aDebug bool) {
	this.debug = aDebug
}

func (this *ClsLogger) StartFile(aPath string, aFilename string, aAppname string) {
	this.console = false
	this.fullLogDir = aPath
	this.fileName = aFilename
	this.fullLogFilePath = fmt.Sprintf("%s%c%s", this.fullLogDir, os.PathSeparator, this.fileName)

	var err error
	this.fs, err = os.OpenFile(this.fullLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	logrus.SetOutput(this.fs)

	this.log = logrus.WithFields(logrus.Fields{"name": aAppname})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
		os.Exit(1)
		// main 함수가 끝나기 직전에 파일을 닫음
	} // end if
}

func (this *ClsLogger) CleanUp() {
	if this.fs != nil {
		this.fs.Close()
		this.fs = nil
	}
	//fmt.Printf("\nCleanUp")
}

/*
Fatalnew("에러");
Fatalnew("%s","발생");
*/

func (this *ClsLogger) FatalNew(aFormat string, aArgs ...interface{}) {
	this.Fatal(fmt.Errorf(aFormat, aArgs...))
}

// Fatal(fmt.Errorf("%v",rErr)) 이런식으로 사용가능
func (this *ClsLogger) Fatal(rErr error) {
	if rErr != nil {
		pc, fn, line, _ := runtime.Caller(1)

		// exit 내장됨!
		this.log.Errorf("{%v} @%s:%d  %s", rErr, filepath.Base(fn), line, runtime.FuncForPC(pc).Name())
		this.CleanUp()
		os.Exit(-1)
	}
}

func (this *ClsLogger) Fatalf(rErr error, aFormat string, aArgs ...interface{}) {
	if rErr != nil {
		pc, fn, line, _ := runtime.Caller(1)
		l_runtimeMsg := fmt.Sprintf(" {%v} @%s:%d  %s", rErr, filepath.Base(fn), line, runtime.FuncForPC(pc).Name())

		this.log.Errorf(aFormat+l_runtimeMsg, aArgs...)
		this.CleanUp()
		os.Exit(-1)
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

func (this *ClsLogger) Info(aFormat string, aArgs ...interface{}) {
	if aArgs == nil {
		this.log.Info(aFormat)
	} else {
		this.log.Infof(aFormat, aArgs...)
	}
}

func (this *ClsLogger) Trace(aFormat string, aArgs ...interface{}) {
	if this.debug {
		pc, fn, line, _ := runtime.Caller(1)
		l_runtimeMsg := fmt.Sprintf(" @%s:%d  %s", filepath.Base(fn), line, runtime.FuncForPC(pc).Name())
		if aArgs == nil {
			this.log.Tracef("%s%s", aFormat, l_runtimeMsg)
		} else {
			this.log.Tracef(aFormat+l_runtimeMsg, aArgs...)
		}
	}
}

func (this *ClsLogger) Console(aFormat string, aArgs ...interface{}) {
	if aArgs == nil {
		//fmt.Printf(format+"\n")
		log.Println(aFormat)
	} else {
		//fmt.Printf(format+"\n", args...)
		log.Printf(aFormat+"\n", aArgs...)
	}
}

func (this *ClsLogger) DebugConsole(aFormat string, aArgs ...interface{}) {
	if this.debug {
		pc, fn, line, _ := runtime.Caller(1)
		l_runtimeMsg := fmt.Sprintf(" @%s:%d  %s", filepath.Base(fn), line, runtime.FuncForPC(pc).Name())
		if aArgs == nil {
			log.Printf("@@DEBUG@@ %s%s\n", aFormat, l_runtimeMsg)
		} else {
			log.Printf("@@DEBUG@@ "+aFormat+l_runtimeMsg+"\n", aArgs...)
		}
	}
	/*
		if this.debug {
			if args == nil {
				//fmt.Printf(format+"\n")
				log.Println("@@DEBUG@@ " + format)
			} else {
				//fmt.Printf(format+"\n", args...)
				log.Printf("@@DEBUG@@ "+format+"\n", args...)
			}
		}
	*/
}

func (this *ClsLogger) Catch() {
	if l_err := recover(); l_err != nil {
		//this.Fatalf(fmt.Errorf("%v", l_err), "CATCH")
		this.Fatal(errors.Errorf("%v", l_err)) // 2017.07.29
	}
}

/********************************************************************************/

var _g_ClsLogger_ *ClsLogger // global 변수임
func init() { // main 보다 우선 수행됨!!!
	_g_ClsLogger_ = BornClsLogger()
	//fmt.Printf("mylogger.init()")
}

func Glog() *ClsLogger {
	return _g_ClsLogger_
}
