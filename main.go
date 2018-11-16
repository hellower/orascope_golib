package main

import (
	"orascope_golib/anotherlib"
	"orascope_golib/orascopeLogger"
	"os"
	"time"
	//	"github.com/pkg/errors"
)

var Glog *orascopeLogger.ClsLogger // global 변수임

func init() {
	Glog = orascopeLogger.Glog()
	Glog.TracingMode(true)
	Glog.Trace("main.init") // 디폴트 콘솔임.
}

func main() {

	Glog.Trace("main()")
	Glog.Trace("%v", time.Now())
	Glog.Info("version=%s", Glog.Version)

	// 아래 안하면 디폴트로 콘솔로 출력!!!
	//Glog.StartFile("r:\\", "abc.log") // StartFile 전까지 모두 콘솔
	var err error

	//Glog.FatalNew("start fatal...") // 자체 에러 생산
	//-----------------------------------------------------------------------------------------
	if _, err = os.OpenFile("r:\\X.X", os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		Glog.Fatal(err)
		Glog.Fatalf(err, "화일 오픈에러#3")
		Glog.Fatalf(err, "화일 오픈에러#3=>%s", "x:\\x.x")

		//		Glog.Fatal(errors.New("file open error"))
		//		Glog.Fatal(errors.Wrap(err, "화일 오픈에러#3"))
		//		Glog.Fatal(errors.Wrapf(err, "화일 오픈에러#3=>%s", "x:\\x.x"))
	}
	//-----------------------------------------------------------------------------------------
	an := new(another.Another)
	if err = an.ErrTest(); err != nil {
		Glog.Fatal(err)
	}
	if err = an.ErrTest2(); err != nil {
		Glog.Fatal(err)
	}

	an.InfoTest()
	//	an.InfoTest2()
	//an.FatalTest()
	//-----------------------------------------------------------------------------------------
	_, err = os.OpenFile("z:/a.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Glog.Trace("%+v", "화일오픈에러#1")
		Glog.Trace("%+v, %s", err, "화일에러#2")
		Glog.Trace("%+v 화일에러#2=>%s", err, "z:/a.txt")

		//		Glog.Trace("%+v", errors.New("화일오픈에러#1"))
		//		Glog.Trace("%+v", errors.Wrap(err, "화일에러#2"))
		//		Glog.Trace("%+v", errors.Wrapf(err, "화일에러#2=>%s", "z:/a.txt"))
		//os.Exit(1)
		// main 함수가 끝나기 직전에 파일을 닫음
	} // end if
	//-----------------------------------------------------------------------------------------

	Glog.Info("%d##%s", 123, "batman")
	Glog.Trace("%d##%s", 123, "batman")

	//Glog.Info("superman")

	defer func() {
		Glog.CleanUp()
	}()
}
