package device

import (
	"crypto/rand"
	"encoding/hex"

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
		if DevInfo.Platform == "browser" {
			print("In browser platform")
			DevInfo.Manufacturer = js.Global.Get("navigator").Get("vendor").String()
			if DevInfo.Manufacturer == "" {
				DevInfo.Manufacturer = "Mozilla"
			}
			res := js.Global.Get("localStorage").Get("cordova_browser_uuid")
			if res == js.Undefined {
				uuid := make([]byte, 16)
				rand.Read(uuid)
				DevInfo.UUID = hex.EncodeToString(uuid)
				js.Global.Get("localStorage").Call("setItem", "cordova_browser_uuid", DevInfo.UUID)
				print("Generating new UUID for browser platform: " + DevInfo.UUID)
			} else {
				DevInfo.UUID = res.String()
				print("browser UUID: " + DevInfo.UUID)
			}
		}
	})
}
