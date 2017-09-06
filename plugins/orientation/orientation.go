package orientation

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

type Heading struct {
	*js.Object

	MagneticHeading float64 `js:"magneticHeading"`
	TrueHeading     float64 `js:"trueHeading"`
	HeadingAccuracy float64 `js:"headingAccuracy"`
	Timestamp       int64   `js:"timestamp"`
}

type Watcher struct {
	id *js.Object
}

func wrapHeading(obj *js.Object) *Heading {
	return &Heading{Object: obj}
}

func CurrentHeading() (heading *Heading, err error) {
	ch := make(chan struct{})
	success := func(obj *js.Object) {
		heading = wrapHeading(obj)
		close(ch)
	}
	fail := func() {
		err = errors.New("Error getting compass data")
		close(ch)
	}

	js.Global.Get("navigator").Get("compass").Call("getCurrentHeading", success, fail)
	<-ch
	return
}

func NewWatcher(cb func(*Heading, error), options map[string]interface{}) *Watcher {
	success := func(obj *js.Object) {
		h := wrapHeading(obj)
		cb(h, nil)
	}

	fail := func() {
		err := errors.New("Error getting compass data")
		cb(nil, err)
	}

	id := js.Global.Get("navigator").Get("compass").Call("watchHeading", success, fail, options)
	return &Watcher{id: id}
}

func (w *Watcher) Close() {
	js.Global.Get("navigator").Get("compass").Call("clearWatch", w.id)
}
