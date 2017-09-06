package geolocation

import (
	"errors"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

type Coords struct {
	*js.Object

	Latitude         float64 `js:"latitude"`
	Longitude        float64 `js:"longitude"`
	Altitude         float64 `js:"altitude"`
	Accuracy         float64 `js:"accuracy"`
	AltitudeAccuracy float64 `js:"altitudeAccuracy"`
	Heading          float64 `js:"heading"`
	Speed            float64 `js:"speed"`
}

type Position struct {
	Coords    *Coords
	Timestamp time.Time
}

type Watcher struct {
	id *js.Object
}

func wrapPosition(obj *js.Object) *Position {
	pos := &Position{}
	uts := obj.Get("timestamp").Int64()
	pos.Timestamp = time.Unix(uts/1000, (uts%1000)*1000000)
	pos.Coords = &Coords{Object: obj.Get("coords")}
	return pos
}

func CurrentPosition(options map[string]interface{}) (pos *Position, err error) {
	ch := make(chan struct{})
	success := func(obj *js.Object) {
		pos = wrapPosition(obj)
		close(ch)
	}
	fail := func(obj *js.Object) {
		err = errors.New(obj.Get("message").String())
		close(ch)
	}

	js.Global.Get("navigator").Get("geolocation").Call("getCurrentPosition", success, fail, options)
	<-ch
	return
}

func NewWatcher(cb func(pos *Position, err error), options map[string]interface{}) *Watcher {
	success := func(obj *js.Object) {
		pos := wrapPosition(obj)
		cb(pos, nil)
	}

	fail := func(obj *js.Object) {
		err := errors.New(obj.Get("message").String())
		cb(nil, err)
	}

	id := js.Global.Get("navigator").Get("geolocation").Call("watchPosition", success, fail, options)
	return &Watcher{id: id}
}

func (w *Watcher) Close() {
	js.Global.Get("navigator").Get("geolocation").Call("clearWatch", w.id)
}
