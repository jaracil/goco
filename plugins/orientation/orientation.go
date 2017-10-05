package orientation

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco/plugins/cordova"
)

type Heading struct {
	*js.Object
	MagneticHeading float64 `js:"magneticHeading"`
	TrueHeading     float64 `js:"trueHeading"`
	HeadingAccuracy float64 `js:"headingAccuracy"`
	Timestamp       int64   `js:"timestamp"`
}

type Watcher struct {
	*js.Object
}

var mo *js.Object

func init() {
	cordova.OnDeviceReady(func() {
		mo = js.Global.Get("navigator").Get("compass")
	})
}

func CurrentHeading() (heading *Heading, err error) {
	ch := make(chan struct{})
	success := func(h *Heading) {
		heading = h
		close(ch)
	}
	fail := func() {
		err = errors.New("Error getting compass data")
		close(ch)
	}

	mo.Call("getCurrentHeading", success, fail)
	<-ch
	return
}

func NewWatcher(cb func(*Heading, error), options map[string]interface{}) *Watcher {
	success := func(h *Heading) {
		cb(h, nil)
	}

	fail := func() {
		err := errors.New("Error getting compass data")
		cb(nil, err)
	}

	id := mo.Call("watchHeading", success, fail, options)
	return &Watcher{Object: id}
}

func (w *Watcher) Close() {
	mo.Call("clearWatch", w)
}
