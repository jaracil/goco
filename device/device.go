// Package device is a GopherJS wrapper for cordova device plugin.
// This plugin provides a public DevInfo var, which describes the device's hardware and software.
//
// Install plugin:
//  cordova plugin add cordova-plugin-device

package device

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco"
)

var (
	// DevInfo is a Instance of device.Info type
	DevInfo *Info
)

// Info declares device's hardware and software info.
type Info struct {
	*js.Object
	Cordova      string `js:"cordova"`      // Version of Cordova running on the device.
	Model        string `js:"model"`        // Name of the device's model or product. The value is set by the device manufacturer and may be different across versions of the same product.
	Platform     string `js:"platform"`     // Device's operating system name.
	UUID         string `js:"uuid"`         // Device's Universally Unique Identifier.
	Version      string `js:"version"`      // Operating system version.
	Manufacturer string `js:"manufacturer"` // Device's manufacturer.
	IsVirtual    bool   `js:"isVirtual"`    // Whether the device is running on a simulator.
	Serial       string `js:"serial"`       // Device's hardware serial number
}

func init() {
	goco.OnDeviceReady(func() {
		DevInfo = &Info{
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
