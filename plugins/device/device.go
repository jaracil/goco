package device

import (
	"github.com/gopherjs/gopherjs/js"
)

var (
	singleton   *DeviceInfo
	deviceReady = false
)

type DeviceInfo struct {
	*js.Object
	Cordova      string `js:"cordova"`
	Model        string `js:"model"`
	Platform     string `js:"platform"`
	UUID         string `js:"uuid"`
	Version      string `js:"version"`
	Manufacturer string `js:"manufacturer"`
	// whether the device is running on a simulator.
	IsVirtual bool `js:"isVirtual"`
	// Get the device hardware serial number
	Serial string `js:"serial"`
}

func Info() *DeviceInfo {
	if singleton != nil {
		return singleton
	}
	singleton = &DeviceInfo{
		Object: js.Global.Get("device"),
	}
	return singleton
}

func IsReady() bool {
	return deviceReady
}

func WaitReady() {
	if deviceReady {
		return
	}
	ch := make(chan struct{}, 0)
	f := func() {
		close(ch)
	}
	js.Global.Get("document").Call("addEventListener", "deviceready", f, false)
	<-ch
	deviceReady = true
	js.Global.Get("document").Call("removeEventListener", "deviceready", f, false)
}
