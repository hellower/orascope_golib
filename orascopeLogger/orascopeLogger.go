package orascopeLogger

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

	this.Version = "20190420"
	this.fs = nil
	this.seq = 0 // default

	this.fileName = ""
	this.fullLogDir = ""
	this.fullLogFilePath = ""
	this.console = true // 디폴트콘솔임 StartFile 전까지!!
	this.debug = true   // 디폴트 디버깅모드임.

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

func (this *ClsLogger) TracingMode(a_debug bool) {
	this.debug = a_debug
}

func (this *ClsLogger) StartFile(a_path string, a_fileName string, a_appname string) {
	this.console = false
	this.fullLogDir = a_path
	this.fileName = a_fileName
	this.fullLogFilePath = fmt.Sprintf("%s%c%s", this.fullLogDir, os.PathSeparator, this.fileName)

	var err error
	this.fs, err = os.OpenFile(this.fullLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	logrus.SetOutput(this.fs)

	this.log = logrus.WithFields(logrus.Fields{"name": a_appname})

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

func (this *ClsLogger) FatalNew(a_format string, a_args ...interface{}) {
	this.Fatal(fmt.Errorf(a_format, a_args...))
}

// Fatal(fmt.Errorf("%v",a_errr)) 이런식으로 사용가능
func (this *ClsLogger) Fatal(a_err error) {
	if a_err != nil {
		pc, fn, line, _ := runtime.Caller(1)

		// exit 내장됨!
		this.log.Errorf("{%v} @%s<%s:%d>", a_err, runtime.FuncForPC(pc).Name(), fn, line)
		this.CleanUp()
		os.Exit(-1)
	}
}

func (this *ClsLogger) Fatalf(a_err error, a_format string, a_args ...interface{}) {
	if a_err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		l_runtimeMsg := fmt.Sprintf(" {%v} @%s<%s:%d>", a_err, runtime.FuncForPC(pc).Name(), fn, line)

		this.log.Errorf(a_format+l_runtimeMsg, a_args...)
		this.CleanUp()
		os.Exit(-1)
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

func (this *ClsLogger) Info(a_format string, a_args ...interface{}) {
	if a_args == nil {
		this.log.Info(a_format)
	} else {
		this.log.Infof(a_format, a_args...)
	}
}

func (this *ClsLogger) Trace(a_format string, a_args ...interface{}) {
	if this.debug {
		pc, fn, line, _ := runtime.Caller(1)
		l_runtimeMsg := fmt.Sprintf(" @%s<%s:%d>", runtime.FuncForPC(pc).Name(), fn, line)
		if a_args == nil {
			this.log.Tracef("%s%s", a_format, l_runtimeMsg)
		} else {
			this.log.Tracef(a_format+l_runtimeMsg, a_args...)
		}
	}
}

func (this *ClsLogger) Console(format string, args ...interface{}) {
	if args == nil {
		//fmt.Printf(format+"\n")
		log.Println(format)
	} else {
		//fmt.Printf(format+"\n", args...)
		log.Printf(format+"\n", args...)
	}
}

func (this *ClsLogger) DebugConsole(format string, args ...interface{}) {
	if args == nil {
		//fmt.Printf(format+"\n")
		log.Println("@@DEBUG@@ " + format)
	} else {
		//fmt.Printf(format+"\n", args...)
		log.Printf("@@DEBUG@@ "+format+"\n", args...)
	}
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
