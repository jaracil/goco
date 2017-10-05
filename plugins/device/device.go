package device

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco/plugins/cordova"
)

var (
	DevInfo *DeviceInfo
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

func init() {
	cordova.OnDeviceReady(func() {
		DevInfo = &DeviceInfo{
			Object: js.Global.Get("device"),
		}
	})
}
