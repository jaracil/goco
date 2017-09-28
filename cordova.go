package cordova

import (
	"github.com/gopherjs/gopherjs/js"
)

var (
	deviceReady = false
)

func IsReady() bool {
	return deviceReady
}

func WaitReady() {
	if deviceReady {
		return
	}
	ch := make(chan struct{}, 0)
	f := func() {
		deviceReady = true
		close(ch)
	}
	js.Global.Get("document").Call("addEventListener", "deviceready", f, false)
	<-ch
	js.Global.Get("document").Call("removeEventListener", "deviceready", f, false)
}

func OnPause(cb func()) {
	js.Global.Get("document").Call("addEventListener", "pause", cb, false)
}

func OnResume(cb func()) {
	js.Global.Get("document").Call("addEventListener", "resume", cb, false)
}
