package another

import (
	"orascope_golib/orascopeLogger"
	"os"

	"github.com/pkg/errors"
)

type Another struct {
	a int
}

func (this *Another) ErrTest() error {
	var err error
	if _, err = os.OpenFile("Y:\\Y.Y", os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return errors.Wrap(err, "another.go.ErrTest()")
	}
	return nil
}

func (this *Another) ErrTest2() error {
	var err error
	if _, err = os.OpenFile("Y:\\Y.Y", os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return errors.Wrapf(err, "another.go.%s", "ErrTest2()")
	}
	return nil
}

func (this *Another) ErrTest3() error {
	return errors.New("ErrTest3()")
}
func (this *Another) ErrTest4() error {
	return errors.Errorf("에러 %s 발생", "ErrTest4()")
}

func (this *Another) FatalTest() {
	var err error
	if _, err = os.OpenFile("Y:\\Y.Y", os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		orascopeLogger.Glog().Fatal(err)
	}
}

func (this *Another) InfoTest() {
	var err error
	if _, err = os.OpenFile("Y:\\Y.Y", os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		orascopeLogger.Glog().Trace("another.InfoTest()")
	}
}

func (this *Another) InfoTest2() {
	var err error
	if _, err = os.OpenFile("Y:\\Y.Y", os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		orascopeLogger.Glog().Trace("%+v", errors.Wrap(err, "another.InfoTest2()"))
	}
}

func BornAnother() (this *Another) {
	this = new(Another)
	return
}

func init() { // main 보다 우선 수행됨!!!
	//fmt.Printf("another.init()")
}
