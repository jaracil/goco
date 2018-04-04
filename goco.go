// Package goco contains idiomatic Go bindings for cordova.
package goco

import (
	"github.com/gopherjs/gopherjs/js"
)

var (
	deviceReady = false
)

// WaitReady waits until cordova api is ready (see deviceready event)
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

// OnDeviceReady registers callback function that runs when device is ready
func OnDeviceReady(cb func()) {
	js.Global.Get("document").Call("addEventListener", "deviceready", cb, false)
}

// OnPause registers callback function that runs when app goes to background
func OnPause(cb func()) {
	js.Global.Get("document").Call("addEventListener", "pause", cb, false)
}

// OnResume registers callback function that runs when app goes to foreground
func OnResume(cb func()) {
	js.Global.Get("document").Call("addEventListener", "resume", cb, false)
}

// OnVolumeUpButton registers callback function that runs when volume-up button is pressed
func OnVolumeUpButton(cb func()) {
	js.Global.Get("document").Call("addEventListener", "volumeupbutton", cb, false)
}

// OnVolumeDownButton registers callback function that runs when volume-down button is pressed
func OnVolumeDownButton(cb func()) {
	js.Global.Get("document").Call("addEventListener", "volumedownbutton", cb, false)
}

// UnDeviceReady clears previous OnDeviceReady registration
func UnDeviceReady(cb func()) {
	js.Global.Get("document").Call("removeEventListener", "deviceready", cb, false)
}

// UnPause clears previous OnPause registration
func UnPause(cb func()) {
	js.Global.Get("document").Call("removeEventListener", "pause", cb, false)
}

// UnResume clears previous OnResume registration
func UnResume(cb func()) {
	js.Global.Get("document").Call("removeEventListener", "resume", cb, false)
}

// UnVolumeUpButton clears previous OnVolumeUpButton registration
func UnVolumeUpButton(cb func()) {
	js.Global.Get("document").Call("removeEventListener", "volumeupbutton", cb, false)
}

// UnVolumeDownButton clears previous OnVolumeDownButton registration
func UnVolumeDownButton(cb func()) {
	js.Global.Get("document").Call("removeEventListener", "volumedownbutton", cb, false)
}
