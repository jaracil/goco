package cordova

import (
	"github.com/gopherjs/gopherjs/js"
)

var (
	deviceReady = false
)

func WaitReady() {
	if deviceReady {
		return
	}
	ch := make(chan struct{}, 0)
	f := func() {
		deviceReady = true
		close(ch)
	}
	OnDeviceReady(f)
	<-ch
	UnDeviceReady(f)
}

func OnDeviceReady(cb func()) {
	js.Global.Get("document").Call("addEventListener", "deviceready", cb, false)
}

func OnPause(cb func()) {
	js.Global.Get("document").Call("addEventListener", "pause", cb, false)
}

func OnResume(cb func()) {
	js.Global.Get("document").Call("addEventListener", "resume", cb, false)
}

func UnDeviceReady(cb func()) {
	js.Global.Get("document").Call("removeEventListener", "deviceready", cb, false)
}

func UnPause(cb func()) {
	js.Global.Get("document").Call("removeEventListener", "pause", cb, false)
}

func UnResume(cb func()) {
	js.Global.Get("document").Call("removeEventListener", "resume", cb, false)
}
