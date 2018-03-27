// Package location is a GopherJS wrapper for cordova diagnostic/location plugin.
//
// Install plugin:
//  cordova plugin add cordova.plugins.diagnostic
package location

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco"
)

// Modes defines Android location modes
type Modes struct {
	*js.Object
	HighAccuracy  string `js:"HIGH_ACCURACY"`
	BatterySaving string `js:"BATTERY_SAVING"`
	DeviceOnly    string `js:"DEVICE_ONLY"`
	LocationOff   string `js:"LOCATION_OFF"`
}

var (
	instance *js.Object

	// Mode is an instance of Modes
	Mode *Modes
)

func init() {
	goco.OnDeviceReady(func() {
		Mode = &Modes{Object: mo().Get("locationMode")}
	})
}

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("cordova").Get("plugins").Get("diagnostic")
	}
	return instance
}

// IsLocationAvailable returns true when location is available.
//  Platforms: Android, iOS and Windows 10 UWP
func IsLocationAvailable() (res bool, err error) {
	ch := make(chan struct{})
	success := func(st bool) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("isLocationAvailable", success, fail)
	<-ch
	return
}

// IsLocationEnabled returns true if the device setting for location is on.
// On Android this returns true if Location Mode is switched on.
// On iOS this returns true if Location Services is switched on.
//  Platforms: Android and iOS
func IsLocationEnabled() (res bool, err error) {
	ch := make(chan struct{})
	success := func(st bool) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("isLocationEnabled", success, fail)
	<-ch
	return
}

// IsGpsLocationAvailable checks if high-accuracy locations are available to the app from GPS hardware.
// Returns true if Location mode is enabled and is set to "Device only" or "High accuracy" AND if the app is authorized to use location.
//  Platforms: Android
func IsGpsLocationAvailable() (res bool, err error) {
	ch := make(chan struct{})
	success := func(st bool) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("isGpsLocationAvailable", success, fail)
	<-ch
	return
}

// IsGpsLocationEnabled checks if the device location setting is set to return high-accuracy locations from GPS hardware.
//  Platforms: Android
func IsGpsLocationEnabled() (res bool, err error) {
	ch := make(chan struct{})
	success := func(st bool) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("isGpsLocationEnabled", success, fail)
	<-ch
	return
}

// IsNetworkLocationAvailable checks if low-accuracy locations are available to the app from network triangulation/WiFi access points.
//  Platforms: Android
func IsNetworkLocationAvailable() (res bool, err error) {
	ch := make(chan struct{})
	success := func(st bool) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("isNetworkLocationAvailable", success, fail)
	<-ch
	return
}

// IsNetworkLocationEnabled checks if location mode is set to return low-accuracy locations from network triangulation/WiFi access points.
//  Platforms: Android
func IsNetworkLocationEnabled() (res bool, err error) {
	ch := make(chan struct{})
	success := func(st bool) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("isNetworkLocationEnabled", success, fail)
	<-ch
	return
}

// GetLocationMode returns the current location mode setting for the device.
//  Platforms: Android
func GetLocationMode() (res string, err error) {
	ch := make(chan struct{})
	success := func(st string) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("getLocationMode", success, fail)
	<-ch
	return
}

// IsLocationAuthorized checks if the application is authorized to use location.
//  Platforms: Android and iOS
func IsLocationAuthorized() (res bool, err error) {
	ch := make(chan struct{})
	success := func(st bool) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("isLocationAuthorized", success, fail)
	<-ch
	return
}

// GetLocationAuthorizationStatus returns the location authorization status for the application.
//  Platforms: Android and iOS
func GetLocationAuthorizationStatus() (res string, err error) {
	ch := make(chan struct{})
	success := func(st string) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	mo().Call("getLocationAuthorizationStatus", success, fail)
	<-ch
	return
}

// RequestLocationAuthorization requests location authorization for the application.
//  Platforms: Android and iOS
func RequestLocationAuthorization(mode ...[]string) (res string, err error) {
	ch := make(chan struct{})
	success := func(st string) {
		res = st
		close(ch)
	}
	fail := func(s string) {
		err = errors.New(s)
		close(ch)
	}
	if len(mode) == 0 {
		mo().Call("requestLocationAuthorization", success, fail)
	} else {
		mo().Call("requestLocationAuthorization", success, fail, mode[0])
	}

	<-ch
	return
}

// SwitchToLocationSettings displays the device location settings to allow user to enable location services/change location mode.
//  Platforms: Android and Windows 10 UWP
func SwitchToLocationSettings() {
	mo().Call("switchToLocationSettings")
}

// RegisterLocationStateChangeHandler registers a function to be called when a change in Location state occurs.
// Pass in nil value to de-register the currently registered function.
//  Platforms: Android and iOS
func RegisterLocationStateChangeHandler(f func(string)) {
	mo().Call("registerLocationStateChangeHandler", f)
}
